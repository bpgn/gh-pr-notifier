package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	githubUsername = "testuser"
	webhookSecret = "testsecret"
	slackBotToken = "test-token"
	slackUserID = "U123"
}

func signPayload(payload string) string {
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestHandleWebhook_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandleWebhook_InvalidSignature(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader("{}"))
	req.Header.Set("X-Hub-Signature-256", "sha256=invalid")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandleWebhook_MissingSignature(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandleWebhook_InvalidJSON(t *testing.T) {
	payload := "not json"
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	req.Header.Set("X-Hub-Signature-256", signPayload(payload))
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleWebhook_PullRequestReview_OwnPR(t *testing.T) {
	payload := `{"action":"submitted","pull_request":{"number":42,"title":"Add feature","user":{"login":"testuser"}},"review":{"state":"approved","user":{"login":"reviewer"}},"sender":{"login":"reviewer"}}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	req.Header.Set("X-Hub-Signature-256", signPayload(payload))
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleWebhook_PullRequestReview_OthersPR(t *testing.T) {
	payload := `{"action":"submitted","pull_request":{"number":42,"title":"Add feature","user":{"login":"otheruser"}},"review":{"state":"approved","user":{"login":"testuser"}},"sender":{"login":"testuser"}}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	req.Header.Set("X-Hub-Signature-256", signPayload(payload))
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleWebhook_IgnoreSelfAction(t *testing.T) {
	payload := `{"action":"submitted","pull_request":{"number":42,"title":"Add feature","user":{"login":"testuser"}},"review":{"state":"commented","user":{"login":"testuser"}},"sender":{"login":"testuser"}}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	req.Header.Set("X-Hub-Signature-256", signPayload(payload))
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleWebhook_PullRequestReviewComment(t *testing.T) {
	payload := `{"action":"created","pull_request":{"number":42,"title":"Add feature","user":{"login":"testuser"}},"comment":{"body":"Nice work!","user":{"login":"commenter"}},"sender":{"login":"commenter"}}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review_comment")
	req.Header.Set("X-Hub-Signature-256", signPayload(payload))
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
