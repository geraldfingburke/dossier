// Package database provides PostgreSQL database connection management and schema migrations
// for the Dossier application. It handles database initialization, connection pooling,
// and versioned schema management for all core tables (dossier configurations, feeds,
// articles, deliveries, and AI tones).
package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// defaultDatabaseURL is the fallback connection string when DATABASE_URL is not set
	// Format: postgres://username:password@host:port/database?sslmode=disable
	defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable"
)

// ============================================================================
// CONNECTION MANAGEMENT
// ============================================================================

// NewDB establishes a new PostgreSQL database connection with the following behavior:
//
// Connection Source:
//   - Reads DATABASE_URL environment variable if set
//   - Falls back to localhost default if not set
//
// Connection Verification:
//   - Opens connection pool
//   - Verifies connectivity with Ping()
//   - Returns error if connection fails
//
// Connection Pooling:
//   - Uses sql.DB connection pool (managed by database/sql)
//   - Pool size configurable via PostgreSQL driver parameters
//
// Returns:
//   - *sql.DB: Active database connection pool
//   - error: Connection or ping failure
//
// Example:
//   db, err := NewDB()
//   if err != nil {
//       log.Fatal("Database connection failed:", err)
//   }
//   defer db.Close()
func NewDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = defaultDatabaseURL
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

// ============================================================================
// SCHEMA MIGRATION
// ============================================================================

