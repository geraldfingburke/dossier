// Package models defines the core domain models for the Dossier application.
//
// This package contains all data structures representing the application's domain:
// dossier configurations, RSS feeds, articles, deliveries, and AI tone presets.
//
// # Model Architecture
//
// The models follow a relational design pattern with clear ownership and relationships:
//
//   1. Configuration Layer: DossierConfig (user's automation settings)
//   2. Content Layer: Feed → Article (RSS sources and fetched content)
//   3. Delivery Layer: DossierDelivery ↔ Article (historical deliveries)
//   4. Customization Layer: Tone (AI tone presets)
//
// # Database Mapping
//
// All models use struct tags for:
//   - JSON serialization: `json:"field_name"` (API responses)
//   - Database mapping: `db:"column_name"` (SQL queries with sqlx)
//
// # Key Relationships
//
//   - DossierConfig → DossierDelivery (one-to-many)
//   - Feed → Article (one-to-many)
//   - DossierDelivery ↔ Article (many-to-many via DeliveryArticle junction)
//
// # Type Safety
//
// Custom types enhance type safety:
//   - StringArray: PostgreSQL string array handling with proper scanning
//
// # Timestamp Conventions
//
// Standard timestamp fields across all models:
//   - CreatedAt: Record creation time (immutable)
//   - UpdatedAt: Last modification time (auto-updated)
//   - Additional timestamps: PublishedAt, DeliveryDate, LastFetched (domain-specific)
package models

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

// ============================================================================
// CONFIGURATION MODELS
// ============================================================================

