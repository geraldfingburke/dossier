package graphql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/geraldfingburke/dossier/server/internal/ai"
	"github.com/geraldfingburke/dossier/server/internal/auth"
	"github.com/geraldfingburke/dossier/server/internal/models"
	"github.com/geraldfingburke/dossier/server/internal/rss"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Handler creates the GraphQL HTTP handler
func NewHandler(db *sql.DB, rssService *rss.Service, aiService *ai.Service) http.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "development-secret-key-change-in-production"
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production!")
	}
	authService := auth.NewService(jwtSecret)

	// Define GraphQL types
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	feedType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Feed",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"url": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"active": &graphql.Field{
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

	articleType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"feedId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"link": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: graphql.String,
			},
			"publishedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	digestType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Digest",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"userId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"date": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"summary": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"articles": &graphql.Field{
				Type: graphql.NewList(graphql.NewNonNull(articleType)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					digest, ok := p.Source.(*models.Digest)
					if !ok {
						return nil, fmt.Errorf("invalid source type")
					}

					rows, err := db.Query(`
						SELECT a.id, a.feed_id, a.title, a.link, a.description, a.content, a.author, a.published_at, a.created_at
						FROM articles a
						JOIN digest_articles da ON a.id = da.article_id
						WHERE da.digest_id = $1
						ORDER BY a.published_at DESC
					`, digest.ID)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var articles []models.Article
					for rows.Next() {
						var article models.Article
						if err := rows.Scan(&article.ID, &article.FeedID, &article.Title, &article.Link, &article.Description, &article.Content, &article.Author, &article.PublishedAt, &article.CreatedAt); err != nil {
							return nil, err
						}
						articles = append(articles, article)
					}
					return articles, nil
				},
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	authPayloadType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthPayload",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"user": &graphql.Field{
				Type: graphql.NewNonNull(userType),
			},
		},
	})

	// Define Query type
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return nil, fmt.Errorf("unauthorized")
					}

					var user models.User
					err := db.QueryRowContext(p.Context, `
						SELECT id, email, name, created_at, updated_at
						FROM users WHERE id = $1
					`, userID).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
					if err != nil {
						return nil, err
					}
					return user, nil
				},
			},
			"feeds": &graphql.Field{
				Type: graphql.NewList(graphql.NewNonNull(feedType)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return nil, fmt.Errorf("unauthorized")
					}

					rows, err := db.QueryContext(p.Context, `
						SELECT id, user_id, url, title, description, active, created_at, updated_at
						FROM feeds WHERE user_id = $1
						ORDER BY created_at DESC
					`, userID)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var feeds []models.Feed
					for rows.Next() {
						var feed models.Feed
						if err := rows.Scan(&feed.ID, &feed.UserID, &feed.URL, &feed.Title, &feed.Description, &feed.Active, &feed.CreatedAt, &feed.UpdatedAt); err != nil {
							return nil, err
						}
						feeds = append(feeds, feed)
					}
					return feeds, nil
				},
			},
			"articles": &graphql.Field{
				Type: graphql.NewList(graphql.NewNonNull(articleType)),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 50,
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return nil, fmt.Errorf("unauthorized")
					}

					limit := p.Args["limit"].(int)
					offset := p.Args["offset"].(int)

					rows, err := db.QueryContext(p.Context, `
						SELECT a.id, a.feed_id, a.title, a.link, a.description, a.content, a.author, a.published_at, a.created_at
						FROM articles a
						JOIN feeds f ON a.feed_id = f.id
						WHERE f.user_id = $1
						ORDER BY a.published_at DESC
						LIMIT $2 OFFSET $3
					`, userID, limit, offset)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var articles []models.Article
					for rows.Next() {
						var article models.Article
						if err := rows.Scan(&article.ID, &article.FeedID, &article.Title, &article.Link, &article.Description, &article.Content, &article.Author, &article.PublishedAt, &article.CreatedAt); err != nil {
							return nil, err
						}
						articles = append(articles, article)
					}
					return articles, nil
				},
			},
			"digests": &graphql.Field{
				Type: graphql.NewList(graphql.NewNonNull(digestType)),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return nil, fmt.Errorf("unauthorized")
					}

					limit := p.Args["limit"].(int)

					rows, err := db.QueryContext(p.Context, `
						SELECT id, user_id, date, summary, created_at
						FROM digests
						WHERE user_id = $1
						ORDER BY date DESC
						LIMIT $2
					`, userID, limit)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var digests []*models.Digest
					for rows.Next() {
						var digest models.Digest
						if err := rows.Scan(&digest.ID, &digest.UserID, &digest.Date, &digest.Summary, &digest.CreatedAt); err != nil {
							return nil, err
						}
						digests = append(digests, &digest)
					}
					return digests, nil
				},
			},
			"latestDigest": &graphql.Field{
				Type: digestType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return nil, fmt.Errorf("unauthorized")
					}

					var digest models.Digest
					err := db.QueryRowContext(p.Context, `
						SELECT id, user_id, date, summary, created_at
						FROM digests
						WHERE user_id = $1
						ORDER BY date DESC
						LIMIT 1
					`, userID).Scan(&digest.ID, &digest.UserID, &digest.Date, &digest.Summary, &digest.CreatedAt)
					if err != nil {
						if err == sql.ErrNoRows {
							return nil, nil
						}
						return nil, err
					}
					return &digest, nil
				},
			},
		},
	})

	// Define Mutation type
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Type: graphql.NewNonNull(authPayloadType),
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					name := p.Args["name"].(string)

					user, err := authService.Register(p.Context, db, email, password, name)
					if err != nil {
						return nil, err
					}

					token, _, err := authService.Login(p.Context, db, email, password)
					if err != nil {
						return nil, err
					}

					return map[string]interface{}{
						"token": token,
						"user":  user,
					}, nil
				},
			},
			"login": &graphql.Field{
				Type: graphql.NewNonNull(authPayloadType),
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)

					token, user, err := authService.Login(p.Context, db, email, password)
					if err != nil {
						return nil, err
					}

					return map[string]interface{}{
						"token": token,
						"user":  user,
					}, nil
				},
			},
			"addFeed": &graphql.Field{
				Type: graphql.NewNonNull(feedType),
				Args: graphql.FieldConfigArgument{
					"url": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return nil, fmt.Errorf("unauthorized")
					}

					url := p.Args["url"].(string)

					// Try to fetch the feed to validate it
					feed, err := rssService.FetchFeed(p.Context, url)
					if err != nil {
						return nil, fmt.Errorf("invalid feed URL: %w", err)
					}

					// Save feed to database
					var dbFeed models.Feed
					err = db.QueryRowContext(p.Context, `
						INSERT INTO feeds (user_id, url, title, description, active)
						VALUES ($1, $2, $3, $4, true)
						RETURNING id, user_id, url, title, description, active, created_at, updated_at
					`, userID, url, feed.Title, feed.Description).Scan(
						&dbFeed.ID, &dbFeed.UserID, &dbFeed.URL, &dbFeed.Title,
						&dbFeed.Description, &dbFeed.Active, &dbFeed.CreatedAt, &dbFeed.UpdatedAt,
					)
					if err != nil {
						return nil, err
					}

					// Fetch initial articles
					go func() {
						if err := rssService.SaveArticles(db, dbFeed.ID, feed.Items); err != nil {
							fmt.Printf("Error saving initial articles: %v\n", err)
						}
					}()

					return dbFeed, nil
				},
			},
			"deleteFeed": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return false, fmt.Errorf("unauthorized")
					}

					feedID := p.Args["id"].(string)

					result, err := db.ExecContext(p.Context, `
						DELETE FROM feeds WHERE id = $1 AND user_id = $2
					`, feedID, userID)
					if err != nil {
						return false, err
					}

					rows, err := result.RowsAffected()
					if err != nil {
						return false, err
					}

					return rows > 0, nil
				},
			},
			"refreshAllFeeds": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return false, fmt.Errorf("unauthorized")
					}

					go func() {
						ctx := context.Background()
						if err := rssService.FetchAllFeeds(ctx, db, userID); err != nil {
							fmt.Printf("Error refreshing feeds: %v\n", err)
						}
					}()

					return true, nil
				},
			},
			"generateDigest": &graphql.Field{
				Type: graphql.NewNonNull(digestType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, ok := auth.GetUserFromContext(p.Context)
					if !ok {
						return nil, fmt.Errorf("unauthorized")
					}

					if err := rssService.GenerateUserDigest(p.Context, db, userID); err != nil {
						return nil, err
					}

					// Fetch the latest digest
					var digest models.Digest
					err := db.QueryRowContext(p.Context, `
						SELECT id, user_id, date, summary, created_at
						FROM digests
						WHERE user_id = $1
						ORDER BY date DESC
						LIMIT 1
					`, userID).Scan(&digest.ID, &digest.UserID, &digest.Date, &digest.Summary, &digest.CreatedAt)
					if err != nil {
						return nil, err
					}

					return &digest, nil
				},
			},
		},
	})

	// Create schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(err)
	}

	// Create handler with authentication middleware
	return &authMiddleware{
		authService: authService,
		handler: handler.New(&handler.Config{
			Schema:   &schema,
			Pretty:   true,
			GraphiQL: false,
		}),
	}
}

// authMiddleware handles authentication
type authMiddleware struct {
	authService *auth.Service
	handler     *handler.Handler
}

func (m *authMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			userID, err := m.authService.ValidateToken(parts[1])
			if err == nil {
				ctx := context.WithValue(r.Context(), "user_id", userID)
				r = r.WithContext(ctx)
			}
		}
	}

	m.handler.ServeHTTP(w, r)
}

// PlaygroundHandler returns the GraphQL playground handler
func PlaygroundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(playgroundHTML))
	}
}

const playgroundHTML = `
<!DOCTYPE html>
<html>
<head>
  <title>GraphQL Playground</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
  <link rel="shortcut icon" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/favicon.png" />
  <script src="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>
<body>
  <div id="root"></div>
  <script>
    window.addEventListener('load', function (event) {
      GraphQLPlayground.init(document.getElementById('root'), {
        endpoint: '/graphql',
        settings: {
          'request.credentials': 'include',
        }
      })
    })
  </script>
</body>
</html>
`
