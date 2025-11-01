package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Feed represents an RSS feed
type Feed struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	URL         string    `json:"url" db:"url"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
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

// Digest represents a daily digest with AI summaries
type Digest struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Date      time.Time `json:"date" db:"date"`
	Summary   string    `json:"summary" db:"summary"`
	Articles  []Article `json:"articles"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DigestArticle is a junction table linking digests to articles
type DigestArticle struct {
	DigestID  int `json:"digest_id" db:"digest_id"`
	ArticleID int `json:"article_id" db:"article_id"`
}
