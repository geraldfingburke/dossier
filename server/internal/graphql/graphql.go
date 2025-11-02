package graphql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/geraldfingburke/dossier/server/internal/ai"
	"github.com/geraldfingburke/dossier/server/internal/email"
	"github.com/geraldfingburke/dossier/server/internal/models"
	"github.com/geraldfingburke/dossier/server/internal/rss"
	"github.com/geraldfingburke/dossier/server/internal/scheduler"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/lib/pq"
)

// Handler creates the GraphQL HTTP handler
func Handler(db *sql.DB, rssService *rss.Service, aiService *ai.Service, emailService *email.Service, schedulerService *scheduler.Service) (*handler.Handler, error) {
	// DossierConfig GraphQL type
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
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var config *models.DossierConfig
					
					// Handle both pointer and value types
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

	// DossierConfigInput input type
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

	// SchedulerStatus GraphQL type
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

	// Define the root query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"dossierConfigs": &graphql.Field{
				Type: graphql.NewList(dossierConfigType),
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
		},
	})

	// Define the root mutation
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
			"testEmailConnection": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					err := emailService.TestSMTPConnection()
					if err != nil {
						log.Printf("SMTP test failed: %v", err)
						return false, err
					}
					return true, nil
				},
			},
		},
	})

	// Create GraphQL schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL schema: %w", err)
	}

	// Create GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return h, nil
}