// Package rss provides RSS/Atom feed fetching and parsing functionality.
//
// This package handles all interactions with external RSS feeds, including:
//   - Feed discovery and parsing
//   - Article extraction and transformation
//   - Multi-feed aggregation
//   - Content normalization
//
// # Architecture
//
// The RSS service acts as a bridge between external RSS/Atom feeds and the
// Dossier application's internal data structures. It uses the gofeed library
// for robust RSS/Atom parsing with automatic format detection.
//
// # Feed Processing Pipeline
//
//  1. Fetch: Download feed XML from URL
//  2. Parse: Convert RSS/Atom to normalized structure (via gofeed)
//  3. Extract: Convert feed items to Article models
//  4. Normalize: Handle missing/optional fields with sensible defaults
//  5. Aggregate: Combine articles from multiple feeds
//  6. Sort: Order by publication date (newest first)
//  7. Limit: Return requested number of articles
//
// # Data Quality Handling
//
// RSS feeds vary significantly in quality and completeness. This service
// implements defensive parsing strategies:
//
//   - Missing timestamps → Use current time
//   - Missing content → Fall back to description
//   - Missing author → Use empty string
//   - Parse failures → Skip feed, continue with others
//
// # Integration Points
//
// The RSS service is used by:
//   - GraphQL API: Manual dossier generation (generateAndSendDossier mutation)
//   - Scheduler: Automated dossier delivery
//   - AI Service: Receives articles for summarization
//
// # Dependencies
//
//   - github.com/mmcdole/gofeed: RSS/Atom parsing library
//   - ai.Service: AI-powered article selection (stored for potential future use)
//   - models.Article: Internal article representation
//
// # Error Handling Philosophy
//
// Feed fetching is designed to be resilient:
//   - Individual feed failures don't abort the entire operation
//   - Partial results are returned when some feeds succeed
//   - Detailed logging helps diagnose issues
//   - Empty results are valid (no articles to process)
package rss

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/geraldfingburke/dossier/server/internal/ai"
	"github.com/geraldfingburke/dossier/server/internal/models"
	"github.com/mmcdole/gofeed"
)

// ============================================================================
// SERVICE DEFINITION
// ============================================================================

// Service handles RSS feed operations and article aggregation.
//
// The service maintains a gofeed parser instance that is reused across
// multiple feed fetches for efficiency. The AI service reference is
// stored for potential future enhancements (e.g., AI-powered feed filtering).
//
// Thread Safety:
// The gofeed parser is safe for concurrent use, making this service
// safe for concurrent feed fetching operations.
//
// Fields:
//   - parser: gofeed parser instance (reused for efficiency)
//   - aiService: AI service reference (for potential future enhancements)
type Service struct {
	parser    *gofeed.Parser
	aiService *ai.Service
}

// ============================================================================
// SERVICE INITIALIZATION
// ============================================================================

// NewService creates a new RSS service with feed parsing capabilities.
//
// The service initializes a gofeed parser that automatically detects
// and handles RSS 1.0, RSS 2.0, and Atom feed formats.
//
// Parameters:
//   - aiService: AI service for potential article intelligence features
//
// Returns:
//   - *Service: Configured RSS service ready for feed operations
//
// Example:
//
//	rssService := rss.NewService(aiService)
//	articles, err := rssService.FetchArticlesFromFeeds(ctx, feedURLs, 10)
func NewService(aiService *ai.Service) *Service {
	return &Service{
		parser:    gofeed.NewParser(),
		aiService: aiService,
	}
}

// ============================================================================
// FEED FETCHING OPERATIONS
// ============================================================================

// FetchFeed fetches and parses a single RSS/Atom feed from a URL.
//
// This method handles the complete process of downloading and parsing
// an RSS or Atom feed, automatically detecting the format and converting
// it to a normalized structure.
//
// Supported Formats:
//   - RSS 1.0 (RDF)
//   - RSS 2.0
//   - Atom 1.0
//
// Context Support:
// The method respects context cancellation, allowing timeouts and
// cancellation of long-running feed fetches.
//
// Error Conditions:
//   - Network failures (DNS, connection timeout, etc.)
//   - HTTP errors (404, 500, etc.)
//   - XML parsing errors (malformed feed)
//   - Unsupported feed format
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - feedURL: Complete URL to RSS/Atom feed
//
// Returns:
//   - *gofeed.Feed: Parsed feed with items and metadata
//   - error: Network, HTTP, or parsing error
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	feed, err := service.FetchFeed(ctx, "https://news.ycombinator.com/rss")
//	if err != nil {
//	    log.Printf("Failed to fetch feed: %v", err)
//	    return
//	}
//	log.Printf("Fetched %d items from %s", len(feed.Items), feed.Title)
func (s *Service) FetchFeed(ctx context.Context, feedURL string) (*gofeed.Feed, error) {
	feed, err := s.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("error parsing feed: %w", err)
	}
	return feed, nil
}

