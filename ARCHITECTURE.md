# Dossier System Architecture

## High-Level System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Vue.js 3 Frontend                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐    │
│  │  Dossiers   │  │   Feeds     │  │    Articles     │    │
│  │  Management │  │  Management │  │    Browser      │    │
│  └─────────────┘  └─────────────┘  └─────────────────┘    │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Dossier Configuration UI                 │  │
│  │  • Schedule Settings • AI Tone Selection             │  │
│  │  • Email Configuration • RSS Feed Management        │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ GraphQL API
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Go Backend Server                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                GraphQL API Layer                     │   │
│  │  • Dossier CRUD Operations                          │   │
│  │  • Feed Management                                   │   │
│  │  • Article Fetching                                  │   │
│  │  • Manual Test Email Triggers                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Automated Scheduler Service             │   │
│  │  • Cron-like Scheduling Engine                      │   │
│  │  • Timezone-aware Delivery                          │   │
│  │  • Frequency Management (Daily/Weekly/Monthly)      │   │
│  │  • Duplicate Prevention via Database Tracking       │   │
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
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Processing Pipeline                         │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              RSS Feed Processing                     │   │
│  │  • Multi-feed Concurrent Fetching                   │   │
│  │  • Content Parsing & Deduplication                  │   │
│  │  • Article Storage & Metadata Extraction            │   │
│  └─────────────────────────────────────────────────────┘   │
│                            │                                 │
│                            ▼                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │         3-Stage AI Processing Pipeline               │   │
│  │  ┌───────────────┐ ┌────────────────┐ ┌──────────┐ │   │
│  │  │  1. Article   │ │  2. Content    │ │ 3. Summary│ │   │
│  │  │   Selection   │ │   Extraction   │ │Generation │ │   │
│  │  │               │ │                │ │           │ │   │
│  │  │ • Relevance   │ │ • Clean Text   │ │• Tone     │ │   │
│  │  │ • Freshness   │ │ • Remove Fluff │ │• Language │ │   │
│  │  │ • Quality     │ │ • Key Facts    │ │• Format   │ │   │
│  │  └───────────────┘ └────────────────┘ └──────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                            │                                 │
│                            ▼                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Email Generation & Delivery             │   │
│  │  • HTML Template Rendering                          │   │
│  │  • SMTP Secure Transmission                         │   │
│  │  • Delivery Status Tracking                         │   │
│  │  • Timezone-aware Scheduling                        │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
            ┌─────────────────────────────────────┐
            │           Data Layer                │
            │         (PostgreSQL)                │
            │  ┌─────────────────────────────┐   │
            │  │  Core Tables:               │   │
            │  │  - dossiers                 │   │
            │  │  - dossier_feeds            │   │
            │  │  - feeds                    │   │
            │  │  - articles                 │   │
            │  │  - dossier_deliveries       │   │
            │  └─────────────────────────────┘   │
            └─────────────────────────────────────┘
