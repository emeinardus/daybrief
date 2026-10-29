package rss

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	input := []byte(`<?xml version="1.0"?>
<rss version="2.0">
  <channel>
    <title>Example News</title>
    <item>
      <title> Technology Update </title>
      <link>https://example.com/article</link>
      <description>A technology update.</description>
      <pubDate>Thu, 29 Oct 2026 10:00:00 GMT</pubDate>
    </item>
    <item>
      <title>Article Without Date</title>
      <link>https://example.com/second</link>
    </item>
  </channel>
</rss>`)

	articles, err := Parse(input, "Example News")
	if err != nil {
		t.Fatalf("Parse() returned an unexpected error: %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("expected 2 articles, got %d", len(articles))
	}

	first := articles[0]

	if first.Title != "Technology Update" {
		t.Errorf("expected trimmed title, got %q", first.Title)
	}

	if first.URL != "https://example.com/article" {
		t.Errorf("unexpected URL: %q", first.URL)
	}

	if first.Source != "Example News" {
		t.Errorf("unexpected source: %q", first.Source)
	}

	if first.PublishedAt == nil {
		t.Fatal("expected publication date to be parsed")
	}

	expected := time.Date(2026, time.October, 29, 10, 0, 0, 0, time.UTC)
	if !first.PublishedAt.Equal(expected) {
		t.Errorf("expected publication date %v, got %v", expected, *first.PublishedAt)
	}

	if articles[1].PublishedAt != nil {
		t.Error("expected missing publication date to remain nil")
	}
}

func TestParseInvalidXML(t *testing.T) {
	_, err := Parse([]byte("<rss><channel>"), "Example News")
	if err == nil {
		t.Fatal("expected an error for malformed XML")
	}
}

func TestParseInvalidRoot(t *testing.T) {
	_, err := Parse([]byte("<not-rss></not-rss>"), "Example News")
	if err == nil {
		t.Fatal("expected an error for an invalid root element")
	}
}
