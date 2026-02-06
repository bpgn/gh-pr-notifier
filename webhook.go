package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
	Review      Review      `json:"review"`
	Comment     Comment     `json:"comment"`
	Sender      User        `json:"sender"`
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")
	log.Printf("Received webhook: %s", eventType)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Only notify for PRs authored by the configured user
	if payload.PullRequest.User.Login != githubUsername {
		log.Printf("Ignoring: PR author %s is not %s", payload.PullRequest.User.Login, githubUsername)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Don't notify for own actions
	if payload.Sender.Login == githubUsername {
		log.Printf("Ignoring: action from self")
		w.WriteHeader(http.StatusOK)
		return
	}

	switch eventType {
	case "pull_request_review":
		log.Printf("[NOTIFY] PR #%d '%s' received %s review from %s",
			payload.PullRequest.Number,
			payload.PullRequest.Title,
			payload.Review.State,
			payload.Review.User.Login)

	case "pull_request_review_comment":
		log.Printf("[NOTIFY] PR #%d '%s' received comment from %s: %s",
			payload.PullRequest.Number,
			payload.PullRequest.Title,
			payload.Comment.User.Login,
			payload.Comment.Body)

	default:
		log.Printf("Ignoring event type: %s", eventType)
	}

	w.WriteHeader(http.StatusOK)
}