// Migrate executes database schema migrations to set up or update the database structure.
//
// Migration Strategy:
//   - Idempotent: Safe to run multiple times
//   - Uses CREATE TABLE IF NOT EXISTS for incremental migrations
//   - Drops legacy tables from previous schema versions
//   - Inserts default data (tones) only if not already present
//
// Schema Components:
//
// 1. Core Tables:
//   - dossier_configs: User dossier configurations and preferences
//   - feeds: RSS feed sources
//   - articles: Fetched articles from feeds
//   - dossier_deliveries: Delivery history and content
//   - delivery_articles: Many-to-many relationship (deliveries ↔ articles)
//   - tones: AI writing styles for summaries
//
// 2. Performance Indexes:
//   - Feed URL lookup
//   - Article queries by feed and date
//   - Active dossier filtering
//   - Delivery history queries
//   - Tone name lookup
//
// 3. Default Data:
//   - 10 predefined AI tones (professional, humorous, apocalyptic, etc.)
//
// Migration Order:
//   1. Drop legacy tables (old schema cleanup)
//   2. Create tables in dependency order
//   3. Create performance indexes
//   4. Insert default tone configurations
//
// Parameters:
//   - db: Active database connection
//
// Returns:
//   - error: Any SQL execution error
//
// Example:
//   if err := Migrate(db); err != nil {
//       log.Fatal("Migration failed:", err)
//   }
func Migrate(db *sql.DB) error {
	schema := `
	-- ========================================================================
	-- CLEANUP: Drop legacy tables from previous schema versions
	-- ========================================================================
	-- These tables are from earlier iterations and are no longer used
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

	-- ========================================================================
	-- TABLE: dossier_configs
	-- ========================================================================
	-- Stores user dossier configurations (replaces old user-centric design)
	-- Each record represents one configured dossier with delivery preferences
	--
	-- Key Fields:
	--   - feed_urls: Array of RSS feed URLs to monitor
	--   - frequency: How often to deliver (daily/weekly/monthly)
	--   - delivery_time: Time of day to send (in specified timezone)
	--   - tone: AI writing style (references tones table)
	--   - active: Enable/disable delivery without deletion
	-- ========================================================================
	CREATE TABLE IF NOT EXISTS dossier_configs (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		feed_urls TEXT[] NOT NULL,
		article_count INTEGER DEFAULT 20 CHECK (article_count >= 1 AND article_count <= 50),
		frequency VARCHAR(50) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly')),
		delivery_time TIME NOT NULL,
		timezone VARCHAR(50) DEFAULT 'UTC',
		tone VARCHAR(50) DEFAULT 'professional',
		language VARCHAR(50) DEFAULT 'English',
		special_instructions TEXT DEFAULT '',
		active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- ========================================================================
	-- TABLE: feeds
	-- ========================================================================
	-- Tracks RSS/Atom feed sources
	-- Decoupled from users - feeds are shared across all dossiers
	--
	-- Key Fields:
	--   - url: Unique feed URL
	--   - active: Can temporarily disable problematic feeds
	--   - last_fetched: Track fetch schedule and detect stale feeds
	-- ========================================================================
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

	-- ========================================================================
	-- TABLE: articles
	-- ========================================================================
	-- Stores fetched articles from RSS feeds
	-- Articles are shared across dossiers - deduplication by link
	--
	-- Key Fields:
	--   - link: Unique article URL (prevents duplicates)
	--   - published_at: Original publication date (for sorting/filtering)
	--   - content: Full article content if available
	--   - description: Article summary/excerpt
	-- ========================================================================
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

	-- ========================================================================
	-- TABLE: dossier_deliveries
	-- ========================================================================
	-- Tracks actual sent dossier emails and their content
	-- Maintains delivery history and enables archive viewing
	--
	-- Key Fields:
	--   - config_id: Links to dossier configuration
	--   - summary: Generated AI summary (stored for archive)
	--   - email_sent: Delivery status tracking
	--   - delivery_date: When the dossier was sent
	-- ========================================================================
	CREATE TABLE IF NOT EXISTS dossier_deliveries (
		id SERIAL PRIMARY KEY,
		config_id INTEGER REFERENCES dossier_configs(id) ON DELETE CASCADE,
		delivery_date TIMESTAMP NOT NULL,
		summary TEXT NOT NULL,
		article_count INTEGER NOT NULL,
		email_sent BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- ========================================================================
	-- TABLE: delivery_articles
	-- ========================================================================
	-- Junction table: Many-to-many relationship between deliveries and articles
	-- Tracks which specific articles were included in each delivery
	-- Enables:
	--   - Article-level delivery tracking
	--   - "Already sent" detection for deduplication
	--   - Historical article viewing in archive
	-- ========================================================================
	CREATE TABLE IF NOT EXISTS delivery_articles (
		delivery_id INTEGER REFERENCES dossier_deliveries(id) ON DELETE CASCADE,
		article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
		PRIMARY KEY (delivery_id, article_id)
	);

	-- ========================================================================
	-- TABLE: tones
	-- ========================================================================
	-- Customizable AI writing styles for summary generation
	-- Each tone contains a prompt that guides the AI's writing style
	--
	-- Key Fields:
	--   - name: Unique tone identifier (e.g., "professional", "humorous")
	--   - prompt: Detailed instructions for the AI model
	--   - is_system_default: System tones vs user-created custom tones
	-- ========================================================================
	CREATE TABLE IF NOT EXISTS tones (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE,
		prompt TEXT NOT NULL,
		is_system_default BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- ========================================================================
	-- PERFORMANCE INDEXES
	-- ========================================================================
	-- Critical indexes for query performance on large datasets
	
	-- Feed lookup optimization
	CREATE INDEX IF NOT EXISTS idx_feeds_url ON feeds(url);
	
	-- Article queries by feed relationship
	CREATE INDEX IF NOT EXISTS idx_articles_feed_id ON articles(feed_id);
	
	-- Article sorting and time-range filtering
	CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at);
	
	-- Active dossier filtering for scheduler
	CREATE INDEX IF NOT EXISTS idx_dossier_configs_active ON dossier_configs(active);
	
	-- Delivery history queries by dossier
	CREATE INDEX IF NOT EXISTS idx_dossier_deliveries_config_id ON dossier_deliveries(config_id);
	
	-- Delivery chronological queries and archive pagination
	CREATE INDEX IF NOT EXISTS idx_dossier_deliveries_date ON dossier_deliveries(delivery_date);
	
	-- Tone lookup by name (most common query pattern)
	CREATE INDEX IF NOT EXISTS idx_tones_name ON tones(name);

	-- ========================================================================
	-- DEFAULT DATA: System Tones
	-- ========================================================================
	-- Pre-configured AI writing styles available to all users
	-- Uses ON CONFLICT to make this idempotent (safe to run multiple times)
	--
	-- Tone Categories:
	--   - Professional/Business: professional, analytical
	--   - Creative/Fun: humorous, casual, apocalyptic, orc, robot
	--   - Regional/Cultural: southern_belle
	--   - Experimental: apologetic, sweary (adult content)
	-- ========================================================================
	INSERT INTO tones (name, prompt, is_system_default) VALUES 
		('professional', 'Write in a professional, formal tone suitable for business communication. Be clear, concise, and authoritative.', true),
		('humorous', 'Write with humor and wit. Use light-hearted commentary, clever observations, and entertaining language while maintaining informative value.', true),
		('analytical', 'Focus on data-driven insights, trends, and deep analysis. Use precise language and highlight statistical significance and implications.', true),
		('casual', 'Write in a friendly, conversational tone as if talking to a colleague. Be approachable and easy to understand.', true),
		('apocalyptic', 'Frame everything as if the world is ending. Use dramatic, urgent language and treat every piece of news as a harbinger of doom.', true),
		('orc', 'Write like a fantasy orc warrior. Use rough, aggressive language with lots of grunts and battle metaphors. WAAAAAGH!', true),
		('robot', 'BEEP BOOP. PROCESSING INFORMATION. USE ROBOTIC LANGUAGE WITH TECHNICAL PRECISION. ELIMINATE EMOTIONAL RESPONSES.', true),
		('southern_belle', 'Write with Southern charm and hospitality. Use sweet, polite language with a touch of sass and regional expressions, darlin''.', true),
		('apologetic', 'Apologize for everything. Feel sorry about all the news being reported. Use hesitant, self-deprecating language.', true),
		('sweary', 'Use uncensored, explicit language. Don''t hold back on profanity when expressing opinions about the news. Adult content warning.', true)
	ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("migration execution failed: %w", err)
	}

	return nil
}
