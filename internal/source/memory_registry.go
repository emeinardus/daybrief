package source

import (
	"context"
	"fmt"
)

// MemoryRegistry stores feed sources in memory.
type MemoryRegistry struct {
	sources []FeedSource
}

// NewMemoryRegistry creates a registry backed by an in-memory source list.
func NewMemoryRegistry(sources []FeedSource) *MemoryRegistry {
	copiedSources := append([]FeedSource(nil), sources...)

	return &MemoryRegistry{
		sources: copiedSources,
	}
}

// ListActive returns all active feed sources.
func (r *MemoryRegistry) ListActive(ctx context.Context) ([]FeedSource, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("list active feed sources: %w", err)
	}

	activeSources := make([]FeedSource, 0, len(r.sources))

	for _, feedSource := range r.sources {
		if feedSource.Active {
			activeSources = append(activeSources, feedSource)
		}
	}

	return activeSources, nil
}
