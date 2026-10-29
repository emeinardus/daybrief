package ingestion

import (
	"context"
	"errors"
	"testing"

	"github.com/emeinardus/daybrief/internal/rss"
)

type stubFetcher struct {
	data []byte
	err  error
}

func (f *stubFetcher) Fetch(
	_ context.Context,
	_ string,
) ([]byte, error) {
	return f.data, f.err
}

func TestLoadFeed(t *testing.T) {
	input := []byte(`<rss version="2.0">
		<channel>
			<title>Example News</title>
			<item>
				<title>Technology Update</title>
				<link>https://example.com/article</link>
			</item>
		</channel>
	</rss>`)

	service := NewService(&stubFetcher{data: input})

	articles, err := service.LoadFeed(
		context.Background(),
		"https://example.com/feed.xml",
		"Example News",
	)
	if err != nil {
		t.Fatalf("LoadFeed() returned an unexpected error: %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(articles))
	}

	if articles[0].Title != "Technology Update" {
		t.Errorf("unexpected article title: %q", articles[0].Title)
	}

	if articles[0].Source != "Example News" {
		t.Errorf("unexpected article source: %q", articles[0].Source)
	}
}

func TestLoadFeedFetchError(t *testing.T) {
	expectedErr := errors.New("network unavailable")
	service := NewService(&stubFetcher{err: expectedErr})

	_, err := service.LoadFeed(
		context.Background(),
		"https://example.com/feed.xml",
		"Example News",
	)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped fetch error, got %v", err)
	}
}

func TestLoadFeedInvalidXML(t *testing.T) {
	service := NewService(&stubFetcher{
		data: []byte("<rss><channel>"),
	})

	_, err := service.LoadFeed(
		context.Background(),
		"https://example.com/feed.xml",
		"Example News",
	)
	if err == nil {
		t.Fatal("expected an error for invalid RSS XML")
	}
}

// Ensure the production fetcher satisfies the application interface.
var _ FeedFetcher = (*rss.Fetcher)(nil)
