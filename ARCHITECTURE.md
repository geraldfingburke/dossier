# Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Vue.js 3 Frontend                       │   │
│  │  ┌───────────┐  ┌──────────┐  ┌──────────────┐    │   │
│  │  │  Login    │  │  Feeds   │  │   Digests    │    │   │
│  │  │  View     │  │  View    │  │   View       │    │   │
│  │  └───────────┘  └──────────┘  └──────────────┘    │   │
│  │                                                      │   │
│  │  ┌──────────────────────────────────────────┐     │   │
│  │  │         State Management (Store)         │     │   │
│  │  │  - User state                            │     │   │
│  │  │  - Feeds state                           │     │   │
│  │  │  - Articles state                        │     │   │
│  │  │  - Digests state                         │     │   │
│  │  └──────────────────────────────────────────┘     │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ GraphQL over HTTP
                            │ (with JWT auth)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                       API Layer (Go)                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           GraphQL Handler & Resolvers                │   │
│  │  ┌────────────────┐  ┌──────────────────────┐     │   │
│  │  │   Queries      │  │    Mutations         │     │   │
│  │  │  - me          │  │  - register          │     │   │
│  │  │  - feeds       │  │  - login             │     │   │
│  │  │  - articles    │  │  - addFeed           │     │   │
│  │  │  - digests     │  │  - deleteFeed        │     │   │
│  │  └────────────────┘  │  - refreshAllFeeds   │     │   │
│  │                      │  - generateDigest    │     │   │
│  │                      └──────────────────────┘     │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Middleware & Auth                       │   │
│  │  - JWT validation                                    │   │
│  │  - CORS handling                                     │   │
│  │  - Request logging                                   │   │
│  │  - Error recovery                                    │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
            │                    │                    │
            ▼                    ▼                    ▼
┌──────────────────┐  ┌──────────────────┐  ┌─────────────────┐
│   RSS Service    │  │   AI Service     │  │  Auth Service   │
│                  │  │                  │  │                 │
│  - Fetch feeds   │  │  - Summarize     │  │  - Register     │
│  - Parse RSS/    │  │    articles      │  │  - Login        │
│    Atom          │  │  - OpenAI API    │  │  - JWT tokens   │
│  - Save articles │  │    integration   │  │  - Password     │
│  - Schedule      │  │  - Mock fallback │  │    hashing      │
│    updates       │  │                  │  │                 │
└──────────────────┘  └──────────────────┘  └─────────────────┘
            │                    │                    │
            └────────────────────┴────────────────────┘
                              │
                              ▼
            ┌─────────────────────────────────────┐
            │        Database Layer               │
            │         (PostgreSQL)                │
            │  ┌─────────────────────────────┐   │
            │  │  Tables:                    │   │
            │  │  - users                    │   │
            │  │  - feeds                    │   │
            │  │  - articles                 │   │
            │  │  - digests                  │   │
            │  │  - digest_articles          │   │
            │  └─────────────────────────────┘   │
            └─────────────────────────────────────┘
```

## Component Details

### Client Layer (Vue.js 3)

**Technology Stack:**
- Vue.js 3 with Composition API
- Vite for build and development
- GraphQL Request for API communication

**Key Components:**
- `LoginView`: User authentication
- `FeedsView`: RSS feed management
- `ArticlesView`: Article listing
- `DigestsView`: AI-generated digest viewing

**State Management:**
- Reactive store using Vue 3 reactivity
- JWT token persistence in localStorage
- GraphQL client with automatic auth headers

### API Layer (Go)

**Technology Stack:**
- Go 1.24
- Chi router for HTTP handling
- graphql-go for GraphQL implementation
- Custom middleware for auth and CORS

**Key Features:**
- Type-safe GraphQL schema
- JWT-based authentication
- Graceful shutdown
- Request logging
- Error handling

**GraphQL Schema:**
```
Types: User, Feed, Article, Digest, AuthPayload
Queries: me, feeds, articles, digests, latestDigest
Mutations: register, login, addFeed, deleteFeed, 
           refreshAllFeeds, generateDigest
```

### Service Layer

#### RSS Service
- **Purpose**: Fetch and parse RSS feeds
- **Library**: gofeed
- **Functions**:
  - Fetch feeds from URLs
  - Parse RSS/Atom formats
  - Save articles to database
  - Schedule periodic updates

#### AI Service
- **Purpose**: Generate article summaries
- **Library**: go-openai
- **Functions**:
  - Summarize multiple articles into digest
  - Summarize individual articles
  - Mock summaries when API key not configured
- **Model**: GPT-3.5-turbo (configurable)

#### Auth Service
- **Purpose**: User authentication and authorization
- **Libraries**: golang-jwt, bcrypt
- **Functions**:
  - User registration with password hashing
  - User login with JWT generation
  - Token validation and refresh
  - Context-based user identification

### Database Layer (PostgreSQL)

**Schema Design:**

```sql
users
├── id (serial, primary key)
├── email (unique)
├── password (hashed)
├── name
└── timestamps

