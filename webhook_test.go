package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleWebhook_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandleWebhook_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader("not json"))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleWebhook_PullRequestReview(t *testing.T) {
	payload := `{
		"action": "submitted",
		"pull_request": {
			"number": 42,
			"title": "Add feature",
			"user": {"login": "author"},
			"html_url": "https://github.com/org/repo/pull/42"
		},
		"review": {
			"state": "approved",
			"user": {"login": "reviewer"},
			"body": "LGTM"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleWebhook_PullRequestReviewComment(t *testing.T) {
	payload := `{
		"action": "created",
		"pull_request": {
			"number": 42,
			"title": "Add feature",
			"user": {"login": "author"}
		},
		"comment": {
			"body": "Nice work!",
			"user": {"login": "commenter"}
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review_comment")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleWebhook_UnknownEvent(t *testing.T) {
	payload := `{"action": "opened"}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "issues")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
