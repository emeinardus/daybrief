package source

import "context"

// Registry provides access to configured feed sources.
type Registry interface {
	ListActive(ctx context.Context) ([]FeedSource, error)
}
