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
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS feeds (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		url TEXT NOT NULL,
		title VARCHAR(255),
		description TEXT,
		active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

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

	CREATE TABLE IF NOT EXISTS digests (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		date DATE NOT NULL,
		summary TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, date)
	);

	CREATE TABLE IF NOT EXISTS digest_articles (
		digest_id INTEGER REFERENCES digests(id) ON DELETE CASCADE,
		article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
		PRIMARY KEY (digest_id, article_id)
	);

	CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id);
	CREATE INDEX IF NOT EXISTS idx_articles_feed_id ON articles(feed_id);
	CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at);
	CREATE INDEX IF NOT EXISTS idx_digests_user_id ON digests(user_id);
	CREATE INDEX IF NOT EXISTS idx_digests_date ON digests(date);
	`

	_, err := db.Exec(schema)
	return err
}
