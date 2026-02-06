package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// GitHub webhook payload structs
type User struct {
	Login string `json:"login"`
}

type PullRequest struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	User    User   `json:"user"`
	HTMLURL string `json:"html_url"`
}

type Issue struct {
	Number      int                    `json:"number"`
	Title       string                 `json:"title"`
	User        User                   `json:"user"`
	HTMLURL     string                 `json:"html_url"`
	PullRequest map[string]interface{} `json:"pull_request"` // present if issue is a PR
}

type Review struct {
	State string `json:"state"`
	User  User   `json:"user"`
	Body  string `json:"body"`
}

type Comment struct {
	Body string `json:"body"`
	User User   `json:"user"`
}

type WebhookPayload struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Issue       Issue       `json:"issue"`
	Review      Review      `json:"review"`
	Comment     Comment     `json:"comment"`
	Sender      User        `json:"sender"`
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func verifySignature(body []byte, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(body)
	return hmac.Equal(sig, mac.Sum(nil))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(body, signature) {
		log.Printf("Invalid signature")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")
	log.Printf("Received webhook: %s", eventType)
	if debug {
		log.Printf("Webhook body: %s", string(body))
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Determine PR author based on event type
	var prAuthor string
	if eventType == "issue_comment" {
		// issue_comment uses "issue" field
		if payload.Issue.PullRequest == nil {
			log.Printf("Ignoring: not a pull request")
			w.WriteHeader(http.StatusOK)
			return
		}
		prAuthor = payload.Issue.User.Login
	} else {
		prAuthor = payload.PullRequest.User.Login
	}

	// Only notify for PRs authored by the configured user
	if prAuthor != githubUsername {
		log.Printf("Ignoring: PR author %s is not %s", prAuthor, githubUsername)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Don't notify for own actions
	if payload.Sender.Login == githubUsername {
		log.Printf("Ignoring: action from self")
		w.WriteHeader(http.StatusOK)
		return
	}

	var msg string
	switch eventType {
	case "pull_request_review":
		msg = fmt.Sprintf("PR #%d <%s|%s> received *%s* review from %s",
			payload.PullRequest.Number,
			payload.PullRequest.HTMLURL,
			payload.PullRequest.Title,
			payload.Review.State,
			payload.Review.User.Login)
		if payload.Review.Body != "" {
			msg += ": " + truncate(payload.Review.Body, 1000)
		}

	case "pull_request_review_comment":
		msg = fmt.Sprintf("PR #%d <%s|%s> new comment from %s: %s",
			payload.PullRequest.Number,
			payload.PullRequest.HTMLURL,
			payload.PullRequest.Title,
			payload.Comment.User.Login,
			truncate(payload.Comment.Body, 1000))

	case "issue_comment":
		msg = fmt.Sprintf("PR #%d <%s|%s> new comment from %s: %s",
			payload.Issue.Number,
			payload.Issue.HTMLURL,
			payload.Issue.Title,
			payload.Comment.User.Login,
			truncate(payload.Comment.Body, 1000))

	default:
		log.Printf("Ignoring event type: %s", eventType)
	}

	if msg != "" {
		if err := sendSlackMessage(msg); err != nil {
			log.Printf("Error sending Slack message: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}
