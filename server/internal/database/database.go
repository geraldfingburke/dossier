package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// NewDB creates a new database connection
func NewDB() (*sql.DB, error) {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *sql.DB) error {
	schema := `
	-- Drop old user-dependent tables if they exist
	DROP TABLE IF EXISTS delivery_articles CASCADE;
	DROP TABLE IF EXISTS dossier_articles CASCADE;
	DROP TABLE IF EXISTS dossier_deliveries CASCADE;
	DROP TABLE IF EXISTS digest_articles CASCADE;
	DROP TABLE IF EXISTS digest_deliveries CASCADE;
	DROP TABLE IF EXISTS digest_configs CASCADE;
	DROP TABLE IF EXISTS digests CASCADE;
	DROP TABLE IF EXISTS articles CASCADE;
	DROP TABLE IF EXISTS feeds CASCADE;
	DROP TABLE IF EXISTS users CASCADE;

	-- Create new dossier configurations table (replaces users)
	CREATE TABLE IF NOT EXISTS dossier_configs (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		feed_urls TEXT[] NOT NULL, -- Array of feed URLs
		article_count INTEGER DEFAULT 20 CHECK (article_count >= 1 AND article_count <= 50),
		frequency VARCHAR(50) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly')),
		delivery_time TIME NOT NULL, -- Time of day to send
		timezone VARCHAR(50) DEFAULT 'UTC',
		tone VARCHAR(50) DEFAULT 'professional',
		language VARCHAR(50) DEFAULT 'English',
		special_instructions TEXT DEFAULT '',
		active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create feeds table (no user reference)
	CREATE TABLE IF NOT EXISTS feeds (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL UNIQUE,
		title VARCHAR(255),
		description TEXT,
		active BOOLEAN DEFAULT true,
		last_fetched TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create articles table 
	CREATE TABLE IF NOT EXISTS articles (
		id SERIAL PRIMARY KEY,
		feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		link TEXT NOT NULL UNIQUE,
		description TEXT,
		content TEXT,
		author VARCHAR(255),
		published_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create dossier deliveries table (tracks actual sent dossiers)
	CREATE TABLE IF NOT EXISTS dossier_deliveries (
		id SERIAL PRIMARY KEY,
		config_id INTEGER REFERENCES dossier_configs(id) ON DELETE CASCADE,
		delivery_date TIMESTAMP NOT NULL,
		summary TEXT NOT NULL,
		article_count INTEGER NOT NULL,
		email_sent BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create junction table for which articles were in each delivery
	CREATE TABLE IF NOT EXISTS delivery_articles (
		delivery_id INTEGER REFERENCES dossier_deliveries(id) ON DELETE CASCADE,
		article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
		PRIMARY KEY (delivery_id, article_id)
	);

	-- Create indexes for performance
	CREATE INDEX IF NOT EXISTS idx_feeds_url ON feeds(url);
	CREATE INDEX IF NOT EXISTS idx_articles_feed_id ON articles(feed_id);
	CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at);
	CREATE INDEX IF NOT EXISTS idx_dossier_configs_active ON dossier_configs(active);
	CREATE INDEX IF NOT EXISTS idx_dossier_deliveries_config_id ON dossier_deliveries(config_id);
	CREATE INDEX IF NOT EXISTS idx_dossier_deliveries_date ON dossier_deliveries(delivery_date);
	`

	_, err := db.Exec(schema)
	return err
}
