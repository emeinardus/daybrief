package ingestion

import (
	"context"
	"fmt"

	"github.com/emeinardus/daybrief/internal/article"
	"github.com/emeinardus/daybrief/internal/rss"
	"github.com/emeinardus/daybrief/internal/source"
)

// Service coordinates fetching and parsing RSS feeds.
type Service struct {
	fetcher  FeedFetcher
	registry source.Registry
}

// NewService creates an ingestion service.
func NewService(fetcher FeedFetcher, registry source.Registry) *Service {
	return &Service{
		fetcher:  fetcher,
		registry: registry,
	}
}

// LoadFeed retrieves an RSS feed and converts it into articles.
func (s *Service) LoadFeed(
	ctx context.Context,
	feedURL string,
	sourceName string,
) ([]article.Article, error) {
	data, err := s.fetcher.Fetch(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("load RSS feed: %w", err)
	}

	articles, err := rss.Parse(data, sourceName)
	if err != nil {
		return nil, fmt.Errorf("load RSS feed: %w", err)
	}

	return articles, nil
}

// LoadActiveFeeds processes every active source and collects successful articles.
func (s *Service) LoadActiveFeeds(
	ctx context.Context,
) ([]article.Article, []FeedResult, error) {
	sources, err := s.registry.ListActive(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load active feed sources: %w", err)
	}

	articles := make([]article.Article, 0)
	results := make([]FeedResult, 0, len(sources))

	for _, feedSource := range sources {
		if err := ctx.Err(); err != nil {
			return articles, results, fmt.Errorf("load active feeds: %w", err)
		}

		feedArticles, err := s.LoadFeed(ctx, feedSource.URL, feedSource.Name)
		if err != nil {
			results = append(results, FeedResult{
				Source: feedSource,
				Error:  err,
			})
			continue
		}

		articles = append(articles, feedArticles...)
		results = append(results, FeedResult{
			Source: feedSource,
			Count:  len(feedArticles),
		})
	}

	return articles, results, nil
}
