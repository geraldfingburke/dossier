package rss

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/geraldfingburke/dossier/server/internal/ai"
	"github.com/geraldfingburke/dossier/server/internal/models"
	"github.com/mmcdole/gofeed"
)

// Service handles RSS feed operations
type Service struct {
	parser    *gofeed.Parser
	aiService *ai.Service
}

// NewService creates a new RSS service
func NewService(aiService *ai.Service) *Service {
	return &Service{
		parser:    gofeed.NewParser(),
		aiService: aiService,
	}
}

// FetchFeed fetches and parses an RSS feed
func (s *Service) FetchFeed(ctx context.Context, feedURL string) (*gofeed.Feed, error) {
	feed, err := s.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("error parsing feed: %w", err)
	}
	return feed, nil
}

// SaveArticles saves articles from a feed to the database
func (s *Service) SaveArticles(db *sql.DB, feedID int, items []*gofeed.Item) error {
	for _, item := range items {
		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		author := ""
		if item.Author != nil {
			author = item.Author.Name
		}

		// Insert article if it doesn't exist
		_, err := db.Exec(`
			INSERT INTO articles (feed_id, title, link, description, content, author, published_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (link) DO NOTHING
		`, feedID, item.Title, item.Link, item.Description, content, author, publishedAt)

		if err != nil {
			return fmt.Errorf("error saving article: %w", err)
		}
	}
	return nil
}

// FetchAllFeeds fetches all active feeds for a user
func (s *Service) FetchAllFeeds(ctx context.Context, db *sql.DB, userID int) error {
	rows, err := db.QueryContext(ctx, `
		SELECT id, url FROM feeds WHERE user_id = $1 AND active = true
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var feedID int
		var url string
		if err := rows.Scan(&feedID, &url); err != nil {
			log.Printf("Error scanning feed: %v", err)
			continue
		}

		feed, err := s.FetchFeed(ctx, url)
		if err != nil {
			log.Printf("Error fetching feed %s: %v", url, err)
			continue
		}

		// Update feed title and description
		_, err = db.ExecContext(ctx, `
			UPDATE feeds SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $3
		`, feed.Title, feed.Description, feedID)
		if err != nil {
			log.Printf("Error updating feed: %v", err)
		}

		// Save articles
		if err := s.SaveArticles(db, feedID, feed.Items); err != nil {
			log.Printf("Error saving articles: %v", err)
		}
	}

	return rows.Err()
}

// GenerateDailyDigests generates AI summaries for all users
func (s *Service) GenerateDailyDigests(ctx context.Context, db *sql.DB) error {
	// Get all users
	rows, err := db.QueryContext(ctx, "SELECT id FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}

		if err := s.GenerateUserDigest(ctx, db, userID); err != nil {
			log.Printf("Error generating digest for user %d: %v", userID, err)
		}
	}

	return rows.Err()
}

// GenerateUserDigest generates a digest for a specific user
func (s *Service) GenerateUserDigest(ctx context.Context, db *sql.DB, userID int) error {
	// First, fetch latest articles
	if err := s.FetchAllFeeds(ctx, db, userID); err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}

	// Get articles from the last 24 hours
	yesterday := time.Now().Add(-24 * time.Hour)
	rows, err := db.QueryContext(ctx, `
		SELECT a.id, a.title, a.description, a.content, a.link, a.author, a.published_at
		FROM articles a
		JOIN feeds f ON a.feed_id = f.id
		WHERE f.user_id = $1 AND a.published_at > $2
		ORDER BY a.published_at DESC
	`, userID, yesterday)
	if err != nil {
		return err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Description, &article.Content, &article.Link, &article.Author, &article.PublishedAt); err != nil {
			return err
		}
		articles = append(articles, article)
	}

	if len(articles) == 0 {
		log.Printf("No new articles for user %d", userID)
		return nil
	}

	// Generate AI summary
	summary, err := s.aiService.SummarizeArticles(ctx, articles)
	if err != nil {
		return fmt.Errorf("error generating summary: %w", err)
	}

	// Create digest
	today := time.Now().Truncate(24 * time.Hour)
	var digestID int
	err = db.QueryRowContext(ctx, `
		INSERT INTO digests (user_id, date, summary)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, date) DO UPDATE SET summary = $3
		RETURNING id
	`, userID, today, summary).Scan(&digestID)
	if err != nil {
		return err
	}

	// Link articles to digest
	for _, article := range articles {
		_, err = db.ExecContext(ctx, `
			INSERT INTO digest_articles (digest_id, article_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, digestID, article.ID)
		if err != nil {
			log.Printf("Error linking article to digest: %v", err)
		}
	}

	log.Printf("Generated digest for user %d with %d articles", userID, len(articles))
	return nil
}
