// Package scheduler provides automated dossier generation and delivery scheduling.
//
// This package implements a time-based scheduling system that periodically checks
// for dossier configurations that need processing and orchestrates the complete
// delivery pipeline: fetching articles, generating AI summaries, and sending emails.
//
// # Architecture
//
// The scheduler uses a ticker-based approach with 1-minute granularity:
//   1. Ticker fires every minute
//   2. Query active dossier configurations
//   3. Check each configuration's schedule
//   4. Generate and send dossiers asynchronously
//   5. Record delivery history
//
// # Scheduling Logic
//
// The scheduler respects user-configured timezones and delivery schedules:
//
// Frequency Support:
//   - Daily: Delivers once per day at specified time
//   - Weekly: Delivers on Mondays at specified time
//   - Monthly: Delivers on 1st of month at specified time
//
// Timezone Handling:
//   - Each configuration has its own timezone (IANA format)
//   - Delivery time is evaluated in the configuration's timezone
//   - Falls back to system time if timezone is invalid
//
// Duplicate Prevention:
//   - Tracks last delivery time per configuration
//   - Prevents multiple deliveries within same period
//   - Uses dossier_deliveries table as delivery log
//
// # Concurrency Model
//
// The scheduler is designed for safe concurrent operation:
//   - Single ticker goroutine checks schedules
//   - Each dossier generation runs in separate goroutine
//   - Thread-safe start/stop via mutex
//   - Graceful shutdown via stop channel
//
// # Error Handling Philosophy
//
// The scheduler is resilient to individual failures:
//   - Configuration parsing errors → Skip, continue with others
//   - Feed fetching errors → Partial results acceptable
//   - Email failures → Logged, don't block other deliveries
//   - Database errors → Logged, scheduler continues
//
// # Integration Points
//
// The scheduler orchestrates four services:
//   - RSS Service: Fetches articles from configured feeds
//   - AI Service: Generates summaries with specified tone
//   - Email Service: Delivers formatted emails
//   - Database: Tracks configurations and delivery history
//
// # Lifecycle
//
//   1. NewService(): Initialize with service dependencies
//   2. Start(): Begin ticker loop in goroutine
//   3. [Runtime]: Automatic dossier processing
//   4. Stop(): Graceful shutdown, stop ticker
//
// # Performance Characteristics
//
//   - Check frequency: 1 minute (configurable via ticker duration)
//   - Dossier generation: Async (doesn't block other deliveries)
//   - Context timeout: 10 minutes per dossier
//   - Database queries: Minimal (one query per check cycle)
//
// # Usage Example
//
//   scheduler := scheduler.NewService(db, rssService, aiService, emailService)
//   scheduler.Start()
//   defer scheduler.Stop()
//   // Scheduler runs in background until Stop() is called
package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/geraldfingburke/dossier/server/internal/ai"
	"github.com/geraldfingburke/dossier/server/internal/email"
	"github.com/geraldfingburke/dossier/server/internal/models"
	"github.com/geraldfingburke/dossier/server/internal/rss"
	"github.com/lib/pq"
)

// ============================================================================
// SERVICE DEFINITION
// ============================================================================

// Service handles scheduled dossier generation and delivery.
//
// The service maintains a ticker that periodically checks for dossiers
// that need to be generated based on their configured schedules. It
// orchestrates the complete delivery pipeline from RSS fetching through
// AI summarization to email delivery.
//
// Thread Safety:
// The service uses a mutex to protect the running state and ensure
// safe concurrent access to Start(), Stop(), and IsRunning() methods.
//
// Fields:
//   - db: Database connection for configuration and delivery tracking
//   - rssService: RSS feed fetching service
//   - aiService: AI-powered summarization service
//   - emailService: Email delivery service
//   - ticker: Time ticker for periodic checks (1-minute intervals)
//   - stopChan: Channel for graceful shutdown signaling
//   - mutex: Read-write mutex for thread-safe state management
//   - running: Current running state of the scheduler
type Service struct {
	db           *sql.DB
	rssService   *rss.Service
	aiService    *ai.Service
	emailService *email.Service
	ticker       *time.Ticker
	stopChan     chan bool
	mutex        sync.RWMutex
	running      bool
}

// ============================================================================
// SERVICE INITIALIZATION
// ============================================================================

