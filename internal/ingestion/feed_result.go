package ingestion

import "github.com/emeinardus/daybrief/internal/source"

// FeedResult describes the outcome of processing one feed source.
type FeedResult struct {
	Source source.FeedSource
	Count  int
	Error  error
}
