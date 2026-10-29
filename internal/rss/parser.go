package rss

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/emeinardus/daybrief/internal/article"
)

// Parse converts an RSS 2.0 document into Daybrief articles.
func Parse(data []byte, source string) ([]article.Article, error) {
	var document rssDocument

	if err := xml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("parse RSS document: %w", err)
	}

	if document.XMLName.Local != "rss" {
		return nil, fmt.Errorf("invalid RSS document: expected root element <rss>")
	}

	articles := make([]article.Article, 0, len(document.Channel.Items))

	for _, item := range document.Channel.Items {
		var publishedAt *time.Time

		if value := strings.TrimSpace(item.PubDate); value != "" {
			parsed, err := parsePubDate(value)
			if err == nil {
				publishedAt = &parsed
			}
		}

		articles = append(articles, article.Article{
			Title:       strings.TrimSpace(item.Title),
			URL:         strings.TrimSpace(item.Link),
			Description: strings.TrimSpace(item.Description),
			PublishedAt: publishedAt,
			Source:      source,
		})
	}

	return articles, nil
}

// parsePubDate parses common date formats used by RSS feeds.
func parsePubDate(value string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, format := range formats {
		parsed, err := time.Parse(format, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported RSS publication date %q", value)
}
