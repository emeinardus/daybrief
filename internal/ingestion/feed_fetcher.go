package ingestion

import "context"

// FeedFetcher retrieves the raw contents of an RSS feed.
type FeedFetcher interface {
	Fetch(ctx context.Context, feedURL string) ([]byte, error)
}
