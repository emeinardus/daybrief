package source

import (
	"context"
	"testing"
)

func TestMemoryRegistryListsOnlyActiveSources(t *testing.T) {
	registry := NewMemoryRegistry([]FeedSource{
		{Name: "Example News", URL: "https://example.com/rss", Active: true},
		{Name: "Disabled News", URL: "https://disabled.example/rss", Active: false},
		{Name: "Tech News", URL: "https://tech.example/rss", Active: true},
	})

	sources, err := registry.ListActive(context.Background())
	if err != nil {
		t.Fatalf("ListActive() returned an unexpected error: %v", err)
	}

	if len(sources) != 2 {
		t.Fatalf("expected 2 active sources, got %d", len(sources))
	}

	if sources[0].Name != "Example News" {
		t.Errorf("unexpected first source: %q", sources[0].Name)
	}

	if sources[1].Name != "Tech News" {
		t.Errorf("unexpected second source: %q", sources[1].Name)
	}
}

func TestMemoryRegistryRespectsCanceledContext(t *testing.T) {
	registry := NewMemoryRegistry(nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := registry.ListActive(ctx)
	if err == nil {
		t.Fatal("expected an error for a canceled context")
	}
}