```

## Architectural Principles

### Single-User Focus

- No user authentication system (simplified deployment)
- Direct dossier management without user isolation
- Streamlined API surface without user context

### Local-First AI Processing

- **Ollama Integration**: Local LLM hosting for privacy
- **Multi-Model Support**: Different models for different tones
- **3-Stage Pipeline**: Structured content processing for quality

### Automated Delivery System

- **Scheduler Service**: Cron-like background processing
- **Time-Aware**: Proper timezone handling for global users
- **Frequency Support**: Daily, weekly, monthly delivery options
- **Duplicate Prevention**: Database tracking prevents re-sending

## Component Details

### Frontend (Vue.js 3)

**Technology Stack:**

- Vue.js 3 with Composition API
- Vite for fast development and building
- VeeValidate for form validation
- Responsive design with modern CSS

**Key Views:**

- `DossiersView`: Main dossier configuration interface
- `FeedsView`: RSS feed management and testing
- `ArticlesView`: Browse and search collected articles

**Features:**

- Real-time form validation
- Test email functionality
- Tone selection with uncensored options
- Multi-language dossier support

### Backend (Go)

**Technology Stack:**

- Go 1.21+ with modern idioms
- gqlgen for type-safe GraphQL
- PostgreSQL with prepared statements
- Structured logging with levels

**Core Services:**

#### GraphQL API Layer

- **Schema**: Type-safe dossier, feed, and article operations
- **Resolvers**: CRUD operations for dossier management
- **Mutations**: Create/update/delete dossiers and feeds
- **Queries**: Fetch dossiers, articles, delivery history

#### Scheduler Service

- **Cron Engine**: Background processing for scheduled deliveries
- **Time Management**: Timezone-aware scheduling logic
- **Frequency Logic**: Daily, weekly, monthly delivery patterns
- **Lifecycle**: Integrated with server startup/shutdown

#### RSS Processing Service

- **Feed Fetching**: Concurrent RSS/Atom feed processing
- **Parsing**: gofeed library for robust content extraction
- **Deduplication**: Intelligent article duplicate detection
- **Storage**: Efficient database operations for article management

#### AI Processing Service

- **Ollama Integration**: Local LLM server communication
- **3-Stage Pipeline**:
  1. **selectArticles()**: Intelligent relevance-based article selection
  2. **extractFactualContent()**: Clean content extraction from articles
  3. **generateSummaryFromCleanedArticles()**: Tone-aware summary generation
- **Multi-Model Support**: Standard and uncensored model options
- **Quality Assurance**: Structured prompts for consistent output

#### Email Service

- **SMTP Integration**: Secure email delivery via TLS
- **HTML Templates**: Rich email formatting with Go templates
- **Provider Support**: Gmail, Outlook, and standard SMTP servers
- **Delivery Tracking**: Database logging of send status and timestamps

### Database Layer (PostgreSQL)

**Schema Design:**

```sql
-- Core dossier configuration
dossiers
├── id (serial, primary key)
├── name (e.g., "Tech News", "Sports Updates")
├── delivery_time (TIME, e.g., "08:00:00")
├── frequency (daily/weekly/monthly)
├── timezone (e.g., "America/New_York")
├── tone (professional/humorous/analytical/etc.)
├── language (default: "english")
├── special_instructions (optional custom prompts)
├── email_to (delivery email address)
├── is_active (boolean)
└── timestamps (created_at, updated_at)

-- RSS feeds associated with dossiers
dossier_feeds
├── id (serial, primary key)
├── dossier_id (foreign key → dossiers)
├── feed_id (foreign key → feeds)
└── timestamps

-- RSS feed metadata
feeds
├── id (serial, primary key)
├── url (unique RSS feed URL)
├── title (extracted from feed)
├── description
├── last_fetched_at
└── timestamps

-- Collected articles from all feeds
articles
├── id (serial, primary key)
├── feed_id (foreign key → feeds)
├── title
├── link (unique URL)
├── description
├── content (full article text)
├── author
├── published_at
└── created_at

-- Delivery tracking to prevent duplicates
dossier_deliveries
├── id (serial, primary key)
├── dossier_id (foreign key → dossiers)
├── delivered_at (timestamp with timezone)
├── status (sent/failed)
├── email_content (generated HTML)
└── article_count
```

**Key Indexes:**

- `idx_dossiers_active`: Fast lookup of active dossiers
- `idx_dossier_feeds_dossier`: Feed associations per dossier
- `idx_articles_feed_published`: Recent articles by feed
- `idx_deliveries_dossier_time`: Delivery history tracking
- `idx_feeds_url`: Unique feed URL constraints

## System Flows

### Dossier Creation Flow

```
1. User configures dossier in UI (name, schedule, tone, email)
2. User adds RSS feeds to dossier
3. Frontend sends createDossier mutation with feed URLs
4. Backend validates feeds and fetches initial metadata
5. Dossier and feed associations stored in database
6. Scheduler automatically picks up new dossier for delivery
7. UI confirms successful creation
```

### Automated Delivery Flow

```
1. Scheduler runs every minute checking for due dossiers
2. For each due dossier:
   a. Check delivery history to prevent duplicates
   b. Fetch recent articles from associated RSS feeds
   c. Run 3-stage AI processing pipeline:
      - selectArticles(): Choose most relevant articles
      - extractFactualContent(): Clean and extract key information
      - generateSummaryFromCleanedArticles(): Create formatted summary
   d. Generate HTML email using template
   e. Send via SMTP with delivery tracking
   f. Record successful delivery in database
