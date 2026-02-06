package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	githubUsername = "testuser"
}

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

func TestHandleWebhook_PullRequestReview_OwnPR(t *testing.T) {
	payload := `{
		"action": "submitted",
		"pull_request": {
			"number": 42,
			"title": "Add feature",
			"user": {"login": "testuser"}
		},
		"review": {
			"state": "approved",
			"user": {"login": "reviewer"}
		},
		"sender": {"login": "reviewer"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleWebhook_PullRequestReview_OthersPR(t *testing.T) {
	payload := `{
		"action": "submitted",
		"pull_request": {
			"number": 42,
			"title": "Add feature",
			"user": {"login": "otheruser"}
		},
		"review": {
			"state": "approved",
			"user": {"login": "testuser"}
		},
		"sender": {"login": "testuser"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	// Should still return OK but not notify (filtered out)
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleWebhook_IgnoreSelfAction(t *testing.T) {
	payload := `{
		"action": "submitted",
		"pull_request": {
			"number": 42,
			"title": "Add feature",
			"user": {"login": "testuser"}
		},
		"review": {
			"state": "commented",
			"user": {"login": "testuser"}
		},
		"sender": {"login": "testuser"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	// Should return OK but not notify (self action)
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
			"user": {"login": "testuser"}
		},
		"comment": {
			"body": "Nice work!",
			"user": {"login": "commenter"}
		},
		"sender": {"login": "commenter"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request_review_comment")
	w := httptest.NewRecorder()

	handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
