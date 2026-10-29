package ingestion

import (
	"context"
	"fmt"

	"github.com/emeinardus/daybrief/internal/article"
	"github.com/emeinardus/daybrief/internal/rss"
)

// Service coordinates fetching and parsing RSS feeds.
type Service struct {
	fetcher FeedFetcher
}

// NewService creates an ingestion service.
func NewService(fetcher FeedFetcher) *Service {
	return &Service{
		fetcher: fetcher,
	}
}

// LoadFeed retrieves an RSS feed and converts it into articles.
func (s *Service) LoadFeed(
	ctx context.Context,
	feedURL string,
	source string,
) ([]article.Article, error) {
	data, err := s.fetcher.Fetch(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("load RSS feed: %w", err)
	}

	articles, err := rss.Parse(data, source)
	if err != nil {
		return nil, fmt.Errorf("load RSS feed: %w", err)
	}

	return articles, nil
}
