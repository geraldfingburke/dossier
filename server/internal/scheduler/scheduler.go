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

// Service handles scheduled dossier generation
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

// NewService creates a new scheduler service
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

// Start begins the scheduler with 1-minute check intervals
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

// Stop stops the scheduler
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

// IsRunning returns whether the scheduler is currently running
func (s *Service) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// checkAndProcessDossiers checks for dossiers that need to be generated and processes them
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

// getActiveDossierConfigs retrieves all active dossier configurations
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

// shouldGenerateDossier determines if a dossier should be generated now based on schedule
func (s *Service) shouldGenerateDossier(config models.DossierConfig) bool {
	// Parse the timezone
	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Printf("Invalid timezone %s for config %d, using UTC", config.Timezone, config.ID)
		location = time.UTC
	}

	// Get current time in the config's timezone
	now := time.Now().In(location)
	
	// Parse delivery time - handle both HH:MM format and full timestamp
	var deliveryTime time.Time
	
	// Try parsing as HH:MM first
	if len(config.DeliveryTime) == 5 && config.DeliveryTime[2] == ':' {
		var err error
		deliveryTime, err = time.Parse("15:04", config.DeliveryTime)
		if err != nil {
			log.Printf("Invalid delivery time %s for config %d", config.DeliveryTime, config.ID)
			return false
		}
	} else {
		// Try parsing as full timestamp and extract time
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

	// Create target time for today in the config's timezone
	targetTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		deliveryTime.Hour(), deliveryTime.Minute(), 0, 0,
		location,
	)

	// Check if we're within the delivery window (current minute matches target minute)
	if now.Hour() != targetTime.Hour() || now.Minute() != targetTime.Minute() {
		return false
	}

	// Check frequency-based conditions
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

// shouldGenerateDaily checks if a daily dossier should be generated
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

	// Check if last generation was today
	lastGeneratedDay := lastGenerated.In(now.Location()).Format("2006-01-02")
	todayDay := now.Format("2006-01-02")
	
	return lastGeneratedDay != todayDay
}

// shouldGenerateWeekly checks if a weekly dossier should be generated (Mondays)
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

	// Check if last generation was this week
	_, thisWeek := now.ISOWeek()
	_, lastWeek := lastGenerated.In(now.Location()).ISOWeek()
	
	return thisWeek != lastWeek
}

// shouldGenerateMonthly checks if a monthly dossier should be generated (1st of month)
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

	// Check if last generation was this month
	thisMonth := now.Format("2006-01")
	lastMonth := lastGenerated.In(now.Location()).Format("2006-01")
	
	return thisMonth != lastMonth
}

// getLastGeneratedTime gets the last time a dossier was generated for a config
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

// generateAndSendDossier generates and sends a dossier for the given configuration
func (s *Service) generateAndSendDossier(config models.DossierConfig) error {
	log.Printf("Generating scheduled dossier for config %d (%s)", config.ID, config.Title)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Fetch articles from all feeds
	var allArticles []models.Article
	for _, feedURL := range config.FeedURLs {
		feed, err := s.rssService.FetchFeed(ctx, feedURL)
		if err != nil {
			log.Printf("Error fetching feed %s: %v", feedURL, err)
			continue
		}
		
		// Convert gofeed.Item to models.Article
		for _, item := range feed.Items {
			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}
			
			publishedAt := time.Now()
			if item.PublishedParsed != nil {
				publishedAt = *item.PublishedParsed
			}
			
			article := models.Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Author:      author,
				PublishedAt: publishedAt,
			}
			if item.Content != "" {
				article.Content = item.Content
			}
			allArticles = append(allArticles, article)
		}
	}

	if len(allArticles) == 0 {
		return fmt.Errorf("no articles found from any feeds")
	}

	// Sort by publish date (newest first) and limit to requested count
	// Note: RSS service should handle this, but we'll ensure it here
	if len(allArticles) > config.ArticleCount {
		allArticles = allArticles[:config.ArticleCount]
	}

	// Generate summary using AI service
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

	// Send email
	err = s.emailService.SendDossier(&config, summary, allArticles)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Record the dossier generation
	err = s.recordDossierGeneration(config.ID, summary, len(allArticles))
	if err != nil {
		log.Printf("Error recording dossier generation: %v", err)
		// Don't return error here since email was sent successfully
	}

	log.Printf("Successfully generated and sent dossier for config %d (%s) to %s", 
		config.ID, config.Title, config.Email)
	
	return nil
}

// recordDossierGeneration records that a dossier was generated
func (s *Service) recordDossierGeneration(configID int, summary string, articleCount int) error {
	_, err := s.db.Exec(`
		INSERT INTO dossier_deliveries (config_id, delivery_date, summary, article_count, email_sent)
		VALUES ($1, $2, $3, $4, $5)
	`, configID, time.Now(), summary, articleCount, true)
	
	return err
}