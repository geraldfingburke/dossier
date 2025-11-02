# Dossier System Architecture

## High-Level System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Vue.js 3 Frontend                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐    │
│  │  Dossier    │  │   Digest    │  │    Articles     │    │
│  │  Configs    │  │  History    │  │    View         │    │
│  └─────────────┘  └─────────────┘  └─────────────────┘    │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Dossier Configuration UI                 │  │
│  │  • Schedule Settings (Time/Frequency/Timezone)       │  │
│  │  • AI Tone Selection (10 defaults + custom)          │  │
│  │  • Email Configuration (Delivery address)            │  │
│  │  • RSS Feed URLs (Inline management)                 │  │
│  │  • Language & Special Instructions                   │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ GraphQL API (http://localhost:8080/graphql)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Go Backend Server                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                GraphQL API Layer                     │   │
│  │  • Dossier Config CRUD (create/read/update/delete)  │   │
│  │  • Tone Management (system + custom)                │   │
│  │  • Delivery History (dossiers query)                │   │
│  │  • Manual Triggers (generateAndSendDossier)         │   │
│  │  • Test Email (sendTestEmail)                       │   │
│  │  • Scheduler Status (running/nextCheck/active)      │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Automated Scheduler Service             │   │
│  │  • Ticker-based (1-minute granularity)              │   │
│  │  • Timezone-aware Delivery                          │   │
│  │  • Frequency Support (Daily/Weekly/Monthly)         │   │
│  │  • Duplicate Prevention (last_generated tracking)   │   │
│  │  • Async Processing (goroutine per dossier)         │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
            │                    │                    │
            ▼                    ▼                    ▼
┌──────────────────┐  ┌──────────────────┐  ┌─────────────────┐
│   RSS Service    │  │   AI Service     │  │  Email Service  │
│                  │  │                  │  │                 │
│  - Fetch feeds   │  │  - Summarize     │  │  - SMTP/TLS     │
│  - Parse RSS/    │  │    articles      │  │  - HTML email   │
│    Atom/RSS 2.0  │  │  - Ollama LLM    │  │  - Delivery     │
│  - Multi-feed    │  │  - Multi-model   │  │    tracking     │
│    aggregation   │  │  - Tone support  │  │  - Test email   │
│    aggregation   │  │  - Tone support  │  │    tracking     │
│  - Deduplication │  │  - Language      │  │                 │
└──────────────────┘  └──────────────────┘  └─────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Processing Pipeline                         │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              RSS Feed Processing                     │   │
│  │  • Multi-feed Concurrent Fetching (goroutines)      │   │
│  │  • Content Parsing (RSS 1.0/2.0, Atom 1.0)         │   │
│  │  • Article Storage & Deduplication (by link URL)    │   │
│  │  • Missing Field Handling (graceful degradation)    │   │
│  └─────────────────────────────────────────────────────┘   │
│                            │                                 │
│                            ▼                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           AI Processing Pipeline (Ollama)            │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │  1. Article Aggregation                       │  │   │
│  │  │     • Fetch from all feed URLs                │  │   │
│  │  │     • Sort by published date                  │  │   │
│  │  │     • Limit to articleCount (configurable)    │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │  2. Content Formatting                        │  │   │
│  │  │     • Format: "Title + Description + Link"    │  │   │
│  │  │     • Prevent token overflow (15 articles)    │  │   │
│  │  │     • Clean text for LLM processing           │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │  3. AI Summary Generation                     │  │   │
│  │  │     • Apply tone-specific system prompt       │  │   │
│  │  │     • Apply language preference               │  │   │
│  │  │     • Apply special instructions              │  │   │
│  │  │     • Model: llama3.2:3b / dolphin-mistral    │  │   │
│  │  │     • Output: Markdown-formatted summary      │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│                            │                                 │
│                            ▼                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Email Generation & Delivery             │   │
│  │  • Markdown to HTML Conversion                      │   │
│  │  • HTML Template Rendering (professional layout)    │   │
│  │  • SMTP/TLS Secure Transmission                     │   │
│  │  • Delivery Status Tracking (dossier_deliveries)    │   │
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
            │  │  - dossier_configs          │   │
            │  │  - feeds                    │   │
            │  │  - articles                 │   │
            │  │  - dossier_deliveries       │   │
            │  │  - delivery_articles        │   │
            │  │  - tones                    │   │
            │  └─────────────────────────────┘   │
            └─────────────────────────────────────┘
```

## Architectural Principles

### Single-User Focus

- **No authentication system**: Simplified deployment without user accounts, passwords, or session management
- **Direct dossier management**: All configs accessible without user context or ownership
- **Streamlined API**: No authorization checks, user filtering, or multi-tenant complexity
- **Personal deployment**: Designed for individual use on private servers or localhost

### Local-First AI Processing

- **Ollama Integration**: Local LLM hosting for complete privacy and data sovereignty
- **Multi-Model Support**:
  - `llama3.2:3b` (default) - Fast, efficient, general-purpose
  - `dolphin-mistral` (uncensored) - For "sweary" tone without content filters
- **Zero External APIs**: All AI processing happens locally, no OpenAI/Anthropic calls
- **Cost-Free Operation**: Unlimited summaries without per-request charges

### Automated Delivery System

- **Scheduler Service**: Ticker-based (1-minute granularity) background processing
- **Timezone-Aware**: Converts delivery times to UTC for accurate global scheduling
- **Frequency Support**: Daily, weekly, and monthly delivery patterns with intelligent date handling
- **Duplicate Prevention**: Tracks `last_generated` timestamp to prevent re-sending

### Data Simplicity

- **Feed URLs in Config**: No separate feed management - URLs stored as TEXT[] array in dossier_configs
- **Shared Articles**: Single articles table, deduplicated by link URL across all dossiers
- **Delivery Archive**: Full summary content and article list stored in dossier_deliveries
- **Customizable Tones**: 10 system defaults + unlimited custom user-defined tones

## Component Details

### Frontend (Vue.js 3)

**Technology Stack:**

- Vue.js 3 with Composition API and `<script setup>`
- Vite for fast development and optimized builds
- Vuex for state management (GraphQL data)
- Modular CSS architecture (9 separate files)
- Responsive design with modern CSS Grid/Flexbox

**Key Views:**

1. **DossierConfigsView** (`/`):

   - Main configuration interface for creating/editing dossiers
   - Form inputs for title, email, feed URLs, schedule, tone, language
   - Active/inactive toggle for temporary disabling
   - Test email button for immediate preview
   - Delete confirmation dialogs

2. **DigestsView** (`/digests`):

   - Historical view of sent dossiers
   - Filter by dossier config
   - View full HTML content of past deliveries
   - Delivery timestamps and article counts

3. **ArticlesView** (`/articles`):
   - Browse all fetched articles from RSS feeds
   - Display title, description, author, published date
   - Links to original article sources

**State Management:**

- Vuex store for dossier configs, tones, and delivery history
- GraphQL query caching and optimistic updates
- Local state for form validation and UI interactions

**Styling:**

- Modular CSS: `variables.css`, `reset.css`, `buttons.css`, `forms.css`, `components.css`, `layout.css`, `modals.css`, `utilities.css`, `main.css`
- Design tokens for consistent spacing, colors, typography
- Dark mode support (via CSS variables)

### Backend (Go)

**Technology Stack:**

- Go 1.21+ with standard library idioms
- `github.com/graphql-go/graphql` for GraphQL server
- `github.com/lib/pq` for PostgreSQL driver
- `github.com/mmcdole/gofeed` for RSS/Atom parsing
- Structured logging to stdout (JSON format)

**Project Structure:**

```
server/
├── cmd/
│   └── main.go              # Entry point, server initialization
└── internal/
    ├── ai/
    │   └── ai.go            # Ollama LLM integration
    ├── auth/
    │   └── auth.go          # (Legacy - not used)
    ├── database/
    │   └── database.go      # PostgreSQL connection & migrations
    ├── email/
    │   └── email.go         # SMTP email delivery
    ├── graphql/
    │   ├── graphql.go       # GraphQL schema & resolvers
    │   └── schema.graphql   # GraphQL type definitions
    ├── models/
    │   └── models.go        # Domain types (DossierConfig, Article, etc.)
    ├── rss/
    │   └── rss.go           # RSS feed fetching & parsing
    └── scheduler/
        └── scheduler.go     # Automated delivery scheduler
```

**Core Services:**

#### GraphQL API Layer (`internal/graphql/graphql.go`)

- **Schema**: Type-safe dossier config, tone, and delivery operations
- **Queries**:
  - `dossierConfigs` - List all configurations
  - `dossierConfig(id)` - Get specific config
  - `dossiers(configId, limit)` - Delivery history
  - `schedulerStatus` - Scheduler state (running/nextCheck/activeDossiers)
  - `tones` - List all AI tones
  - `tone(id)` - Get specific tone
- **Mutations**:
  - `createDossierConfig(input)` - Create new configuration
  - `updateDossierConfig(id, input)` - Update existing config
  - `deleteDossierConfig(id)` - Delete configuration
  - `toggleDossierConfig(id, active)` - Enable/disable
  - `generateAndSendDossier(configId)` - Manual trigger
  - `sendTestEmail(configId)` - Test email delivery
  - `testEmailConnection(...)` - Test SMTP credentials
  - `createTone(input)` - Create custom tone
  - `updateTone(id, input)` - Update custom tone
  - `deleteTone(id)` - Delete custom tone (system tones protected)

#### Scheduler Service (`internal/scheduler/scheduler.go`)

- **Ticker-Based**: 1-minute interval check for due dossiers
- **Timezone Handling**: Converts delivery_time + timezone to UTC for comparison
- **Frequency Logic**:
  - **Daily**: Generates if current time matches delivery_time
  - **Weekly**: Generates if 7+ days since last delivery
  - **Monthly**: Generates if 30+ days since last delivery
- **Duplicate Prevention**: Queries `dossier_deliveries` for last delivery timestamp
- **Async Processing**: Each dossier processed in separate goroutine
- **Lifecycle**: Started on server boot, graceful shutdown on SIGTERM

#### RSS Processing Service (`internal/rss/rss.go`)

- **Multi-Feed Fetching**: Concurrent fetching with goroutines
- **Parsing**: Supports RSS 1.0, RSS 2.0, Atom 1.0 via `gofeed` library
- **Deduplication**: Uses article link URL as unique identifier
- **Error Handling**: Individual feed failures don't stop entire process
- **Data Quality**: Handles missing fields (author, description, content) gracefully
- **Functions**:
  - `NewService(db)` - Initialize RSS service
  - `FetchFeed(url)` - Fetch single feed and return metadata
  - `FetchArticlesFromFeeds(feedUrls, limit)` - Aggregate articles from multiple feeds

#### AI Processing Service (`internal/ai/ai.go`)

- **Ollama Integration**: HTTP client for local Ollama server
- **Model Selection**:
  - Default: `llama3.2:3b` (fast, balanced)
  - Uncensored: `dolphin-mistral` (for "sweary" tone)
- **Processing Pipeline**:
  1. Format articles (title + description + link)
  2. Apply tone-specific system prompt from `tones` table
  3. Apply language preference (e.g., "english", "spanish")
  4. Apply special instructions (custom user prompts)
  5. Generate summary via Ollama `/api/generate` endpoint
  6. Return markdown-formatted summary
- **Configuration**:
  - Temperature: 0.7 (balanced creativity)
  - Max Tokens: 2000 (comprehensive summaries)
  - Stream: false (wait for complete response)
- **Error Handling**: Detailed logging, graceful fallback on failure

#### Email Service (`internal/email/email.go`)

- **SMTP Integration**: TLS encryption via `smtp.SendMail`
- **HTML Templates**: Professional email layout with inline CSS
- **Configuration**: Environment variables (SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS, SMTP_FROM)
- **Provider Support**: Gmail, Outlook, generic SMTP servers
- **Delivery Tracking**: Records status in `dossier_deliveries` table
- **Test Mode**: `sendTestEmail` for immediate preview without recording history

### Database Layer (PostgreSQL)

**Schema Design:**

```sql
-- ========================================================================
-- TABLE: dossier_configs
-- ========================================================================
-- Core dossier configuration table (single-user, no user_id foreign key)
-- Each row represents one automated digest configuration

CREATE TABLE dossier_configs (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,                  -- Display name (e.g., "Tech News Daily")
  email VARCHAR(255) NOT NULL,                  -- Delivery email address
  feed_urls TEXT[] NOT NULL,                    -- RSS/Atom feed URLs (array)
  article_count INTEGER DEFAULT 20              -- Number of articles per digest
    CHECK (article_count >= 1 AND article_count <= 50),
  frequency VARCHAR(50) NOT NULL                -- 'daily', 'weekly', 'monthly'
    CHECK (frequency IN ('daily', 'weekly', 'monthly')),
  delivery_time TIME NOT NULL,                  -- HH:MM:SS time of day
  timezone VARCHAR(50) DEFAULT 'UTC',           -- IANA timezone (e.g., 'America/New_York')
  tone VARCHAR(50) DEFAULT 'professional',      -- AI tone name (references tones.name)
  language VARCHAR(50) DEFAULT 'English',       -- Summary language
  special_instructions TEXT DEFAULT '',         -- Custom AI instructions
  active BOOLEAN DEFAULT true,                  -- Enable/disable without deletion
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================================
-- TABLE: feeds
-- ========================================================================
-- RSS feed metadata cache (NOT directly linked to dossier_configs)
-- Used for feed validation and metadata display

CREATE TABLE feeds (
  id SERIAL PRIMARY KEY,
  url TEXT NOT NULL UNIQUE,                     -- RSS/Atom feed URL
  title VARCHAR(255),                           -- Extracted from feed metadata
  description TEXT,                             -- Feed description
  active BOOLEAN DEFAULT true,                  -- Can disable problematic feeds
  last_fetched TIMESTAMP,                       -- Track fetch schedule
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================================
-- TABLE: articles
-- ========================================================================
-- Fetched articles from all RSS feeds (shared across dossiers)
-- Deduplication by link URL

CREATE TABLE articles (
  id SERIAL PRIMARY KEY,
  feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
  title TEXT NOT NULL,                          -- Article headline
  link TEXT NOT NULL UNIQUE,                    -- Article URL (unique constraint)
  description TEXT,                             -- Article summary/excerpt
  content TEXT,                                 -- Full article content
  author VARCHAR(255),                          -- Article author
  published_at TIMESTAMP NOT NULL,              -- Original publication date
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================================
-- TABLE: dossier_deliveries
-- ========================================================================
-- Historical records of sent dossiers
-- Archives full summary content for viewing

CREATE TABLE dossier_deliveries (
  id SERIAL PRIMARY KEY,
  config_id INTEGER REFERENCES dossier_configs(id) ON DELETE CASCADE,
  delivery_date TIMESTAMP NOT NULL,             -- When dossier was sent
  summary TEXT NOT NULL,                        -- Generated AI summary (stored)
  article_count INTEGER NOT NULL,               -- Number of articles included
  email_sent BOOLEAN DEFAULT false,             -- Delivery status
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================================
-- TABLE: delivery_articles
-- ========================================================================
-- Junction table: Which articles were in which deliveries
-- Enables "already sent" detection and archive viewing

CREATE TABLE delivery_articles (
  delivery_id INTEGER REFERENCES dossier_deliveries(id) ON DELETE CASCADE,
  article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
  PRIMARY KEY (delivery_id, article_id)
);

-- ========================================================================
-- TABLE: tones
-- ========================================================================
-- AI writing styles (10 system defaults + custom user tones)
-- Each tone contains a prompt that guides LLM behavior

CREATE TABLE tones (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,            -- Tone identifier (e.g., "professional")
  prompt TEXT NOT NULL,                         -- System prompt for AI model
  is_system_default BOOLEAN DEFAULT false,      -- System vs user-created
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Key Indexes:**

```sql
-- Feed lookup optimization
CREATE INDEX idx_feeds_url ON feeds(url);

-- Article queries by feed relationship
CREATE INDEX idx_articles_feed_id ON articles(feed_id);

-- Article sorting by publication date
CREATE INDEX idx_articles_published_at ON articles(published_at DESC);

-- Active dossier filtering for scheduler
CREATE INDEX idx_dossier_configs_active ON dossier_configs(active);

-- Delivery history queries
CREATE INDEX idx_dossier_deliveries_config_id ON dossier_deliveries(config_id);
CREATE INDEX idx_dossier_deliveries_delivery_date ON dossier_deliveries(delivery_date DESC);

-- Tone lookup optimization
CREATE INDEX idx_tones_name ON tones(name);
```

**System Default Tones:**

The database is seeded with 10 default tones during migration:

1. **professional** - Standard business communication
2. **humorous** - Witty and entertaining
3. **analytical** - Data-driven insights
4. **casual** - Relaxed, conversational
5. **apocalyptic** - Dramatic, foreboding with biblical references
6. **orc** - Warcraft-style blunt ("Me Grognak!")
7. **robot** - Mechanical, technical language
8. **southern_belle** - Polite, charming Southern style
9. **apologetic** - Sympathetic and reassuring
10. **sweary** - Adult language (requires uncensored model)

## System Flows

### Dossier Creation Flow

```
1. User fills out form in DossierConfigsView:
   - Title, email, feed URLs (multi-line textarea)
   - Article count, frequency, delivery time, timezone
   - Tone selection (dropdown from tones table)
   - Language and special instructions (optional)

2. Frontend validates inputs:
   - Email format validation
   - Feed URL format validation
   - Time format (HH:MM)
   - Timezone validation (IANA format)

3. Frontend sends createDossierConfig GraphQL mutation

4. Backend processes:
   - Validates all inputs
   - Fetches feed metadata (validates feed URLs)
   - Inserts row into dossier_configs table
   - Creates feed records if they don't exist
   - Returns created config to frontend

5. Scheduler automatically picks up new config on next tick (within 1 minute)

6. UI updates list of dossier configs
```

### Automated Delivery Flow (Scheduler)

```
1. Scheduler tick (every 1 minute):
   - Queries dossier_configs WHERE active = true
   - For each active config:

2. Check if dossier is due:
   a. Convert delivery_time + timezone to UTC
   b. Compare with current UTC time
   c. Check frequency rules:
      - Daily: If current time matches delivery_time
      - Weekly: If 7+ days since last delivery
      - Monthly: If 30+ days since last delivery
   d. Query dossier_deliveries for last delivery timestamp
   e. Skip if already delivered today/this period

3. If dossier is due, process in goroutine:
   a. Fetch articles from all feed_urls (RSS service)
   b. Sort by published_at DESC, limit to article_count
   c. Format articles for AI (title + description + link)
   d. Generate summary via Ollama (AI service):
      - Apply tone prompt from tones table
      - Apply language preference
      - Apply special instructions
      - Generate markdown summary
   e. Convert markdown to HTML
   f. Send email via SMTP (Email service)
   g. Record delivery in dossier_deliveries table
   h. Insert article mappings in delivery_articles table

4. Log all operations (success/failure) to stdout

5. Continue to next config (individual failures don't stop scheduler)
```

### Manual Test Flow

```
1. User clicks "Send Test Email" button in DossierConfigsView

2. Frontend sends sendTestEmail(configId) GraphQL mutation

3. Backend immediately processes (bypasses scheduler timing):
   a. Loads dossier config from database
   b. Fetches articles from feed_urls
   c. Generates summary via AI pipeline (same as automated)
   d. Sends email with "TEST:" prefix in subject
   e. Does NOT record in dossier_deliveries (test only)

4. Backend returns success/failure boolean

5. UI displays success toast or error message

6. User receives test email at configured address
```

### Custom Tone Creation Flow

```
1. User creates custom tone in UI:
   - Enters unique tone name (e.g., "pirate")
   - Writes custom prompt (e.g., "Write like a pirate with 'arr' and 'matey'")

2. Frontend sends createTone mutation

3. Backend validates:
   - Name uniqueness (no duplicate tone names)
   - Prompt length (must be non-empty)
   - Inserts into tones table with is_system_default = false

4. New tone immediately available in dossier config dropdowns

5. Users can update/delete custom tones (system defaults protected)
```

## AI Processing Pipeline

### Overview

The AI processing is handled by a single-stage pipeline that formats articles and generates a summary using the Ollama LLM service.

### Pipeline: Summary Generation

**Input:**

- Array of Article objects from RSS feeds
- DossierConfig with tone, language, special_instructions

**Process:**

1. **Article Formatting:**

   ```
   For each article (up to 15 to prevent token overflow):
     Format: "N. **Title**\n   Description\n   Link: URL\n\n"
   ```

2. **Prompt Construction:**

   ```
   - Base prompt: "Please provide a concise daily digest summary..."
   - Apply tone-specific system prompt from tones.prompt
   - Apply language preference: "Write in {language}"
   - Apply special_instructions if provided
   ```

3. **Ollama API Call:**

   ```
   POST http://localhost:11434/api/generate
   {
     "model": "llama3.2:3b" (or "dolphin-mistral" for sweary tone),
     "prompt": "<constructed prompt>",
     "stream": false,
     "options": {
       "temperature": 0.7,
       "num_predict": 2000
     }
   }
   ```

4. **Response Handling:**
   - Parse JSON response
   - Extract `response` field (markdown-formatted summary)
   - Handle errors gracefully (return error, log details)

**Output:**

- Markdown-formatted summary string
- Ready for HTML conversion and email delivery

**Model Selection:**

- Default: `llama3.2:3b` (3 billion parameters, fast, balanced)
- Uncensored: `dolphin-mistral` (for "sweary" tone, no content filters)

**Tone Application:**

Each tone in the `tones` table has a unique system prompt:

- **professional**: "You are a professional news analyst..."
- **humorous**: "You are a witty comedian summarizing news..."
- **sweary**: "You are an uncensored news commentator..." (uses dolphin-mistral)

The tone prompt is prepended to the article content before sending to Ollama.

## Security Architecture

### Single-User Design Benefits

- **No Authentication Attack Surface**: Zero user accounts, passwords, sessions, or JWT tokens
- **No Authorization Logic**: No complex permission checks, role-based access, or data isolation
- **Simplified Deployment**: Direct configuration without user management overhead
- **Reduced Complexity**: Fewer moving parts = fewer security vulnerabilities
- **Personal Use**: Designed for private servers or localhost development

### Core Security Measures

1. **Local AI Processing**:

   - No external API calls to OpenAI, Anthropic, etc.
   - All data stays on premises
   - No API keys to leak or manage

2. **Environment Variables**:

   - SMTP credentials in .env files
   - Database connection strings
   - Ollama service URLs
   - Never committed to Git

3. **SQL Injection Prevention**:

   - Prepared statements with parameterized queries
   - No string concatenation for SQL
   - PostgreSQL driver handles escaping

4. **SMTP Security**:

   - TLS encryption for email transmission
   - SMTP authentication required
   - Port 587 (STARTTLS) recommended

5. **Input Validation**:

   - URL format validation for RSS feeds
   - Email format validation (regex)
   - Timezone validation (IANA database)
   - Frequency enum validation (daily/weekly/monthly)

6. **Docker Isolation**:
   - Each service in separate container
   - Internal Docker network for service communication
   - Only necessary ports exposed (5173, 8080)

### Infrastructure Security

- **Container Isolation**: Services communicate over Docker internal networks only
- **Network Security**: PostgreSQL port 5432 NOT exposed to host (internal only)
- **Credential Rotation**: Easy SMTP password updates via environment variables
- **Access Control**: File system permissions for database volumes
- **HTTPS**: Nginx reverse proxy with Let's Encrypt SSL in production

## Deployment Architecture

### Development Environment (Docker Compose)

```
┌──────────────┐   ┌─────────────┐   ┌──────────────┐   ┌─────────────┐
│ PostgreSQL   │   │ Ollama      │   │ Go Backend   │   │ Vite Dev    │
│ Database     │───│ AI Server   │───│ + Scheduler  │───│ Server      │
│ Port: 5432   │   │ Port: 11434 │   │ Port: 8080   │   │ Port: 5173  │
│ (internal)   │   │             │   │ /graphql     │   │ Hot Reload  │
└──────────────┘   └─────────────┘   └──────────────┘   └─────────────┘
        │                  │                  │                  │
        └──────────────────┴──────────────────┴──────────────────┘
                            Docker Network
```

**Services:**

- **postgres**: PostgreSQL 15 with persistent volume (`postgres_data`)
- **ollama**: Ollama AI server with models volume (`ollama_data`)
- **backend**: Go server (auto-rebuild on code changes)
- **frontend**: Vite dev server (hot module replacement)

**Commands:**

```bash
# Start all services
docker-compose up

# Build and start (after code changes)
docker-compose up --build

# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down --volumes
```

### Single-Server Production Deployment

```
                    ┌──────────────────────────────────────┐
                    │          Production Server           │
                    │                                      │
                    │  ┌─────────────┐  ┌─────────────┐   │
                    │  │ Nginx Proxy │  │ Go Backend  │   │
                    │  │ (Port 80/   │──│ (Port 8080) │   │
                    │  │  443 SSL)   │  │ + Scheduler │   │
                    │  │             │  │ + GraphQL   │   │
                    │  └─────────────┘  └─────────────┘   │
                    │         │                │          │
                    │  ┌──────▼──────┐  ┌──────▼──────┐   │
                    │  │ PostgreSQL  │  │ Ollama AI   │   │
                    │  │ (Port 5432) │  │ (Port 11434)│   │
                    │  │ Internal    │  │ Internal    │   │
                    │  └─────────────┘  └─────────────┘   │
                    └──────────────────────────────────────┘
                              │
                              ▼ SMTP
                    ┌──────────────────┐
                    │   Email Provider │
                    │ (Gmail/Outlook)  │
                    └──────────────────┘
```

**Nginx Configuration:**

```nginx
server {
    listen 80;
    server_name dossier.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name dossier.example.com;

    ssl_certificate /etc/letsencrypt/live/dossier.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/dossier.example.com/privkey.pem;

    # Backend API
    location /graphql {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Frontend static files
    location / {
        root /var/www/dossier/dist;
        try_files $uri /index.html;
    }
}
```

### Cloud Deployment Options

**Option 1: VPS Deployment (DigitalOcean, Linode, Vultr)**

- **Requirements**: 4GB RAM minimum (for Ollama), 2 CPU cores, 50GB storage
- **Setup**: Docker Compose for all services
- **SSL**: Let's Encrypt via Certbot
- **Backups**: Automated PostgreSQL dumps to S3/Backblaze
- **Cost**: ~$20-40/month

**Option 2: Railway / Render**

- **Database**: Managed PostgreSQL service
- **Backend**: Docker container deployment
- **Frontend**: Static site deployment
- **Ollama**: Self-hosted on separate VPS (Ollama needs GPU/high RAM)
- **Cost**: ~$15-30/month + VPS for Ollama

**Option 3: Self-Hosted (Home Server / Raspberry Pi)**

- **Hardware**: Raspberry Pi 4 (8GB) or mini PC
- **Network**: Port forwarding or Cloudflare Tunnel
- **Power**: UPS recommended for 24/7 operation
- **Cost**: One-time hardware cost only

## Performance & Scaling

### Current Optimizations

- **Concurrent RSS Fetching**: Goroutines for parallel feed processing (multiple feeds fetched simultaneously)
- **Database Connection Pooling**: Efficient PostgreSQL connections via `database/sql` pool
- **Local AI Processing**: No external API latency (Ollama runs on localhost)
- **Smart Scheduling**: Duplicate prevention via `dossier_deliveries` timestamp queries
- **Efficient Queries**: Proper indexes on articles.published_at, dossier_configs.active, feeds.url
- **Async Delivery**: Each dossier processed in separate goroutine (non-blocking)

### Performance Characteristics

- **RSS Fetching**: ~1-3 seconds per feed (parallel)
- **AI Summary**: ~10-30 seconds (depends on article count and model)
- **Email Delivery**: ~1-2 seconds (SMTP send)
- **Total Per Dossier**: ~15-60 seconds (varies by article count)
- **Scheduler Overhead**: Negligible (<1% CPU) with 1-minute ticker

### Scaling Strategies

**Vertical Scaling (Single Server):**

- **More CPU**: Faster Ollama inference (2-4 cores recommended)
- **More RAM**: Required for Ollama models (4GB minimum, 8GB+ ideal)
- **SSD Storage**: Faster database queries and model loading

**Optimization Techniques:**

- **Smaller AI Models**: Use `llama3.2:1b` instead of `llama3.2:3b` for faster summaries
- **Reduce Article Count**: Limit to 10-15 articles instead of 20+
- **Database Optimization**: Add more indexes, tune PostgreSQL settings
- **Caching**: Redis for feed metadata caching (future enhancement)
- **CDN**: Serve static frontend assets via CDN (Cloudflare, etc.)

**Horizontal Scaling (Advanced):**

- **Separate AI Server**: Dedicate one server to Ollama, multiple backends
- **Read Replicas**: PostgreSQL read replicas for article queries
- **Load Balancer**: Distribute GraphQL requests across multiple backends
- **Message Queue**: Use RabbitMQ/Redis for async dossier generation

## Technology Decisions

### Why Go for Backend?

- **Concurrency**: Goroutines perfect for RSS fetching (parallel) and scheduling (async)
- **Performance**: Fast execution for real-time email generation (compiled binary)
- **Deployment**: Single binary deployment - easy to ship and run
- **Libraries**: Excellent GraphQL (`graphql-go/graphql`), HTTP (`net/http`), and database (`lib/pq`) libraries
- **Memory**: Efficient memory usage for always-running scheduler service
- **Simplicity**: Standard library covers most needs, minimal dependencies
- **Cross-Platform**: Compile once, run anywhere (Linux, macOS, Windows)

### Why Local AI (Ollama)?

- **Privacy**: No external API calls means complete data sovereignty
- **Cost**: No per-request charges like OpenAI ($0.01-0.03 per 1K tokens)
- **Unlimited**: Generate unlimited summaries without billing concerns
- **Reliability**: No network dependencies for AI processing (works offline)
- **Customization**: Multiple models (llama3.2:3b, dolphin-mistral), adjustable parameters
- **Speed**: Local processing after initial model download (~2GB one-time)
- **Control**: Fine-tune prompts, temperature, max tokens without API restrictions

### Why GraphQL?

- **Type Safety**: Compile-time schema validation catches errors early
- **Flexibility**: Frontend requests exactly the data needed (no over/under-fetching)
- **Single Endpoint**: Simplified API surface (`/graphql` for everything)
- **Introspection**: Self-documenting API with built-in schema exploration
- **Development**: Excellent tooling (GraphiQL playground, schema generators)
- **Evolution**: Easy schema changes without REST versioning (v1, v2, etc.)
- **Nested Queries**: Fetch dossier configs with nested deliveries in one request

### Why Vue.js 3?

- **Reactivity**: Proxy-based reactivity perfect for real-time dossier configuration
- **Composition API**: Clean component logic organization with `<script setup>`
- **Performance**: Efficient rendering with Virtual DOM and compiler optimizations
- **Developer Experience**: Excellent Vite integration (HMR, fast builds) and Vue DevTools
- **Ecosystem**: Rich component libraries (Vue Router, Vuex/Pinia)
- **Learning Curve**: Easier than React for beginners, simpler than Angular
- **Size**: Small bundle size (~30KB core) for fast page loads

### Why PostgreSQL?

- **Reliability**: ACID compliance for delivery tracking (critical for "already sent" logic)
- **Time Support**: Native TIME, TIMESTAMP, TIMEZONE types for scheduling
- **Arrays**: Native TEXT[] array support for feed_urls (no JSON parsing needed)
- **JSON**: JSONB support for flexible configuration storage (future use)
- **Performance**: Excellent indexing for article queries (published_at, feed_id)
- **Ecosystem**: Mature tooling (pg_dump for backups, pgAdmin for management)
- **Constraints**: CHECK constraints for data validation (frequency enum, article_count range)

### Why Docker?

- **Consistency**: Identical development and production environments
- **Isolation**: Each service in dedicated container (no port conflicts)
- **Dependencies**: AI models and database bundled cleanly (no manual setup)
- **Deployment**: Simple `docker-compose up` for entire stack
- **Scaling**: Easy horizontal scaling with orchestration (Kubernetes, Docker Swarm)
- **Development**: Hot-reload for backend/frontend without container rebuilds

## Monitoring Strategy

### Current Implementation

- **Structured Logging**: JSON logs to stdout with severity levels (debug/info/warn/error)
- **Scheduler Logging**: Detailed logs for every tick, dossier processing, and delivery attempt
- **AI Pipeline Logging**: Debug information for Ollama requests/responses, token counts
- **RSS Fetching Logging**: Feed fetch success/failure, article counts, parsing errors
- **Email Tracking**: Database records in `dossier_deliveries` for all delivery attempts
- **Error Handling**: Graceful degradation with detailed error messages (never crashes)
- **GraphQL Errors**: Structured error responses in GraphQL format

**Example Logs:**

```
[INFO] Scheduler started, checking every 1 minute
[INFO] Checking 3 active dossier configs for scheduled delivery
[INFO] Dossier "Tech News Daily" is due for delivery
[INFO] Fetching articles from 3 feed URLs
[INFO] Fetched 45 articles, limited to 15 for AI processing
[INFO] Generating summary with Ollama (model: llama3.2:3b, tone: professional)
[INFO] AI summary generated successfully (567 tokens)
[INFO] Sending email to user@example.com
[INFO] Email sent successfully, recorded delivery ID: 123
```

### Production Monitoring (Recommended)

**Health Checks:**

- `/health` endpoint for uptime monitoring (Uptime Robot, Pingdom)
- Check scheduler status via `schedulerStatus` GraphQL query
- Database connection health checks

**Metrics Collection:**

- Prometheus for system metrics (CPU, RAM, disk, network)
- Custom metrics: dossiers sent per day, AI generation times, RSS fetch failures
- Grafana dashboards for visualization

**Log Aggregation:**

- ELK Stack (Elasticsearch, Logstash, Kibana) for log analysis
- Or: Loki + Grafana for lightweight log aggregation
- Or: Cloud logging (CloudWatch, Datadog, Papertrail)

**Alerting:**

- Email delivery failures (SMTP errors)
- AI model unavailable (Ollama down)
- Database connection failures
- RSS feed fetch failures (all feeds)
- Scheduler stopped unexpectedly

**Database Performance:**

- Query times and slow query logs
- Connection pool usage and exhaustion
- Table sizes and growth rates
- Index usage statistics

**Email Delivery Monitoring:**

- Track bounce rates and delivery success
- Monitor SMTP authentication failures
- Alert on consecutive failures (3+ in a row)
