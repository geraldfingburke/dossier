// Package graphql provides a GraphQL API for the Dossier application.
//
// This package implements a complete GraphQL schema for managing automated news digests,
// including dossier configurations, RSS feed management, AI-powered summarization,
// email delivery, and tone customization.
//
// # Architecture Overview
//
// The GraphQL API is structured into three main layers:
//   1. Type Definitions: GraphQL types representing domain objects
//   2. Query Operations: Read-only data retrieval
//   3. Mutation Operations: State-changing operations (create, update, delete)
//
// # Key Features
//
//   - Dossier Configuration Management: Create, read, update, delete dossier configs
//   - Delivery History: Query past dossier deliveries with filtering
//   - Manual Triggering: Generate and send dossiers on-demand
//   - Email Testing: Validate SMTP configuration and send test emails
//   - Tone Customization: Manage AI tone presets for summary generation
//   - Scheduler Monitoring: Check scheduler status and active configurations
//
// # Integration Points
//
// The GraphQL API orchestrates multiple services:
//   - Database: PostgreSQL for persistent storage
//   - RSS Service: Feed fetching and article aggregation
//   - AI Service: Ollama-powered content summarization
//   - Email Service: SMTP-based email delivery
//   - Scheduler Service: Automated delivery scheduling
//
// # GraphQL Schema Structure
//
// Types:
//   - DossierConfig: User configuration for automated digests
//   - Dossier: Historical delivery record
//   - Tone: AI tone preset with system/custom variants
//   - SchedulerStatus: Real-time scheduler information
//
// Queries:
//   - dossierConfigs: List all active configurations
//   - dossierConfig(id): Get single configuration by ID
//   - dossiers(configId, limit): Query delivery history
//   - tones: List all available AI tones
//   - tone(id): Get single tone by ID
//   - schedulerStatus: Current scheduler state
//
// Mutations:
//   - createDossierConfig: Create new configuration
//   - updateDossierConfig: Update existing configuration
//   - deleteDossierConfig: Delete configuration
//   - generateAndSendDossier: Manually trigger delivery
//   - sendTestEmail: Send test email with sample data
//   - testEmailConnection: Validate SMTP settings
//   - createTone: Create custom AI tone
//   - updateTone: Update custom tone
//   - deleteTone: Delete custom tone (system defaults protected)
package graphql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/geraldfingburke/dossier/server/internal/ai"
	"github.com/geraldfingburke/dossier/server/internal/email"
	"github.com/geraldfingburke/dossier/server/internal/models"
	"github.com/geraldfingburke/dossier/server/internal/rss"
	"github.com/geraldfingburke/dossier/server/internal/scheduler"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/lib/pq"
)

// ============================================================================
// GRAPHQL HANDLER
// ============================================================================

