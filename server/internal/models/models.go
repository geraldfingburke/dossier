package models

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

// DossierConfig represents a dossier configuration that defines how and when dossiers are generated
type DossierConfig struct {
	ID                  int       `json:"id" db:"id"`
	Title               string    `json:"title" db:"title"`
	Email               string    `json:"email" db:"email"`
	FeedURLs            []string  `json:"feed_urls" db:"feed_urls"`
	ArticleCount        int       `json:"article_count" db:"article_count"`
	Frequency           string    `json:"frequency" db:"frequency"` // daily, weekly, monthly
	DeliveryTime        string    `json:"delivery_time" db:"delivery_time"` // HH:MM format
	Timezone            string    `json:"timezone" db:"timezone"`
	Tone                string    `json:"tone" db:"tone"`
	Language            string    `json:"language" db:"language"`
	SpecialInstructions string    `json:"special_instructions" db:"special_instructions"`
	Active              bool      `json:"active" db:"active"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// StringArray is a custom type for PostgreSQL string arrays
type StringArray []string

// Value implements the driver.Valuer interface for database storage
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return pq.Array(a).Value()
}

// Scan implements the sql.Scanner interface for database retrieval
func (a *StringArray) Scan(value interface{}) error {
	return pq.Array(a).Scan(value)
}

// Feed represents an RSS feed (no user association)
type Feed struct {
	ID          int       `json:"id" db:"id"`
	URL         string    `json:"url" db:"url"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	LastFetched time.Time `json:"last_fetched" db:"last_fetched"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Article represents a fetched article from an RSS feed
type Article struct {
	ID          int       `json:"id" db:"id"`
	FeedID      int       `json:"feed_id" db:"feed_id"`
	Title       string    `json:"title" db:"title"`
	Link        string    `json:"link" db:"link"`
	Description string    `json:"description" db:"description"`
	Content     string    `json:"content" db:"content"`
	Author      string    `json:"author" db:"author"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// DossierDelivery represents an actual delivered dossier
type DossierDelivery struct {
	ID           int       `json:"id" db:"id"`
	ConfigID     int       `json:"config_id" db:"config_id"`
	DeliveryDate time.Time `json:"delivery_date" db:"delivery_date"`
	Summary      string    `json:"summary" db:"summary"`
	ArticleCount int       `json:"article_count" db:"article_count"`
	EmailSent    bool      `json:"email_sent" db:"email_sent"`
	Articles     []Article `json:"articles"` // Populated via join
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DeliveryArticle is a junction table linking deliveries to articles
type DeliveryArticle struct {
	DeliveryID int `json:"delivery_id" db:"delivery_id"`
	ArticleID  int `json:"article_id" db:"article_id"`
}