3. Log all operations with structured logging
```

### Manual Test Flow

```
1. User clicks "Test Email" button in UI
2. Frontend sends testDossier mutation
3. Backend immediately processes dossier (bypassing scheduler)
4. Uses same AI pipeline as automated delivery
5. Sends test email with "TEST" subject prefix
6. Returns success/failure status to UI
7. Does not record in delivery history
```

## AI Processing Pipeline

### Stage 1: Article Selection (`selectArticles()`)

- **Input**: All articles from dossier's RSS feeds (last 24-48 hours)
- **Process**: AI evaluates articles for relevance, freshness, and quality
- **Output**: Filtered list of 10-20 most important articles
- **Criteria**: Topic relevance, source credibility, recency, uniqueness

### Stage 2: Content Extraction (`extractFactualContent()`)

- **Input**: Selected articles with full content
- **Process**: AI extracts key facts, removes promotional content and fluff
- **Output**: Clean, factual summaries of each article
- **Focus**: Core information, key statistics, important quotes, actionable insights

### Stage 3: Summary Generation (`generateSummaryFromCleanedArticles()`)

- **Input**: Clean factual content from stage 2
- **Process**: AI synthesizes information into cohesive summary with chosen tone
- **Output**: Final formatted summary ready for email delivery
- **Features**: Tone adaptation, language selection, custom instructions

## Security Architecture

### Single-User Design Benefits

- **No Authentication Attack Surface**: No user accounts, passwords, or session management
- **Simplified Deployment**: Direct configuration without user isolation concerns
- **Reduced Complexity**: No authorization logic or user-specific data access controls

### Core Security Measures

1. **Local AI Processing**: No external API calls, all data stays on premises
2. **Environment Variables**: All sensitive data (SMTP credentials, secrets) in env files
3. **SQL Injection Prevention**: Prepared statements with parameterized queries
4. **SMTP Security**: TLS encryption for email transmission
5. **Input Validation**: URL validation for RSS feeds, email format verification
6. **Secret Management**: Git history cleaned of exposed credentials

### Infrastructure Security

- **Container Isolation**: Each service runs in isolated Docker containers
- **Network Security**: Services communicate over Docker internal networks
- **Credential Rotation**: Easy SMTP password updates via environment variables
- **Access Control**: File system permissions for database and configuration files

## Deployment Architecture

### Development Environment (Docker Compose)

```
┌──────────────┐   ┌─────────────┐   ┌──────────────┐   ┌─────────────┐
│ PostgreSQL   │   │ Ollama      │   │ Go Backend   │   │ Vue.js      │
│ Database     │───│ AI Server   │───│ API Server   │───│ Frontend    │
│ Port: 5432   │   │ Port: 11434 │   │ Port: 8080   │   │ Port: 5173  │
└──────────────┘   └─────────────┘   └──────────────┘   └─────────────┘
```

### Single-Server Production Deployment

```
                    ┌──────────────────────────────────────┐
                    │          Production Server           │
                    │                                      │
                    │  ┌─────────────┐  ┌─────────────┐   │
                    │  │ Nginx Proxy │  │ Go Backend  │   │
                    │  │ (Port 80/   │──│ (Port 8080) │   │
                    │  │  443)       │  │ + Scheduler │   │
                    │  └─────────────┘  └─────────────┘   │
                    │         │                           │
                    │  ┌──────▼──────┐  ┌─────────────┐   │
                    │  │ PostgreSQL  │  │ Ollama AI   │   │
                    │  │ (Port 5432) │  │ (Port 11434)│   │
                    │  └─────────────┘  └─────────────┘   │
                    └──────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   SMTP Provider   │
                    │ (Gmail/Outlook)   │
                    └──────────────────┘
