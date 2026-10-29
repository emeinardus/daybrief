package article

import "time"

// Article represents a news article normalized into Daybrief's domain model.
type Article struct {
	Title       string
	URL         string
	Description string
	PublishedAt *time.Time
	Source      string
}
