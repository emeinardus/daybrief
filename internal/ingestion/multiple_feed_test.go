package ingestion

import (
	"context"
	"errors"
	"testing"

	"github.com/emeinardus/daybrief/internal/source"
)

type mapFetcher struct {
	responses map[string][]byte
	errors    map[string]error
}

func (f *mapFetcher) Fetch(
	_ context.Context,
	feedURL string,
) ([]byte, error) {
	if err := f.errors[feedURL]; err != nil {
		return nil, err
	}

	data, ok := f.responses[feedURL]
	if !ok {
		return nil, errors.New("feed URL not configured in test")
	}

	return data, nil
}

func TestLoadActiveFeedsCollectsArticlesAndContinuesAfterError(t *testing.T) {
	const (
		firstURL  = "https://example.com/rss"
		secondURL = "https://broken.example/rss"
		thirdURL  = "https://tech.example/rss"
	)

	fetcher := &mapFetcher{
		responses: map[string][]byte{
			firstURL: []byte(`<rss><channel><item>
				<title>First Article</title>
				<link>https://example.com/1</link>
			</item></channel></rss>`),
			thirdURL: []byte(`<rss><channel><item>
				<title>Third Article</title>
				<link>https://tech.example/3</link>
			</item></channel></rss>`),
		},
		errors: map[string]error{
			secondURL: errors.New("temporary network failure"),
		},
	}

	registry := source.NewMemoryRegistry([]source.FeedSource{
		{Name: "Example News", URL: firstURL, Active: true},
		{Name: "Broken News", URL: secondURL, Active: true},
		{Name: "Tech News", URL: thirdURL, Active: true},
		{Name: "Disabled News", URL: "https://disabled.example/rss", Active: false},
	})

	service := NewService(fetcher, registry)

	articles, results, err := service.LoadActiveFeeds(context.Background())
	if err != nil {
		t.Fatalf("LoadActiveFeeds() returned an unexpected error: %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("expected 2 articles, got %d", len(articles))
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 feed results, got %d", len(results))
	}

	if results[0].Count != 1 || results[0].Error != nil {
		t.Errorf("unexpected first feed result: %+v", results[0])
	}

	if results[1].Error == nil {
		t.Error("expected an error for the second feed")
	}

	if results[2].Count != 1 || results[2].Error != nil {
		t.Errorf("unexpected third feed result: %+v", results[2])
	}

	if articles[0].Source != "Example News" {
		t.Errorf("unexpected first article source: %q", articles[0].Source)
	}

	if articles[1].Source != "Tech News" {
		t.Errorf("unexpected second article source: %q", articles[1].Source)
	}
}