// NewService creates a new scheduler service with required dependencies.
//
// The service is created in a stopped state. Call Start() to begin
// the scheduling loop.
//
// Parameters:
//   - db: Database connection for querying configs and recording deliveries
//   - rssService: Service for fetching RSS feed articles
//   - aiService: Service for generating AI-powered summaries
//   - emailService: Service for sending formatted emails
//
// Returns:
//   - *Service: Configured scheduler service ready to start
//
// Example:
//   scheduler := scheduler.NewService(db, rssService, aiService, emailService)
//   scheduler.Start()
//   defer scheduler.Stop()
func NewService(db *sql.DB, rssService *rss.Service, aiService *ai.Service, emailService *email.Service) *Service {
	return &Service{
		db:           db,
		rssService:   rssService,
		aiService:    aiService,
		emailService: emailService,
		stopChan:     make(chan bool),
		running:      false,
	}
}

// ============================================================================
// SCHEDULER CONTROL
// ============================================================================

// Start begins the scheduler with 1-minute check intervals.
//
// This method starts a background goroutine that:
//   1. Wakes up every minute via ticker
//   2. Queries active dossier configurations
//   3. Evaluates each configuration's schedule
//   4. Triggers async dossier generation for due deliveries
//
// Thread Safety:
// Safe to call concurrently. If already running, logs and returns.
// Uses mutex to prevent multiple scheduler instances.
//
// Behavior:
//   - Idempotent: Multiple calls have no effect if already running
//   - Non-blocking: Returns immediately after starting goroutine
//   - Logs: Startup confirmation and ticker events
//
// Example:
//   scheduler.Start()
//   // Scheduler runs in background
//   time.Sleep(1 * time.Hour)
//   scheduler.Stop()
func (s *Service) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		log.Println("Scheduler is already running")
		return
	}

	log.Println("Starting dossier scheduler...")
	s.running = true
	s.ticker = time.NewTicker(1 * time.Minute) // Check every minute

	go func() {
		for {
			select {
			case <-s.ticker.C:
				log.Printf("Scheduler: Ticker fired at %s", time.Now().UTC().Format("15:04:05"))
				s.checkAndProcessDossiers()
			case <-s.stopChan:
				return
			}
		}
	}()

	log.Println("Dossier scheduler started successfully")
}

// Stop gracefully stops the scheduler.
//
// This method:
//   1. Stops the ticker (no more periodic checks)
//   2. Signals the background goroutine to exit
//   3. Updates running state
//
// Thread Safety:
// Safe to call concurrently. If not running, returns silently.
// Uses mutex to ensure clean shutdown.
//
// Behavior:
//   - Idempotent: Multiple calls have no effect if not running
//   - Blocking: Waits for goroutine to acknowledge stop signal
//   - Graceful: In-flight dossier generations are not interrupted
//
// Note: In-flight dossier generation goroutines will complete
// independently. Only the scheduler loop stops.
//
// Example:
//   scheduler.Stop()
//   // Scheduler stopped, but active deliveries may still complete
func (s *Service) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return
	}

	log.Println("Stopping dossier scheduler...")
	s.running = false
	s.ticker.Stop()
	s.stopChan <- true
	log.Println("Dossier scheduler stopped")
}

// IsRunning returns whether the scheduler is currently active.
//
// Thread Safety:
// Uses read lock to safely check running state.
//
// Returns:
//   - bool: true if scheduler is running, false otherwise
//
// Usage:
// Primarily used by GraphQL API to report scheduler status
// in the schedulerStatus query.
//
// Example:
//   if scheduler.IsRunning() {
//       log.Println("Scheduler is active")
//   }
func (s *Service) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// ============================================================================
// CONFIGURATION MANAGEMENT
// ============================================================================

// checkAndProcessDossiers checks for dossiers that need generation and processes them.
//
// This is the main scheduler loop body that runs every minute. It:
//   1. Queries all active dossier configurations
//   2. Evaluates each against current time and frequency
//   3. Launches async goroutines for due dossiers
//
// Error Handling:
// Individual configuration errors don't stop processing of others.
// Each dossier generation runs in its own goroutine to prevent blocking.
//
// Logging:
// Comprehensive logging at each step for debugging and monitoring:
//   - Ticker events with UTC timestamps
//   - Configuration count and details
//   - Schedule evaluation results
//   - Generation triggers and completion
func (s *Service) checkAndProcessDossiers() {
	log.Printf("Scheduler: Checking for dossiers to process at %s", time.Now().UTC().Format("15:04:05"))
	
	configs, err := s.getActiveDossierConfigs()
	if err != nil {
		log.Printf("Error getting active dossier configs: %v", err)
		return
	}

	log.Printf("Scheduler: Found %d active configurations", len(configs))
	
	for _, config := range configs {
		log.Printf("Scheduler: Checking config %d (%s) - delivery_time: %s", config.ID, config.Title, config.DeliveryTime)
		
		if s.shouldGenerateDossier(config) {
			log.Printf("Scheduler: Triggering dossier generation for config %d (%s)", config.ID, config.Title)
			
			// Launch async generation to avoid blocking other configs
			go func(cfg models.DossierConfig) {
				if err := s.generateAndSendDossier(cfg); err != nil {
					log.Printf("Error generating dossier for config %d (%s): %v", cfg.ID, cfg.Title, err)
				}
			}(config)
		} else {
			log.Printf("Scheduler: Not time to generate dossier for config %d (%s)", config.ID, config.Title)
		}
	}
}