feeds
├── id (serial, primary key)
├── user_id (foreign key → users)
├── url
├── title
├── description
├── active
└── timestamps

articles
├── id (serial, primary key)
├── feed_id (foreign key → feeds)
├── title
├── link (unique)
├── description
├── content
├── author
├── published_at
└── created_at

digests
├── id (serial, primary key)
├── user_id (foreign key → users)
├── date (unique per user)
├── summary
└── created_at

digest_articles (junction table)
├── digest_id (foreign key → digests)
└── article_id (foreign key → articles)
```

**Indexes:**
- `idx_feeds_user_id`: Fast feed lookup by user
- `idx_articles_feed_id`: Fast article lookup by feed
- `idx_articles_published_at`: Chronological article queries
- `idx_digests_user_id`: Fast digest lookup by user
- `idx_digests_date`: Chronological digest queries

## Data Flow

### Feed Addition Flow
```
1. User enters RSS URL in UI
2. Frontend sends addFeed mutation
3. GraphQL resolver validates auth
4. RSS service fetches and validates feed
5. Feed metadata saved to database
6. Background job fetches initial articles
7. Success response sent to client
8. UI updates feed list
```

### Digest Generation Flow
```
1. User clicks "Generate Digest" or scheduled job runs
2. System fetches all active feeds for user
3. RSS service fetches latest articles (last 24h)
4. Articles saved to database
5. AI service groups and summarizes articles
6. Digest created and linked to articles
7. Summary displayed to user
```

### Authentication Flow
```
1. User submits login credentials
2. Auth service validates email/password
3. JWT token generated (24h expiry)
4. Token sent to client and stored
5. Subsequent requests include token in header
6. Middleware validates token and extracts user ID
7. User ID added to request context
8. Resolvers access user ID from context
```

## Security Considerations

1. **Password Security**: Bcrypt hashing with salt
2. **JWT Tokens**: 24-hour expiration, HMAC-SHA256 signing
3. **SQL Injection**: Prepared statements with parameterized queries
4. **CORS**: Configured for specific origins
5. **API Keys**: Stored in environment variables
6. **Input Validation**: Feed URL validation, email format checks

## Scalability Considerations

### Current Design
- Single server instance
- Direct database connections
- Synchronous article fetching
- In-memory rate limiting

### Future Enhancements
- **Horizontal Scaling**: Load balancer + multiple server instances
- **Caching**: Redis for feed/article caching
- **Queue System**: Background job processing for feed updates
- **CDN**: Static asset delivery
- **Database**: Read replicas for query scaling
- **Rate Limiting**: Redis-based distributed rate limiting

## Monitoring and Observability

### Current Features
- Request logging via Chi middleware
- Error logging to stdout
- Health check endpoint

### Recommended Additions
- Structured logging (JSON format)
- Metrics collection (Prometheus)
- Distributed tracing (OpenTelemetry)
- Error tracking (Sentry)
- Performance monitoring
- Database query analytics

## Deployment Architecture

### Docker Compose (Development)
```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  PostgreSQL  │────▶│  Backend Go  │────▶│  Frontend    │
│  Container   │     │  Container   │     │  Dev Server  │
└──────────────┘     └──────────────┘     └──────────────┘
    Port 5432            Port 8080            Port 5173
```

### Production (Recommended)
```
                    ┌──────────────────┐
                    │   Load Balancer  │
                    │     (Nginx)      │
                    └────────┬─────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
    ┌─────────▼────────┐         ┌─────────▼────────┐
    │  Backend Server  │         │  Backend Server  │
    │    Instance 1    │         │    Instance 2    │
    └─────────┬────────┘         └─────────┬────────┘
              │                             │
              └──────────────┬──────────────┘
                             │
                    ┌────────▼─────────┐
                    │   PostgreSQL     │
                    │  (Managed DB)    │
                    └──────────────────┘

    ┌──────────────────────────────────────┐
    │  Frontend (Static Files on CDN)      │
    └──────────────────────────────────────┘
```

## Technology Choices

### Why Go for Backend?
- Fast compilation and execution
- Excellent concurrency support (goroutines)
- Strong standard library
- Easy deployment (single binary)
- Good GraphQL libraries

### Why GraphQL?
- Flexible data fetching
- Strong typing
- Single endpoint
- Efficient data loading
- Self-documenting API

### Why Vue.js?
- Gentle learning curve
- Reactive and performant
- Composition API for better code organization
- Excellent developer experience with Vite
- Small bundle size

### Why PostgreSQL?
- Robust and reliable
- ACID compliance
- Excellent JSON support
- Rich indexing capabilities
- Great ecosystem and tools

## Performance Considerations

### Backend
- Connection pooling for database
- Goroutines for concurrent feed fetching
- HTTP timeout configuration
- Efficient SQL queries with indexes

### Frontend
- Code splitting with Vite
- Lazy loading of routes
- Efficient state management
- Debounced API calls

### Database
- Proper indexing strategy
- Foreign key constraints
- Efficient JOIN queries
- Regular VACUUM and ANALYZE