```

### Cloud Deployment Options

**Option 1: VPS Deployment**

- Single virtual server (4GB RAM minimum for Ollama)
- Docker Compose for all services
- Automated backups for PostgreSQL
- SSL certificate via Let's Encrypt

**Option 2: Container Platform**

- Deploy to platforms like Railway, Render, or DigitalOcean App Platform
- External PostgreSQL service
- Ollama on dedicated AI server or cloud AI service
- Static frontend deployment

## Performance & Scaling

### Current Optimizations

- **Concurrent RSS Fetching**: Goroutines for parallel feed processing
- **Database Connection Pooling**: Efficient PostgreSQL connections
- **Local AI Processing**: No external API latency
- **Smart Scheduling**: Prevents duplicate deliveries
- **Efficient Queries**: Proper database indexing

### Scaling Strategies

- **Vertical Scaling**: More CPU/RAM for better Ollama performance
- **AI Model Optimization**: Smaller models for faster processing
- **Database Optimization**: Query optimization and proper indexing
- **Caching Layer**: Redis for feed caching (future enhancement)
- **Content Delivery**: Static asset CDN for frontend

## Technology Decisions

### Why Go for Backend?

- **Concurrency**: Goroutines perfect for RSS fetching and scheduling
- **Performance**: Fast execution for real-time email generation
- **Deployment**: Single binary deployment for easy server management
- **Libraries**: Excellent GraphQL, HTTP, and database libraries
- **Memory**: Efficient memory usage for always-running scheduler

### Why Local AI (Ollama)?

- **Privacy**: No external API calls, complete data sovereignty
- **Cost**: No per-request charges, unlimited processing
- **Reliability**: No network dependencies for AI processing
- **Customization**: Multiple models, uncensored options
- **Speed**: Local processing after initial model download

### Why GraphQL?

- **Type Safety**: Compile-time schema validation
- **Flexibility**: Frontend can request exactly needed data
- **Single Endpoint**: Simplified API surface
- **Development**: Excellent tooling and introspection
- **Evolution**: Easy schema changes without versioning

### Why Vue.js 3?

- **Reactivity**: Perfect for real-time dossier configuration
- **Composition API**: Clean component logic organization
- **Performance**: Efficient rendering with proxy-based reactivity
- **Developer Experience**: Excellent Vite integration and DevTools
- **Ecosystem**: Rich component libraries and tooling

### Why PostgreSQL?

- **Reliability**: ACID compliance for delivery tracking
- **Time Support**: Native TIME/TIMESTAMP handling for scheduling
- **JSON**: Flexible configuration storage capabilities
- **Performance**: Excellent indexing for article queries
- **Ecosystem**: Mature tooling and backup solutions

### Why Docker?

- **Consistency**: Identical development and production environments
- **Isolation**: Each service in dedicated container
- **Dependencies**: AI models and database bundled cleanly
- **Deployment**: Simple docker-compose deployment
- **Scaling**: Easy horizontal scaling when needed

## Monitoring Strategy

### Current Implementation

- **Structured Logging**: JSON logs with levels (debug/info/warn/error)
- **Scheduler Logging**: Detailed logs for delivery processing
- **AI Pipeline Logging**: Debug information for content processing
- **Email Tracking**: Database records of all delivery attempts
- **Error Handling**: Graceful degradation with detailed error messages

### Production Monitoring (Recommended)

- **Health Checks**: `/health` endpoint for uptime monitoring
- **Metrics Collection**: Prometheus for system metrics
- **Log Aggregation**: ELK stack or similar for log analysis
- **Email Delivery Monitoring**: Track bounce rates and delivery success
- **AI Model Performance**: Response times and quality metrics
- **Database Performance**: Query times and connection pool usage