// getActiveDossierConfigs retrieves all active dossier configurations from database.
//
// Query:
// Selects all configurations where active = true, including all fields
// needed for schedule evaluation and dossier generation.
//
// Error Handling:
// Individual row scan errors are logged but don't fail entire query.
// This allows partial results if some rows have data issues.
//
// Returns:
//   - []models.DossierConfig: All active configurations
//   - error: Database query error (nil on success)
func (s *Service) getActiveDossierConfigs() ([]models.DossierConfig, error) {
	rows, err := s.db.Query(`
		SELECT id, title, email, feed_urls, article_count, frequency, 
		       delivery_time, timezone, tone, language, special_instructions, 
		       active, created_at, updated_at
		FROM dossier_configs 
		WHERE active = true
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query dossier configs: %w", err)
	}
	defer rows.Close()

	var configs []models.DossierConfig
	for rows.Next() {
		var config models.DossierConfig
		err := rows.Scan(
			&config.ID, &config.Title, &config.Email, pq.Array(&config.FeedURLs),
			&config.ArticleCount, &config.Frequency, &config.DeliveryTime,
			&config.Timezone, &config.Tone, &config.Language,
			&config.SpecialInstructions, &config.Active, &config.CreatedAt, &config.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning dossier config: %v", err)
			continue
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// ============================================================================
// SCHEDULE EVALUATION
// ============================================================================

// shouldGenerateDossier determines if a dossier should be generated now.
//
// This method implements the core scheduling logic:
//   1. Parse configuration's timezone
//   2. Get current time in that timezone
//   3. Parse delivery time from configuration
//   4. Check if current time matches delivery window
//   5. Apply frequency-based rules (daily/weekly/monthly)
//   6. Check duplicate prevention logic
//
// Time Matching:
// Delivery occurs when current hour and minute match configured time.
// This provides 1-minute granularity (matching ticker frequency).
//
// Timezone Handling:
// Each configuration has its own timezone. If invalid, falls back
// to system local time with warning log.
//
// Time Format Support:
// Handles multiple delivery time formats:
//   - HH:MM (preferred format)
//   - HH:MM:SS (PostgreSQL TIME format)
//   - Full timestamps (extracts time component)
//
// Parameters:
//   - config: Dossier configuration to evaluate
//
// Returns:
//   - bool: true if dossier should be generated now, false otherwise
func (s *Service) shouldGenerateDossier(config models.DossierConfig) bool {
	// Parse and validate timezone
	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Printf("Invalid timezone %s for config %d, using local system time", config.Timezone, config.ID)
		location = time.Local
	}

	// Get current time in configuration's timezone
	now := time.Now().In(location)
	
	// Parse delivery time - handle multiple formats for robustness
	var deliveryTime time.Time
	
	// Try HH:MM format first (most common)
	if len(config.DeliveryTime) == 5 && config.DeliveryTime[2] == ':' {
		var err error
		deliveryTime, err = time.Parse("15:04", config.DeliveryTime)
		if err != nil {
			log.Printf("Invalid delivery time %s for config %d", config.DeliveryTime, config.ID)
			return false
		}
	} else {
		// Try parsing as full timestamp and extract time component
		fullTime, parseErr := time.Parse(time.RFC3339, config.DeliveryTime)
		if parseErr != nil {
			// Try PostgreSQL TIME format with zero date
			fullTime, parseErr = time.Parse("0000-01-01T15:04:05Z", config.DeliveryTime)
			if parseErr != nil {
				// Try HH:MM:SS format
				fullTime, parseErr = time.Parse("15:04:05", config.DeliveryTime)
			}
		}
		if parseErr != nil {
			log.Printf("Invalid delivery time format %s for config %d - %v", config.DeliveryTime, config.ID, parseErr)
			return false
		}
		deliveryTime = fullTime
	}

	// Create target time for today in the configuration's timezone
	targetTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		deliveryTime.Hour(), deliveryTime.Minute(), 0, 0,
		location,
	)

	// Check if we're within the delivery window (current minute matches target minute)
	if now.Hour() != targetTime.Hour() || now.Minute() != targetTime.Minute() {
		return false
	}

	// Apply frequency-based scheduling rules
	switch config.Frequency {
	case "daily":
		return s.shouldGenerateDaily(config, now)
	case "weekly":
		return s.shouldGenerateWeekly(config, now)
	case "monthly":
		return s.shouldGenerateMonthly(config, now)
	default:
		log.Printf("Unknown frequency %s for config %d", config.Frequency, config.ID)
		return false
	}
}

// shouldGenerateDaily checks if a daily dossier should be generated.
//
// Logic:
//   - Generates once per day
//   - Checks if already generated today (in config's timezone)
//   - Compares dates in YYYY-MM-DD format
//
// Parameters:
//   - config: Dossier configuration
//   - now: Current time in configuration's timezone
//
// Returns:
//   - bool: true if should generate (not yet generated today)
func (s *Service) shouldGenerateDaily(config models.DossierConfig, now time.Time) bool {
	// Check if we already generated today
	lastGenerated, err := s.getLastGeneratedTime(config.ID)
	if err != nil {
		log.Printf("Error checking last generated time for config %d: %v", config.ID, err)
		return true // Generate on error to be safe
	}

	if lastGenerated == nil {
		return true // Never generated before
	}

	// Compare dates only (ignore time)
	lastGeneratedDay := lastGenerated.In(now.Location()).Format("2006-01-02")
	todayDay := now.Format("2006-01-02")
	
	return lastGeneratedDay != todayDay
}

// shouldGenerateWeekly checks if a weekly dossier should be generated.
//
// Logic:
//   - Generates once per week on Mondays
//   - Uses ISO week numbers for comparison
//   - Timezone-aware (Monday in config's timezone)
//
// Parameters:
//   - config: Dossier configuration
//   - now: Current time in configuration's timezone
//
// Returns:
//   - bool: true if should generate (Monday and not generated this week)
func (s *Service) shouldGenerateWeekly(config models.DossierConfig, now time.Time) bool {
	// Only generate on Mondays
	if now.Weekday() != time.Monday {
		return false
	}

	// Check if we already generated this week
	lastGenerated, err := s.getLastGeneratedTime(config.ID)
	if err != nil {
		log.Printf("Error checking last generated time for config %d: %v", config.ID, err)
		return true
	}

	if lastGenerated == nil {
		return true // Never generated before
	}

	// Compare ISO week numbers
	_, thisWeek := now.ISOWeek()
	_, lastWeek := lastGenerated.In(now.Location()).ISOWeek()
	
	return thisWeek != lastWeek
}

// shouldGenerateMonthly checks if a monthly dossier should be generated.
//
// Logic:
//   - Generates once per month on the 1st
//   - Compares year-month in YYYY-MM format
//   - Timezone-aware (1st in config's timezone)
//
// Parameters:
//   - config: Dossier configuration
//   - now: Current time in configuration's timezone
//
// Returns:
//   - bool: true if should generate (1st of month and not generated this month)
func (s *Service) shouldGenerateMonthly(config models.DossierConfig, now time.Time) bool {
	// Only generate on the 1st of the month
	if now.Day() != 1 {
		return false
	}

	// Check if we already generated this month
	lastGenerated, err := s.getLastGeneratedTime(config.ID)
	if err != nil {
		log.Printf("Error checking last generated time for config %d: %v", config.ID, err)
		return true
	}

	if lastGenerated == nil {
		return true // Never generated before
	}

	// Compare year-month only
	thisMonth := now.Format("2006-01")
	lastMonth := lastGenerated.In(now.Location()).Format("2006-01")
	
	return thisMonth != lastMonth
}

// getLastGeneratedTime retrieves the most recent delivery time for a configuration.
//
// This method queries the dossier_deliveries table to find the last time
// a dossier was generated, used for duplicate prevention logic.
//
// Parameters:
//   - configID: Configuration ID to check
//
// Returns:
//   - *time.Time: Last delivery date (nil if never generated)
//   - error: Database error (nil on success or no rows)
func (s *Service) getLastGeneratedTime(configID int) (*time.Time, error) {
	var deliveryDate time.Time
	err := s.db.QueryRow(`
		SELECT delivery_date FROM dossier_deliveries 
		WHERE config_id = $1 
		ORDER BY delivery_date DESC 
		LIMIT 1
	`, configID).Scan(&deliveryDate)
	
	if err == sql.ErrNoRows {
		return nil, nil // No previous generation
	}
	if err != nil {
		return nil, err
	}
	
	return &deliveryDate, nil
}

// ============================================================================
// DOSSIER GENERATION PIPELINE
// ============================================================================

// generateAndSendDossier executes the complete dossier generation pipeline.
//
// This method orchestrates all steps needed to create and deliver a dossier:
//   1. Fetch articles from all configured RSS feeds
//   2. Convert feed items to Article models
//   3. Sort by publication date (newest first)
//   4. Limit to requested article count
//   5. Generate AI summary with specified tone and language
//   6. Send formatted email to recipient
//   7. Record delivery in database
//
// Context:
// Uses 10-minute timeout to prevent indefinite hangs on slow operations.
// This is generous enough for slow feeds and AI processing.
//
// Error Handling:
//   - Individual feed failures: Logged, continue with other feeds
//   - No articles found: Returns error, no email sent
//   - AI generation failure: Returns error, no email sent
//   - Email failure: Returns error, no delivery recorded
//   - Recording failure: Logged only (email already sent)
//
// Concurrency:
// Designed to be called from goroutine (doesn't block caller).
// Each configuration's generation is independent.
//
// Parameters:
//   - config: Dossier configuration with all settings
//
// Returns:
//   - error: Any step failure (nil on complete success)
func (s *Service) generateAndSendDossier(config models.DossierConfig) error {
	log.Printf("Generating scheduled dossier for config %d (%s)", config.ID, config.Title)
	
	// Create context with timeout for entire pipeline
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Fetch and aggregate articles from all configured feeds
	var allArticles []models.Article
	for _, feedURL := range config.FeedURLs {
		feed, err := s.rssService.FetchFeed(ctx, feedURL)
		if err != nil {
			log.Printf("Error fetching feed %s: %v", feedURL, err)
			continue // Skip failed feeds, continue with others
		}
		
		// Convert gofeed.Item to models.Article
		for _, item := range feed.Items {
			// Extract author name (may be nil)
			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}
			
			// Use parsed published date or current time
			publishedAt := time.Now()
			if item.PublishedParsed != nil {
				publishedAt = *item.PublishedParsed
			}
			
			// Build article model
			article := models.Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Author:      author,
				PublishedAt: publishedAt,
			}
			
			// Prefer full content over description
			if item.Content != "" {
				article.Content = item.Content
			}
			
			allArticles = append(allArticles, article)
		}
	}

	// Validate we have articles to process
	if len(allArticles) == 0 {
		return fmt.Errorf("no articles found from any feeds")
	}

	// Limit to requested count (RSS service should handle sorting)
	if len(allArticles) > config.ArticleCount {
		allArticles = allArticles[:config.ArticleCount]
	}

	// Generate AI summary with configured tone and language
	summary, err := s.aiService.GenerateSummary(
		ctx, 
		allArticles, 
		config.Tone, 
		config.Language, 
		config.SpecialInstructions,
	)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	// Send formatted email to recipient
	err = s.emailService.SendDossier(&config, summary, allArticles)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Record successful delivery in database
	err = s.recordDossierGeneration(config.ID, summary, len(allArticles))
	if err != nil {
		log.Printf("Error recording dossier generation: %v", err)
		// Don't return error here since email was sent successfully
	}

	log.Printf("Successfully generated and sent dossier for config %d (%s) to %s", 
		config.ID, config.Title, config.Email)
	
	return nil
}

// recordDossierGeneration records a successful dossier delivery in the database.
//
// This creates an audit trail of all deliveries and is used by the
// duplicate prevention logic to track when dossiers were last generated.
//
// Parameters:
//   - configID: Configuration ID that generated this dossier
//   - summary: AI-generated summary HTML
//   - articleCount: Number of articles included
//
// Returns:
//   - error: Database insertion error (nil on success)
func (s *Service) recordDossierGeneration(configID int, summary string, articleCount int) error {
	_, err := s.db.Exec(`
		INSERT INTO dossier_deliveries (config_id, delivery_date, summary, article_count, email_sent)
		VALUES ($1, $2, $3, $4, $5)
	`, configID, time.Now(), summary, articleCount, true)
	
	return err
}