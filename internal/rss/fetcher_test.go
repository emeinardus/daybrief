package rss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetcherFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = fmt.Fprint(w, "<rss><channel></channel></rss>")
	}))
	defer server.Close()

	fetcher := NewFetcher()
	body, err := fetcher.Fetch(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("Fetch() returned an unexpected error: %v", err)
	}

	if string(body) != "<rss><channel></channel></rss>" {
		t.Errorf("unexpected response body: %q", string(body))
	}
}

func TestFetcherRejectsNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "temporarily unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, err := NewFetcher().Fetch(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected an error for a non-200 response")
	}
}

func TestFetcherRejectsOversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, strings.Repeat("x", int(maxResponseBytes+1)))
	}))
	defer server.Close()

	_, err := NewFetcher().Fetch(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected an error for an oversized response")
	}
}
