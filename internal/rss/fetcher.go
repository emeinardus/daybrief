package rss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultRequestTimeout = 10 * time.Second
	maxResponseBytes      = 5 << 20 // 5 MiB
)

// Fetcher retrieves RSS documents over HTTP.
type Fetcher struct {
	client *http.Client
}

// NewFetcher creates a Fetcher with a bounded request timeout.
func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: defaultRequestTimeout,
		},
	}
}

// Fetch downloads an RSS document from the given URL.
func (f *Fetcher) Fetch(ctx context.Context, feedURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create RSS request: %w", err)
	}

	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")

	response, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch RSS feed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch RSS feed: unexpected HTTP status %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read RSS response: %w", err)
	}

	if int64(len(body)) > maxResponseBytes {
		return nil, fmt.Errorf("fetch RSS feed: response exceeds maximum size of %d bytes", maxResponseBytes)
	}

	return body, nil
}