// DossierConfig represents a user's automated news digest configuration.
//
// This is the primary configuration model that defines how, when, and what content
// is aggregated and delivered to users. Each configuration represents one automated
// dossier that runs on a schedule.
//
// Lifecycle:
//   - Created via GraphQL createDossierConfig mutation
//   - Updated via GraphQL updateDossierConfig mutation
//   - Soft-deleted by setting Active = false
//   - Hard-deleted via GraphQL deleteDossierConfig mutation
//
// Field Descriptions:
//   - ID: Unique identifier (auto-generated)
//   - Title: User-friendly name for the dossier (e.g., "Morning Tech News")
//   - Email: Recipient email address for delivery
//   - FeedURLs: Array of RSS feed URLs to aggregate
//   - ArticleCount: Maximum number of articles to include per delivery
//   - Frequency: Delivery schedule - "daily", "weekly", "monthly"
//   - DeliveryTime: Time of day for delivery in HH:MM:SS format
//   - Timezone: IANA timezone for delivery scheduling (e.g., "America/New_York")
//   - Tone: AI tone preset name (references Tone.Name)
//   - Language: Target language for AI summaries (e.g., "English", "Spanish")
//   - SpecialInstructions: Custom AI instructions (optional)
//   - Active: Whether automated delivery is enabled
//   - CreatedAt: Configuration creation timestamp
//   - UpdatedAt: Last modification timestamp
//
// Validation:
//   - Title: Required, non-empty
//   - Email: Required, valid email format
//   - FeedURLs: Required, at least one valid URL
//   - ArticleCount: Required, positive integer (typically 5-20)
//   - Frequency: Required, one of: "daily", "weekly", "monthly"
//   - DeliveryTime: Required, valid time in HH:MM:SS format
//   - Timezone: Required, valid IANA timezone
//   - Tone: Defaults to "professional" if not specified
//   - Language: Defaults to "English" if not specified
//
// Example:
//   config := models.DossierConfig{
//       Title:        "Daily AI News",
//       Email:        "user@example.com",
//       FeedURLs:     []string{"https://news.ycombinator.com/rss"},
//       ArticleCount: 10,
//       Frequency:    "daily",
//       DeliveryTime: "08:00:00",
//       Timezone:     "America/New_York",
//       Tone:         "professional",
//       Language:     "English",
//       Active:       true,
//   }
type DossierConfig struct {
	ID                  int       `json:"id" db:"id"`
	Title               string    `json:"title" db:"title"`
	Email               string    `json:"email" db:"email"`
	FeedURLs            []string  `json:"feed_urls" db:"feed_urls"`
	ArticleCount        int       `json:"article_count" db:"article_count"`
	Frequency           string    `json:"frequency" db:"frequency"`
	DeliveryTime        string    `json:"delivery_time" db:"delivery_time"`
	Timezone            string    `json:"timezone" db:"timezone"`
	Tone                string    `json:"tone" db:"tone"`
	Language            string    `json:"language" db:"language"`
	SpecialInstructions string    `json:"special_instructions" db:"special_instructions"`
	Active              bool      `json:"active" db:"active"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// DATABASE TYPE HELPERS
// ============================================================================

// StringArray is a custom type for PostgreSQL string array columns.
//
// PostgreSQL uses a different representation for arrays than Go's native []string,
// requiring custom serialization and deserialization logic. This type implements
// the driver.Valuer and sql.Scanner interfaces to handle the conversion.
//
// Usage:
// This type is used internally by the database layer but is transparent to
// application code. Models use []string directly, and the database driver
// handles the conversion automatically.
//
// Implementation:
//   - Delegates to github.com/lib/pq for actual PostgreSQL array handling
//   - Handles empty arrays correctly (returns "{}" instead of null)
//
// Example SQL:
//   CREATE TABLE example (
//       tags TEXT[]  -- PostgreSQL array column
//   );
//
// Example Go:
//   type Model struct {
//       Tags []string `db:"tags"`  -- Automatically converted via StringArray
//   }
type StringArray []string

// Value implements the driver.Valuer interface for database storage.
//
// This method converts a Go []string to a PostgreSQL-compatible array format.
// Called automatically when inserting or updating records.
//
// Behavior:
//   - Empty slice → "{}" (empty PostgreSQL array)
//   - Non-empty → Delegates to pq.Array for proper escaping
//
// Returns:
//   - driver.Value: PostgreSQL array representation
//   - error: Conversion error (rare, usually nil)
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return pq.Array(a).Value()
}

// Scan implements the sql.Scanner interface for database retrieval.
//
// This method converts a PostgreSQL array column to a Go []string.
// Called automatically when querying records.
//
// Parameters:
//   - value: Raw database value (PostgreSQL array)
//
// Returns:
//   - error: Conversion error if value is not a valid array
func (a *StringArray) Scan(value interface{}) error {
	return pq.Array(a).Scan(value)
}

// ============================================================================
// CONTENT SOURCE MODELS
// ============================================================================

// Feed represents an RSS feed source for article aggregation.
//
// Feeds are the content sources for dossiers. Each feed represents one RSS/Atom
// feed URL from which articles are periodically fetched.
//
// Note: In the current architecture, feeds are referenced by URL in DossierConfig.FeedURLs
// rather than by foreign key. This table tracks feed metadata and fetch history.
//
// Field Descriptions:
//   - ID: Unique identifier
//   - URL: RSS/Atom feed URL (must be unique)
//   - Title: Feed title (extracted from RSS metadata)
//   - Description: Feed description (from RSS metadata)
//   - Active: Whether this feed is available for use
//   - LastFetched: Timestamp of most recent successful fetch
//   - CreatedAt: Feed registration timestamp
//   - UpdatedAt: Last modification timestamp
//
// Lifecycle:
//   - Auto-created when first referenced in a DossierConfig
//   - Updated on each successful fetch
//   - Marked inactive if feed becomes unavailable
//
// Example:
//   feed := models.Feed{
//       URL:         "https://news.ycombinator.com/rss",
//       Title:       "Hacker News",
//       Description: "Links for the intellectually curious",
//       Active:      true,
//       LastFetched: time.Now(),
//   }
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

// ============================================================================
// CONTENT MODELS
// ============================================================================

// Article represents a single article fetched from an RSS feed.
//
// Articles are the core content units that get aggregated into dossiers.
// Each article is parsed from RSS/Atom XML and stored for processing by
// the AI summarization service.
//
// Field Descriptions:
//   - ID: Unique identifier
//   - FeedID: Reference to source Feed (0 if not tracked)
//   - Title: Article headline
//   - Link: Canonical URL to full article
//   - Description: Article excerpt or summary (from RSS)
//   - Content: Full article text (if available in feed)
//   - Author: Article author name
//   - PublishedAt: Original publication timestamp from feed
//   - CreatedAt: When article was fetched and stored
//
// Data Quality:
//   - Title: Always present (required by RSS spec)
//   - Link: Always present (required by RSS spec)
//   - Description: Usually present, may be empty
//   - Content: Often empty (many feeds only provide excerpts)
//   - Author: Often empty or generic
//   - PublishedAt: Usually accurate, but may be missing or incorrect
//
// Usage in Pipeline:
//   1. Fetched from RSS feeds by rss.Service
//   2. Selected by AI for inclusion in summary
//   3. Summarized by ai.Service
//   4. Included in email by email.Service
//   5. Linked to delivery via DeliveryArticle
//
// Example:
//   article := models.Article{
//       FeedID:      1,
//       Title:       "Go 1.22 Released",
//       Link:        "https://go.dev/blog/go1.22",
//       Description: "The Go team is happy to announce...",
//       Author:      "Go Team",
//       PublishedAt: time.Now().Add(-2 * time.Hour),
//   }
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

// ============================================================================
// DELIVERY MODELS
// ============================================================================

// DossierDelivery represents a historical dossier that was generated and sent.
//
// This model provides an audit trail of all dossier deliveries, storing the
// AI-generated summary and metadata about what was sent. It forms a many-to-many
// relationship with Articles through the DeliveryArticle junction table.
//
// Field Descriptions:
//   - ID: Unique delivery identifier
//   - ConfigID: Reference to DossierConfig that triggered this delivery
//   - DeliveryDate: When the dossier was generated and sent
//   - Summary: AI-generated HTML summary of articles
//   - ArticleCount: Number of articles included
//   - EmailSent: Whether email was successfully delivered
//   - Articles: Populated list of articles (via SQL join, not in DB)
//   - CreatedAt: Record creation timestamp
//
// Lifecycle:
//   1. Created after successful dossier generation
//   2. Updated if email delivery fails initially
//   3. Never deleted (permanent audit trail)
//
// Query Patterns:
//   - List deliveries for specific config: WHERE config_id = ?
//   - Recent deliveries: ORDER BY delivery_date DESC LIMIT 10
//   - Failed deliveries: WHERE email_sent = false
//
// Example:
//   delivery := models.DossierDelivery{
//       ConfigID:     1,
//       DeliveryDate: time.Now(),
//       Summary:      "<h2>Today's Top Stories</h2>...",
//       ArticleCount: 10,
//       EmailSent:    true,
//   }
type DossierDelivery struct {
	ID           int       `json:"id" db:"id"`
	ConfigID     int       `json:"config_id" db:"config_id"`
	DeliveryDate time.Time `json:"delivery_date" db:"delivery_date"`
	Summary      string    `json:"summary" db:"summary"`
	ArticleCount int       `json:"article_count" db:"article_count"`
	EmailSent    bool      `json:"email_sent" db:"email_sent"`
	Articles     []Article `json:"articles"` // Populated via join, not stored in this table
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DeliveryArticle is a junction table linking deliveries to articles.
//
// This model implements the many-to-many relationship between DossierDelivery
// and Article, allowing:
//   - One delivery to contain multiple articles
//   - One article to appear in multiple deliveries
//
// Database Structure:
//   - Composite primary key: (delivery_id, article_id)
//   - Foreign keys: delivery_id → dossier_deliveries.id
//                   article_id → articles.id
//
// Usage:
// This is typically managed automatically when recording deliveries and is
// rarely accessed directly in application code. The Articles field in
// DossierDelivery is populated via SQL joins.
//
// Example SQL:
//   SELECT a.* FROM articles a
//   JOIN delivery_articles da ON a.id = da.article_id
//   WHERE da.delivery_id = $1
type DeliveryArticle struct {
	DeliveryID int `json:"delivery_id" db:"delivery_id"`
	ArticleID  int `json:"article_id" db:"article_id"`
}

// ============================================================================
// AI CUSTOMIZATION MODELS
// ============================================================================

// Tone represents an AI tone preset for summary generation.
//
// Tones control the style, voice, and presentation of AI-generated summaries.
// The system includes default tones (protected) and allows users to create
// custom tones for specific use cases.
//
// Field Descriptions:
//   - ID: Unique tone identifier
//   - Name: Display name (e.g., "Professional", "Casual", "Pirate")
//   - Prompt: AI instruction text that defines the tone style
//   - IsSystemDefault: Whether this is a protected system tone
//   - CreatedAt: Tone creation timestamp
//   - UpdatedAt: Last modification timestamp
//
// System Default Tones:
// The database migration creates 10 default tones:
//   - Professional: Business-appropriate, clear, objective
//   - Casual: Friendly, conversational, accessible
//   - Academic: Scholarly, formal, precise
//   - Creative: Engaging, narrative, colorful
//   - Technical: Developer-focused, detailed, accurate
//   - Concise: Brief, bullet-points, essential facts only
//   - Detailed: Comprehensive, thorough, in-depth
//   - Humorous: Light-hearted, witty, entertaining
//   - Serious: Factual, somber, no-nonsense
//   - Pirate: Ahoy matey! (Yes, really. It's fun.)
//
// Protection Rules:
//   - System default tones cannot be updated via GraphQL
//   - System default tones cannot be deleted via GraphQL
//   - Custom tones (IsSystemDefault = false) can be modified/deleted
//
// Custom Tone Creation:
// Users can create custom tones for specific needs:
//   - Organization-specific voice (e.g., "Company Newsletter")
//   - Audience-specific (e.g., "Executive Summary", "Engineering Team")
//   - Experimental variations (e.g., "Sarcastic", "ELI5")
//
// Example:
//   tone := models.Tone{
//       Name:            "Executive Summary",
//       Prompt:          "Present information as a brief executive summary...",
//       IsSystemDefault: false,
//   }
type Tone struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Prompt          string    `json:"prompt" db:"prompt"`
	IsSystemDefault bool      `json:"is_system_default" db:"is_system_default"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
