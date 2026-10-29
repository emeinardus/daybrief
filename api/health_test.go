package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected application/json, got %q", contentType)
	}

	expected := "{\"status\":\"ok\"}\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q",
			expected, rec.Body.String())
	}
}