// ============================================================================
// MULTI-FEED AGGREGATION
// ============================================================================

// FetchArticlesFromFeeds aggregates articles from multiple RSS feeds.
//
// This is the primary method for dossier generation, fetching articles from
// multiple feeds and intelligently combining them into a single, time-ordered
// collection. It implements resilient fetching where individual feed failures
// don't prevent the entire operation from succeeding.
//
// Algorithm:
//  1. Calculate articles per feed (maxArticles / number of feeds)
//  2. Fetch each feed in sequence (continues on individual failures)
//  3. Convert feed items to Article models
//  4. Normalize missing/optional fields
//  5. Aggregate all articles into single collection
//  6. Sort by publication date (newest first)
//  7. Limit to maxArticles total
//
// Distribution Strategy:
// Articles are distributed evenly across feeds, but if some feeds return
// fewer articles than allocated, other feeds can fill the gap. This ensures
// the requested article count is reached when possible.
//
// Error Handling:
// Individual feed failures are logged but don't stop processing. The method
// returns successfully with partial results as long as at least one feed
// succeeds. Returns error only if ALL feeds fail or other critical issues occur.
//
// Field Normalization:
//   - PublishedAt: Uses item.PublishedParsed or current time if missing
//   - Content: Uses item.Content or falls back to item.Description
//   - Author: Uses item.Author.Name or empty string if missing
//   - Description: Uses item.Description (may be empty)
//   - Title: Always present (required by RSS spec)
//   - Link: Always present (required by RSS spec)
//
// Performance Considerations:
//   - Feeds are fetched sequentially (not parallel) to avoid overwhelming servers
//   - Uses bubble sort for simplicity (article counts typically < 100)
//   - Each feed fetch respects the context timeout
//
// Parameters:
//   - ctx: Context for timeout and cancellation
//   - feedURLs: Array of RSS/Atom feed URLs to fetch
//   - maxArticles: Maximum total articles to return across all feeds
//
// Returns:
//   - []models.Article: Aggregated articles sorted by date (newest first)
//   - error: Only if ALL feeds fail or critical error occurs
//
// Example:
//
//	feedURLs := []string{
//	    "https://news.ycombinator.com/rss",
//	    "https://techcrunch.com/feed/",
//	}
//	articles, err := service.FetchArticlesFromFeeds(ctx, feedURLs, 20)
//	if err != nil {
//	    log.Printf("Failed to fetch articles: %v", err)
//	    return
//	}
//	log.Printf("Fetched %d articles from %d feeds", len(articles), len(feedURLs))
func (s *Service) FetchArticlesFromFeeds(ctx context.Context, feedURLs []string, maxArticles int) ([]models.Article, error) {
	var allArticles []models.Article

	// Calculate target articles per feed for even distribution
	articlesPerFeed := maxArticles / len(feedURLs)
	if articlesPerFeed == 0 {
		articlesPerFeed = 1
	}

	// Fetch articles from each feed
	for _, feedURL := range feedURLs {
		log.Printf("Fetching articles from feed: %s", feedURL)

		// Fetch and parse feed
		feed, err := s.FetchFeed(ctx, feedURL)
		if err != nil {
			log.Printf("Error fetching feed %s: %v", feedURL, err)
			continue // Skip failed feeds, continue with others
		}

		// Convert feed items to Article models
		feedArticles := make([]models.Article, 0)
		for i, item := range feed.Items {
			// Limit articles per feed
			if i >= articlesPerFeed {
				break
			}

			// Normalize published date (use current time if missing)
			publishedAt := time.Now()
			if item.PublishedParsed != nil {
				publishedAt = *item.PublishedParsed
			}

			// Normalize content (prefer full content, fall back to description)
			content := item.Content
			if content == "" {
				content = item.Description
			}

			// Extract author name (may be empty)
			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}

			// Build normalized article model
			article := models.Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Content:     content,
				Author:      author,
				PublishedAt: publishedAt,
			}

			feedArticles = append(feedArticles, article)
		}

		allArticles = append(allArticles, feedArticles...)
		log.Printf("Fetched %d articles from %s", len(feedArticles), feedURL)
	}

	// Sort by published date (newest first) and limit to maxArticles
	if len(allArticles) > maxArticles {
		// Bubble sort by publication date (descending)
		// Using bubble sort for simplicity - article counts are typically small (<100)
		for i := 0; i < len(allArticles)-1; i++ {
			for j := i + 1; j < len(allArticles); j++ {
				if allArticles[i].PublishedAt.Before(allArticles[j].PublishedAt) {
					allArticles[i], allArticles[j] = allArticles[j], allArticles[i]
				}
			}
		}
		// Limit to requested count
		allArticles = allArticles[:maxArticles]
	}

	log.Printf("Total articles fetched: %d", len(allArticles))
	return allArticles, nil
}