// Handler creates the GraphQL HTTP handler for the Dossier API.
//
// This function constructs the complete GraphQL schema including all type definitions,
// queries, and mutations. It integrates with multiple backend services to provide
// a unified API for the frontend application.
//
// Service Dependencies:
//   - db: PostgreSQL database connection for persistent storage
//   - rssService: RSS feed fetching and article aggregation
//   - aiService: AI-powered content summarization via Ollama
//   - emailService: SMTP-based email delivery
//   - schedulerService: Automated delivery scheduling and monitoring
//
// GraphQL Configuration:
//   - Pretty: Formatted JSON responses for readability
//   - GraphiQL: Interactive GraphQL IDE enabled for development
//
// Parameters:
//   - db: Database connection
//   - rssService: RSS feed service
//   - aiService: AI summarization service
//   - emailService: Email delivery service
//   - schedulerService: Delivery scheduler service
//
// Returns:
//   - *handler.Handler: Configured GraphQL HTTP handler
//   - error: Schema creation or validation error
func Handler(db *sql.DB, rssService *rss.Service, aiService *ai.Service, emailService *email.Service, schedulerService *scheduler.Service) (*handler.Handler, error) {
	// ========================================================================
	// TYPE DEFINITIONS
	// ========================================================================

	// DossierConfig GraphQL type represents a user's automated digest configuration.
	//
	// This type maps to the dossier_configs database table and includes all settings
	// needed to generate and deliver personalized news digests.
	//
	// Fields:
	//   - id: Unique configuration identifier
	//   - title: User-friendly name for the dossier
	//   - email: Recipient email address
	//   - feedUrls: Array of RSS feed URLs to aggregate
	//   - articleCount: Number of articles to include per digest
	//   - frequency: Delivery schedule (daily, weekly, etc.)
	//   - deliveryTime: Time of day for scheduled delivery (HH:MM format)
	//   - timezone: IANA timezone for delivery scheduling
	//   - tone: AI tone preset for summary generation
	//   - language: Target language for summaries
	//   - specialInstructions: Custom AI instructions
	//   - active: Whether automated delivery is enabled
	//   - createdAt: Configuration creation timestamp
	dossierConfigType := graphql.NewObject(graphql.ObjectConfig{
		Name: "DossierConfig",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"feedUrls": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
			},
			"articleCount": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"frequency": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"deliveryTime": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				// Custom resolver to format delivery time as HH:MM for frontend.
				// Database stores time as HH:MM:SS but UI only needs hour and minute.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var config *models.DossierConfig

					// Handle both pointer and value types for flexibility
					switch v := p.Source.(type) {
					case *models.DossierConfig:
						config = v
					case models.DossierConfig:
						config = &v
					default:
						return nil, fmt.Errorf("unexpected source type: %T", v)
					}

					// Extract HH:MM from HH:MM:SS format
					if len(config.DeliveryTime) >= 5 && config.DeliveryTime[2] == ':' {
						return config.DeliveryTime[:5], nil
					}
					return config.DeliveryTime, nil
				},
			},
			"timezone": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"tone": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"language": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"specialInstructions": &graphql.Field{
				Type: graphql.String,
			},
			"active": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	// DossierConfigInput input type for create/update mutations.
	//
	// This input type defines the structure for creating and updating dossier configurations.
	// All fields are required except tone, language, and specialInstructions which have defaults.
	//
	// Default Values:
	//   - tone: "professional" (applied in resolver)
	//   - language: "English" (applied in resolver)
	//   - specialInstructions: "" (empty string)
	dossierConfigInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "DossierConfigInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"title": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"feedUrls": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
			},
			"articleCount": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"frequency": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"deliveryTime": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"timezone": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"tone": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"language": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"specialInstructions": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})

	// SchedulerStatus GraphQL type represents the current state of the delivery scheduler.
	//
	// This type provides real-time information about the automated delivery system,
	// useful for admin dashboards and monitoring.
	//
	// Fields:
	//   - running: Whether the scheduler is actively running
	//   - nextCheck: Timestamp of next scheduled check (currently null, TODO)
	//   - activeDossiers: Count of enabled dossier configurations
	schedulerStatusType := graphql.NewObject(graphql.ObjectConfig{
		Name: "SchedulerStatus",
		Fields: graphql.Fields{
			"running": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"nextCheck": &graphql.Field{
				Type: graphql.String,
			},
			"activeDossiers": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	})

	// Dossier (delivery) GraphQL type represents a historical dossier delivery.
	//
	// This type maps to the dossier_deliveries table and provides access to
	// past deliveries for viewing and auditing purposes.
	//
	// Fields:
	//   - id: Unique delivery identifier
	//   - configId: Reference to the dossier configuration
	//   - subject: Email subject line (derived from config title)
	//   - content: AI-generated summary content
	//   - sentAt: Delivery timestamp
	dossierType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Dossier",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"configId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"subject": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"content": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"sentAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	// Tone GraphQL type represents an AI tone preset for summary generation.
	//
	// Tones control the style and voice of AI-generated summaries. The system
	// includes default tones (protected from deletion) and supports custom tones.
	//
	// Fields:
	//   - id: Unique tone identifier
	//   - name: Display name (e.g., "Professional", "Casual", "Pirate")
	//   - prompt: AI instruction text for tone application
	//   - isSystemDefault: Whether this is a protected system tone
	//   - createdAt: Tone creation timestamp
	//   - updatedAt: Last modification timestamp
	//
	// System Default Tones:
	//   - Professional, Casual, Academic, Creative, Technical
	//   - Concise, Detailed, Humorous, Serious, Pirate
	//
	// Note: System default tones cannot be updated or deleted via GraphQL mutations.
	toneType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Tone",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"prompt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"isSystemDefault": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	// ToneInput GraphQL input type for tone create/update mutations.
	//
	// This simplified input type is used when creating or updating custom tones.
	// System default tones are managed via database migrations.
	//
	// Fields:
	//   - name: Display name for the tone
	//   - prompt: AI instruction text (how to apply this tone)
	toneInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "ToneInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"prompt": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	// ========================================================================
	// QUERY OPERATIONS
	// ========================================================================

	// Define the root query with all read-only operations.
	//
	// Query operations provide data retrieval without side effects:
	//   - dossierConfigs: List all active configurations
	//   - dossierConfig: Get single configuration by ID
	//   - schedulerStatus: Get scheduler state and active count
	//   - dossiers: Query delivery history with optional filtering
	//   - tones: List all available AI tones
	//   - tone: Get single tone by ID
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"dossierConfigs": &graphql.Field{
				Type: graphql.NewList(dossierConfigType),
				// Retrieves all active dossier configurations ordered by creation date.
				//
				// Returns:
				//   - List of DossierConfig objects
				//   - Only includes active configurations (active = true)
				//   - Sorted by created_at descending (newest first)
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rows, err := db.QueryContext(p.Context, `
						SELECT id, title, email, feed_urls, article_count, frequency, 
							   delivery_time::text, timezone, tone, language, special_instructions, 
							   active, created_at
						FROM dossier_configs
						WHERE active = true
						ORDER BY created_at DESC
					`)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var configs []models.DossierConfig
					for rows.Next() {
						var config models.DossierConfig
						err := rows.Scan(&config.ID, &config.Title, &config.Email, pq.Array(&config.FeedURLs),
							&config.ArticleCount, &config.Frequency, &config.DeliveryTime,
							&config.Timezone, &config.Tone, &config.Language,
							&config.SpecialInstructions, &config.Active, &config.CreatedAt)
						if err != nil {
							return nil, err
						}
						configs = append(configs, config)
					}
					return configs, nil
				},
			},
			"dossierConfig": &graphql.Field{
				Type: dossierConfigType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				// Retrieves a single dossier configuration by ID.
				//
				// Arguments:
				//   - id: Configuration ID (required)
				//
				// Returns:
				//   - DossierConfig object if found
				//   - null if ID doesn't exist
				//   - error for database issues
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)

					var config models.DossierConfig
					err := db.QueryRowContext(p.Context, `
						SELECT id, title, email, feed_urls, article_count, frequency, 
							   delivery_time::text, timezone, tone, language, special_instructions, 
							   active, created_at
						FROM dossier_configs WHERE id = $1
					`, id).Scan(&config.ID, &config.Title, &config.Email, pq.Array(&config.FeedURLs),
						&config.ArticleCount, &config.Frequency, &config.DeliveryTime,
						&config.Timezone, &config.Tone, &config.Language,
						&config.SpecialInstructions, &config.Active, &config.CreatedAt)
					if err != nil {
						if err == sql.ErrNoRows {
							return nil, nil
						}
						return nil, err
					}
					return &config, nil
				},
			},
			"schedulerStatus": &graphql.Field{
				Type: schedulerStatusType,
				// Retrieves current scheduler status and counts active dossiers.
				//
				// Returns:
				//   - running: Boolean indicating if scheduler is active
				//   - nextCheck: Next scheduled check time (currently null/TODO)
				//   - activeDossiers: Count of enabled dossier configurations
				//
				// Use Cases:
				//   - Admin dashboard monitoring
				//   - System health checks
				//   - User feedback on automation status
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Count active dossier configs
					var activeCount int
					err := db.QueryRowContext(p.Context, `
						SELECT COUNT(*) FROM dossier_configs WHERE active = true
					`).Scan(&activeCount)
					if err != nil {
						return nil, err
					}

					return map[string]interface{}{
						"running":        schedulerService.IsRunning(),
						"nextCheck":      nil, // TODO: implement next check time
						"activeDossiers": activeCount,
					}, nil
				},
			},
			"dossiers": &graphql.Field{
				Type: graphql.NewList(dossierType),
				Args: graphql.FieldConfigArgument{
					"configId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				// Retrieves historical dossier deliveries with optional filtering.
				//
				// Arguments:
				//   - configId: Filter by specific dossier configuration (optional)
				//   - limit: Maximum number of results to return (optional)
				//
				// Returns:
				//   - List of Dossier (delivery) objects
				//   - Sorted by delivery_date descending (newest first)
				//
				// Use Cases:
				//   - Viewing delivery archive for a specific dossier
				//   - Displaying recent deliveries across all dossiers
				//   - Audit trail for email delivery
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					configId, hasConfigId := p.Args["configId"]
					limit, hasLimit := p.Args["limit"]

					query := `
						SELECT dd.id, dd.config_id, dc.title as subject, dd.summary as content, dd.delivery_date
						FROM dossier_deliveries dd
						JOIN dossier_configs dc ON dd.config_id = dc.id
					`
					args := []interface{}{}
					argIndex := 1

					if hasConfigId {
						query += " WHERE dd.config_id = $" + fmt.Sprintf("%d", argIndex)
						args = append(args, configId)
						argIndex++
					}

					query += " ORDER BY dd.delivery_date DESC"

					if hasLimit {
						query += " LIMIT $" + fmt.Sprintf("%d", argIndex)
						args = append(args, limit)
					}

					rows, err := db.QueryContext(p.Context, query, args...)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var dossiers []map[string]interface{}
					for rows.Next() {
						var id, configId int
						var subject, content, sentAt string

						err := rows.Scan(&id, &configId, &subject, &content, &sentAt)
						if err != nil {
							return nil, err
						}

						dossiers = append(dossiers, map[string]interface{}{
							"id":       fmt.Sprintf("%d", id),
							"configId": fmt.Sprintf("%d", configId),
							"subject":  subject,
							"content":  content,
							"sentAt":   sentAt,
						})
					}

					return dossiers, nil
				},
			},
			"tones": &graphql.Field{
				Type: graphql.NewList(toneType),
				// Retrieves all available AI tone presets.
				//
				// Returns:
				//   - List of Tone objects
				//   - Sorted by system default status (system tones first), then by name
				//   - Includes both system defaults and custom user-created tones
				//
				// Sorting:
				//   1. System default tones appear first
				//   2. Custom tones appear after
				//   3. Alphabetically sorted within each group
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rows, err := db.QueryContext(p.Context, `
						SELECT id, name, prompt, is_system_default, created_at, updated_at 
						FROM tones 
						ORDER BY is_system_default DESC, name ASC
					`)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var tones []models.Tone
					for rows.Next() {
						var tone models.Tone
						err := rows.Scan(&tone.ID, &tone.Name, &tone.Prompt, &tone.IsSystemDefault, &tone.CreatedAt, &tone.UpdatedAt)
						if err != nil {
							return nil, err
						}
						tones = append(tones, tone)
					}
					return tones, nil
				},
			},
			"tone": &graphql.Field{
				Type: toneType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				// Retrieves a single tone by ID.
				//
				// Arguments:
				//   - id: Tone ID (required)
				//
				// Returns:
				//   - Tone object if found
				//   - error if ID doesn't exist or database issue
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					var tone models.Tone
					err := db.QueryRowContext(p.Context, `
						SELECT id, name, prompt, is_system_default, created_at, updated_at 
						FROM tones WHERE id = $1
					`, id).Scan(&tone.ID, &tone.Name, &tone.Prompt, &tone.IsSystemDefault, &tone.CreatedAt, &tone.UpdatedAt)
					if err != nil {
						return nil, err
					}
					return &tone, nil
				},
			},
		},
	})

	// ========================================================================
	// MUTATION OPERATIONS
	// ========================================================================

	// Define the root mutation with all state-changing operations.
	//
	// Mutation operations modify data and trigger side effects:
	//
	// Dossier Configuration:
	//   - createDossierConfig: Create new configuration
	//   - updateDossierConfig: Update existing configuration
	//   - deleteDossierConfig: Delete configuration
	//
	// Dossier Generation & Delivery:
	//   - generateAndSendDossier: Manually trigger delivery (fetch, summarize, send)
	//   - sendTestEmail: Send test email with sample data
	//   - testEmailConnection: Validate SMTP configuration
	//
	// Tone Management:
	//   - createTone: Create custom tone
	//   - updateTone: Update custom tone (system defaults protected)
	//   - deleteTone: Delete custom tone (system defaults protected)
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createDossierConfig": &graphql.Field{
				Type: dossierConfigType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(dossierConfigInputType),
					},
				},
				// Creates a new dossier configuration for automated delivery.
				//
				// Arguments:
				//   - input: DossierConfigInput with all required fields
				//
				// Default Values Applied:
				//   - tone: "professional" if not specified
				//   - language: "English" if not specified
				//   - specialInstructions: "" (empty) if not specified
				//   - active: true (set by database default)
				//
				// Returns:
				//   - Newly created DossierConfig object with generated ID
				//   - error for validation failures or database issues
				//
				// Side Effects:
				//   - Scheduler will begin monitoring this configuration
				//   - Automated deliveries will start based on frequency and delivery_time
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})

					title := input["title"].(string)
					email := input["email"].(string)
					feedUrls := input["feedUrls"].([]interface{})
					articleCount := input["articleCount"].(int)
					frequency := input["frequency"].(string)
					deliveryTime := input["deliveryTime"].(string)
					timezone := input["timezone"].(string)

					// Handle optional fields with defaults
					tone := "professional"
					if input["tone"] != nil {
						tone = input["tone"].(string)
					}

					language := "English"
					if input["language"] != nil {
						language = input["language"].(string)
					}

					specialInstructions := ""
					if input["specialInstructions"] != nil {
						specialInstructions = input["specialInstructions"].(string)
					}

					// Convert feedUrls from []interface{} to []string
					feedURLStrings := make([]string, len(feedUrls))
					for i, url := range feedUrls {
						feedURLStrings[i] = url.(string)
					}

					var config models.DossierConfig
					err := db.QueryRowContext(p.Context, `
						INSERT INTO dossier_configs (title, email, feed_urls, article_count, frequency, 
							delivery_time, timezone, tone, language, special_instructions)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
						RETURNING id, title, email, feed_urls, article_count, frequency, 
							delivery_time, timezone, tone, language, special_instructions, 
							active, created_at
					`, title, email, pq.Array(feedURLStrings), articleCount, frequency, deliveryTime,
						timezone, tone, language, specialInstructions).Scan(
						&config.ID, &config.Title, &config.Email, pq.Array(&config.FeedURLs),
						&config.ArticleCount, &config.Frequency, &config.DeliveryTime,
						&config.Timezone, &config.Tone, &config.Language,
						&config.SpecialInstructions, &config.Active, &config.CreatedAt)
					if err != nil {
						return nil, err
					}

					log.Printf("Created new dossier config: %s", config.Title)
					return &config, nil
				},
			},
			"updateDossierConfig": &graphql.Field{
				Type: dossierConfigType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(dossierConfigInputType),
					},
				},
				// Updates an existing dossier configuration.
				//
				// Arguments:
				//   - id: Configuration ID to update (required)
				//   - input: DossierConfigInput with updated values
				//
				// Behavior:
				//   - Replaces all fields with new values
				//   - Sets updated_at timestamp automatically
				//   - Applies same default values as createDossierConfig
				//
				// Returns:
				//   - Updated DossierConfig object
				//   - error if ID doesn't exist or validation fails
				//
				// Side Effects:
				//   - Scheduler will use updated settings for next delivery
				//   - No impact on already-sent dossiers
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					input := p.Args["input"].(map[string]interface{})

					title := input["title"].(string)
					email := input["email"].(string)
					feedUrls := input["feedUrls"].([]interface{})
					articleCount := input["articleCount"].(int)
					frequency := input["frequency"].(string)
					deliveryTime := input["deliveryTime"].(string)
					timezone := input["timezone"].(string)

					// Handle optional fields with defaults
					tone := "professional"
					if input["tone"] != nil {
						tone = input["tone"].(string)
					}

					language := "English"
					if input["language"] != nil {
						language = input["language"].(string)
					}

					specialInstructions := ""
					if input["specialInstructions"] != nil {
						specialInstructions = input["specialInstructions"].(string)
					}

					// Convert feedUrls from []interface{} to []string
					feedURLStrings := make([]string, len(feedUrls))
					for i, url := range feedUrls {
						feedURLStrings[i] = url.(string)
					}

					var config models.DossierConfig
					err := db.QueryRowContext(p.Context, `
						UPDATE dossier_configs 
						SET title = $2, email = $3, feed_urls = $4, article_count = $5, 
							frequency = $6, delivery_time = $7, timezone = $8, tone = $9, 
							language = $10, special_instructions = $11, updated_at = CURRENT_TIMESTAMP
						WHERE id = $1
						RETURNING id, title, email, feed_urls, article_count, frequency, 
							delivery_time::text, timezone, tone, language, special_instructions, 
							active, created_at
					`, id, title, email, pq.Array(feedURLStrings), articleCount, frequency, deliveryTime,
						timezone, tone, language, specialInstructions).Scan(

						&config.ID, &config.Title, &config.Email, pq.Array(&config.FeedURLs),
						&config.ArticleCount, &config.Frequency, &config.DeliveryTime,
						&config.Timezone, &config.Tone, &config.Language,
						&config.SpecialInstructions, &config.Active, &config.CreatedAt)
					if err != nil {
						return nil, err
					}

					log.Printf("Updated dossier config: %s", config.Title)
					return &config, nil
				},
			},
			"deleteDossierConfig": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				// Permanently deletes a dossier configuration.
				//
				// Arguments:
				//   - id: Configuration ID to delete (required)
				//
				// Returns:
				//   - true if deletion successful
				//   - false if deletion failed
				//
				// Side Effects:
				//   - Configuration permanently removed from database
				//   - Scheduler stops monitoring this configuration
				//   - Historical deliveries remain in dossier_deliveries table
				//
				// Warning: This is a hard delete, not soft delete. Cannot be undone.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)

					_, err := db.ExecContext(p.Context, "DELETE FROM dossier_configs WHERE id = $1", id)
					if err != nil {
						return false, err
					}

					log.Printf("Deleted dossier config ID: %s", id)
					return true, nil
				},
			},
			"generateAndSendDossier": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"configId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				// Manually generates and sends a dossier immediately.
				//
				// This mutation performs the complete dossier generation pipeline:
				//   1. Fetch configuration from database
				//   2. Fetch articles from configured RSS feeds
				//   3. Generate AI summary using specified tone and language
				//   4. Send email with summary and article links
				//   5. Record delivery in dossier_deliveries table
				//
				// Arguments:
				//   - configId: Configuration ID to process (required)
				//
				// Returns:
				//   - true if entire pipeline succeeds
				//   - false with error message if any step fails
				//
				// Error Conditions:
				//   - Configuration not found or inactive
				//   - No articles found from RSS feeds
				//   - AI summary generation fails
				//   - Email delivery fails
				//
				// Use Cases:
				//   - Testing configuration before enabling automation
				//   - Manual on-demand dossier generation
				//   - Debugging delivery issues
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					configId := p.Args["configId"].(string)

					// Get dossier config
					var config models.DossierConfig
					err := db.QueryRowContext(p.Context, `
						SELECT id, title, email, feed_urls, article_count, frequency, 
							   delivery_time::text, timezone, tone, language, special_instructions, 
							   active, created_at
						FROM dossier_configs WHERE id = $1 AND active = true
					`, configId).Scan(&config.ID, &config.Title, &config.Email, pq.Array(&config.FeedURLs),
						&config.ArticleCount, &config.Frequency, &config.DeliveryTime,
						&config.Timezone, &config.Tone, &config.Language,
						&config.SpecialInstructions, &config.Active, &config.CreatedAt)
					if err != nil {
						if err == sql.ErrNoRows {
							return false, fmt.Errorf("dossier configuration not found or inactive")
						}
						return false, err
					}

					// Fetch articles from RSS feeds
					articles, err := rssService.FetchArticlesFromFeeds(p.Context, config.FeedURLs, config.ArticleCount)
					if err != nil {
						return false, fmt.Errorf("failed to fetch articles: %w", err)
					}

					if len(articles) == 0 {
						return false, fmt.Errorf("no articles found from the configured feeds")
					}

					// Generate AI summary
					summary, err := aiService.GenerateSummary(p.Context, articles, config.Tone, config.Language, config.SpecialInstructions)
					if err != nil {
						return false, fmt.Errorf("failed to generate summary: %w", err)
					}

					// Send email
					err = emailService.SendDossier(&config, summary, articles)
					if err != nil {
						return false, fmt.Errorf("failed to send dossier email: %w", err)
					}

					// Record delivery in database
					_, err = db.ExecContext(p.Context, `
						INSERT INTO dossier_deliveries (config_id, delivery_date, summary, article_count, email_sent)
						VALUES ($1, CURRENT_TIMESTAMP, $2, $3, true)
					`, config.ID, summary, len(articles))
					if err != nil {
						log.Printf("Failed to record dossier delivery: %v", err)
					}

					log.Printf("Successfully generated and sent dossier '%s' to %s", config.Title, config.Email)
					return true, nil
				},
			},
			"sendTestEmail": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"configId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				// Sends a test email with sample data to validate configuration.
				//
				// This mutation sends a test email using the configuration's settings
				// but with hardcoded sample articles instead of real RSS content.
				//
				// Process:
				//   1. Fetch configuration from database
				//   2. Create sample articles with Lorem Ipsum content
				//   3. Build test email with configuration details display
				//   4. Send email using configured SMTP settings
				//   5. Modifies subject line to include "- Test Email" suffix
				//
				// Arguments:
				//   - configId: Configuration ID to test (required)
				//
				// Returns:
				//   - true if test email sent successfully
				//   - false with error message if sending fails
				//
				// Use Cases:
				//   - Validating email address is correct
				//   - Testing SMTP configuration
				//   - Previewing email template formatting
				//   - Verifying delivery settings before enabling automation
				//
				// Note: Does not record delivery in dossier_deliveries table.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					configId := p.Args["configId"].(string)

					// Get dossier config
					var config models.DossierConfig
					err := db.QueryRowContext(p.Context, `
						SELECT id, title, email, feed_urls, article_count, frequency, 
							   delivery_time::text, timezone, tone, language, special_instructions, 
							   active, created_at
						FROM dossier_configs WHERE id = $1
					`, configId).Scan(&config.ID, &config.Title, &config.Email, pq.Array(&config.FeedURLs),
						&config.ArticleCount, &config.Frequency, &config.DeliveryTime,
						&config.Timezone, &config.Tone, &config.Language,
						&config.SpecialInstructions, &config.Active, &config.CreatedAt)
					if err != nil {
						if err == sql.ErrNoRows {
							return false, fmt.Errorf("dossier configuration not found")
						}
						return false, err
					}

					testContent := `This is a test email from your Dossier system.

**Configuration Details:**
- Title: ` + config.Title + `
- Frequency: ` + config.Frequency + `
- Delivery Time: ` + config.DeliveryTime + ` (` + config.Timezone + `)
- Article Count: ` + fmt.Sprintf("%d", config.ArticleCount) + `
- AI Tone: ` + config.Tone + `
- Language: ` + config.Language + `

**RSS Feeds:**`

					for i, feedURL := range config.FeedURLs {
						testContent += fmt.Sprintf("\n%d. %s", i+1, feedURL)
					}

					testContent += `

**Sample Articles:** _(This is test data)_
1. **Breaking News: Technology Advances Continue** - Lorem ipsum dolor sit amet, consectetur adipiscing elit.
2. **Market Update: Economic Trends Show Growth** - Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
3. **Innovation Spotlight: New Developments** - Ut enim ad minim veniam, quis nostrud exercitation ullamco.

---
*This was a test email sent at ` + time.Now().Format("2006-01-02 15:04:05 MST") + `*
*Your actual dossiers will contain real articles from your configured RSS feeds.*`

					// Modify config title to indicate test email
					testConfig := config
					testConfig.Title = config.Title + " - Test Email"

					// Create sample articles for email template rendering
					sampleArticles := []models.Article{
						{
							ID:          1,
							Title:       "Breaking News: Technology Advances Continue",
							Link:        "https://example.com/article1",
							Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
							Author:      "Test Author",
							PublishedAt: time.Now(),
						},
						{
							ID:          2,
							Title:       "Market Update: Economic Trends Show Growth",
							Link:        "https://example.com/article2",
							Description: "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
							Author:      "Test Reporter",
							PublishedAt: time.Now().Add(-1 * time.Hour),
						},
						{
							ID:          3,
							Title:       "Innovation Spotlight: New Developments",
							Link:        "https://example.com/article3",
							Description: "Ut enim ad minim veniam, quis nostrud exercitation ullamco.",
							Author:      "Tech Writer",
							PublishedAt: time.Now().Add(-2 * time.Hour),
						},
					}

					err = emailService.SendDossier(&testConfig, testContent, sampleArticles)
					if err != nil {
						return false, fmt.Errorf("failed to send test email: %w", err)
					}

					log.Printf("Successfully sent test email for dossier '%s' to %s", config.Title, config.Email)
					return true, nil
				},
			},
			"testEmailConnection": &graphql.Field{
				Type: graphql.Boolean,
				// Tests SMTP connection without sending actual email.
				//
				// This mutation validates the email service configuration by:
				//   1. Connecting to SMTP server
				//   2. Performing TLS handshake
				//   3. Authenticating with credentials
				//   4. Disconnecting without sending message
				//
				// Returns:
				//   - true if connection successful
				//   - false with error message if connection fails
				//
				// Use Cases:
				//   - Validating SMTP settings on initial setup
				//   - Troubleshooting email delivery issues
				//   - Verifying credentials after password change
				//
				// Note: Uses environment variables for SMTP configuration.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					err := emailService.TestSMTPConnection()
					if err != nil {
						log.Printf("SMTP test failed: %v", err)
						return false, err
					}
					return true, nil
				},
			},
			"createTone": &graphql.Field{
				Type: toneType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(toneInputType),
					},
				},
				// Creates a new custom AI tone preset.
				//
				// Arguments:
				//   - input: ToneInput with name and prompt
				//
				// Behavior:
				//   - is_system_default automatically set to false
				//   - created_at and updated_at timestamps auto-generated
				//
				// Returns:
				//   - Newly created Tone object with generated ID
				//   - error for validation failures or duplicate names
				//
				// Use Cases:
				//   - Creating organization-specific tones
				//   - Custom tone for specific audience
				//   - Experimental tone variations
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})
					var tone models.Tone

					err := db.QueryRowContext(p.Context, `
						INSERT INTO tones (name, prompt) 
						VALUES ($1, $2) 
						RETURNING id, name, prompt, is_system_default, created_at, updated_at
					`, input["name"], input["prompt"]).Scan(
						&tone.ID, &tone.Name, &tone.Prompt, &tone.IsSystemDefault, &tone.CreatedAt, &tone.UpdatedAt)
					if err != nil {
						return nil, err
					}
					return &tone, nil
				},
			},
			"updateTone": &graphql.Field{
				Type: toneType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(toneInputType),
					},
				},
				// Updates an existing custom tone.
				//
				// Arguments:
				//   - id: Tone ID to update (required)
				//   - input: ToneInput with updated name and prompt
				//
				// Behavior:
				//   - Only updates custom tones (is_system_default = false)
				//   - Sets updated_at timestamp automatically
				//   - Returns error if attempting to modify system default tone
				//
				// Returns:
				//   - Updated Tone object
				//   - error if ID doesn't exist or tone is system default
				//
				// Protection: System default tones cannot be modified.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					input := p.Args["input"].(map[string]interface{})
					var tone models.Tone

					err := db.QueryRowContext(p.Context, `
						UPDATE tones 
						SET name = $1, prompt = $2, updated_at = CURRENT_TIMESTAMP 
						WHERE id = $3 AND is_system_default = false
						RETURNING id, name, prompt, is_system_default, created_at, updated_at
					`, input["name"], input["prompt"], id).Scan(
						&tone.ID, &tone.Name, &tone.Prompt, &tone.IsSystemDefault, &tone.CreatedAt, &tone.UpdatedAt)
					if err != nil {
						return nil, err
					}
					return &tone, nil
				},
			},
			"deleteTone": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				// Deletes a custom tone.
				//
				// Arguments:
				//   - id: Tone ID to delete (required)
				//
				// Behavior:
				//   - Only deletes custom tones (is_system_default = false)
				//   - Returns false if tone is system default
				//
				// Returns:
				//   - true if deletion successful
				//   - false if tone is system default or doesn't exist
				//   - error for database issues
				//
				// Protection: System default tones cannot be deleted.
				//
				// Note: Dossier configurations using this tone should be updated
				// before deletion to avoid reference errors.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)

					result, err := db.ExecContext(p.Context, `
						DELETE FROM tones WHERE id = $1 AND is_system_default = false
					`, id)
					if err != nil {
						return false, err
					}

					rowsAffected, err := result.RowsAffected()
					if err != nil {
						return false, err
					}

					return rowsAffected > 0, nil
				},
			},
		},
	})

	// ========================================================================
	// SCHEMA CREATION
	// ========================================================================

	// Create GraphQL schema with query and mutation roots.
	//
	// The schema defines the complete GraphQL API contract including:
	//   - All queryable types and their relationships
	//   - Available query operations (read-only)
	//   - Available mutation operations (state-changing)
	//
	// Schema validation occurs at creation time, catching type mismatches
	// and invalid references before the API is exposed.
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL schema: %w", err)
	}

	// ========================================================================
	// HTTP HANDLER CONFIGURATION
	// ========================================================================

	// Create GraphQL HTTP handler with development-friendly settings.
	//
	// Configuration:
	//   - Pretty: Formats JSON responses for readability (should be false in production)
	//   - GraphiQL: Enables interactive GraphQL IDE at same endpoint
	//
	// GraphiQL IDE:
	//   - Access at http://localhost:8080/graphql in browser
	//   - Provides documentation explorer, query builder, autocomplete
	//   - Useful for development and API exploration
	//
	// Production Considerations:
	//   - Set Pretty: false to reduce bandwidth
	//   - Set GraphiQL: false to disable IDE in production
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return h, nil
}
