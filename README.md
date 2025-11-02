# Dossier

An automated RSS digest system that sends personalized email summaries on your schedule. Built with Go, GraphQL, and Vue.js, Dossier transforms your RSS feeds into intelligent email digests using local AI processing.

## Features

- 📧 **Automated Email Delivery**: Scheduled digests sent directly to your inbox
- 🤖 **Local AI Processing**: Free AI summaries using Ollama (no OpenAI required)
- 📰 **Multi-Feed Support**: Combine articles from multiple RSS feeds per dossier
- 🎭 **Customizable Tones**: 10 system defaults + custom user-defined tones
- 🌍 **Multi-language Support**: Generate summaries in any language
- ⏰ **Flexible Scheduling**: Daily, weekly, or monthly delivery with timezone support
- 🎯 **Custom Instructions**: Fine-tune AI behavior with special prompts
- 👤 **Single-User Design**: No authentication needed, perfect for self-hosting
- 📱 **Modern UI**: Clean, responsive Vue.js 3 interface with modular CSS
- 📊 **GraphQL API**: Flexible and type-safe data fetching
- 🐳 **Docker Support**: Easy deployment with Docker Compose
- 📜 **Delivery History**: Full archive of sent digests with content

## Tech Stack

### Backend

- **Go 1.21+**: Server-side language with goroutines for concurrent processing
- **GraphQL**: API layer (github.com/graphql-go/graphql)
- **PostgreSQL 15**: Database with timezone-aware scheduling
- **Ollama**: Local LLM inference (llama3.2:3b, dolphin-mistral)
- **gofeed**: RSS/Atom feed parsing (supports RSS 1.0, RSS 2.0, Atom 1.0)
- **SMTP/TLS**: Secure email delivery
- **Scheduler**: Ticker-based (1-minute) automated delivery system

### Frontend

- **Vue.js 3**: Frontend framework with Composition API and `<script setup>`
- **Vite**: Build tool and dev server with hot module replacement
- **Vuex**: State management for dossier configs and tones
- **Modular CSS**: 9 separate stylesheets with design tokens

### AI & Processing

- **Ollama**: Local LLM inference (privacy-focused, no external APIs)
- **Models**: llama3.2:3b (default), dolphin-mistral (uncensored)
- **Single-Stage Pipeline**: Article formatting → AI summary generation
- **Tone System**: Customizable prompts for different writing styles

## Quick Start

### Prerequisites

- Docker and Docker Compose
- SMTP email credentials (Gmail, Outlook, etc.)

### 5-Minute Setup

1. **Clone and configure**

   ```bash
   git clone https://github.com/geraldfingburke/dossier.git
   cd dossier

   # Copy and edit environment file
   cp .env.example .env
   # Edit .env with your SMTP credentials
   ```

2. **Start the application**

   ```bash
   docker-compose up -d

   # Wait for Ollama to download model (first run only, ~2GB)
   docker-compose logs -f ollama
   ```

3. **Access and configure**
   - Open http://localhost:5173
   - Click "New Dossier Config"
   - Add RSS feed URLs (one per line)
   - Set schedule, timezone, and delivery email
   - Choose AI tone and language
   - Test with the "Send Test Email" button

That's it! Your automated dossiers will be delivered on schedule.

For detailed setup instructions, see [QUICKSTART.md](QUICKSTART.md)

## How It Works

### 1. Configuration

- Create dossier configurations with RSS feed URLs, delivery preferences, and AI settings
- Set flexible schedules: daily, weekly, or monthly delivery with timezone support
- Configure multiple dossiers for different topics (tech news, sports, finance, etc.)
- Customize AI behavior with tone selection and special instructions

### 2. Automated Processing

- **Scheduler**: Checks every minute for due dossiers (timezone-aware)
- **Article Fetching**: Concurrently fetches articles from all configured RSS feeds
- **Article Aggregation**: Sorts by published date, limits to configured article count
- **AI Summary**: Generates personalized summary using Ollama with chosen tone
- **Email Delivery**: Sends HTML-formatted email via SMTP with TLS encryption

### 3. Email Delivery

- Generates professional HTML email with article summaries
- Includes article titles, descriptions, and links to original sources
- Sends at your scheduled time in your specified timezone
- Tracks delivery history to prevent duplicates
- Archives full content for later viewing

## Usage

### Dossier Management

- **Create Config**: Click "New Dossier Config" to create a configuration
- **Configure RSS Feeds**: Add feed URLs (one per line) in the textarea
- **Set Schedule**: Choose frequency (daily/weekly/monthly), time, and timezone
- **Customize AI**: Select tone, language, and add special instructions
- **Test**: Use "Send Test Email" button to verify configuration
- **Toggle Active**: Enable/disable configs without deletion
- **View History**: Click "View Digests" to see past deliveries

### Available Tones (10 System Defaults)

- **professional**: Standard business communication style
- **humorous**: Witty and entertaining summaries
- **analytical**: Data-driven insights and trends
- **casual**: Relaxed, conversational tone
- **apocalyptic**: Dramatic, foreboding style with biblical references
- **orc**: Warcraft-style blunt communication ("Me Grognak!")
- **robot**: Mechanical, technical language ("EXECUTING SUMMARY PROTOCOL")
- **southern_belle**: Polite, charming Southern style ("Well, bless your heart")
- **apologetic**: Sympathetic and reassuring ("I'm so sorry to report...")
- **sweary**: Adult language (requires uncensored dolphin-mistral model)

**Custom Tones**: Create your own tones with custom prompts via the UI

### Multi-language Support

Generate summaries in any language by setting the language field: English, Spanish, French, German, Japanese, etc.

### Scheduler Behavior

- **Granularity**: Checks every 1 minute for due dossiers
- **Daily**: Delivers at specified time each day
- **Weekly**: Delivers same day of week, 7+ days after last delivery
- **Monthly**: Delivers same day of month, 30+ days after last delivery
- **Duplicate Prevention**: Tracks last delivery to avoid re-sending

## Development

### Project Architecture

The system follows a clean, single-user architecture:

- **Frontend**: Vue.js 3 SPA with modular CSS (9 files) and Vuex state management
- **Backend**: Go GraphQL API with automated scheduler service (1-minute ticker)
- **Database**: PostgreSQL with 6 tables (configs, feeds, articles, deliveries, delivery_articles, tones)
- **AI Processing**: Local Ollama integration (no external APIs)
- **Email Service**: SMTP/TLS with HTML template rendering
- **Containerization**: Docker Compose for complete development environment

### Local Development

1. **Prerequisites**

   - Docker and Docker Compose
   - Git

2. **Clone and Start**

   ```bash
   git clone https://github.com/geraldfingburke/dossier.git
   cd dossier

   # Copy environment file
   cp .env.example .env
   # Edit .env with SMTP credentials

   # Start all services
   docker-compose up -d
   ```

3. **Access Services**
   - Frontend: http://localhost:5173 (Vite dev server with HMR)
   - GraphQL API: http://localhost:8080/graphql
   - PostgreSQL: localhost:5432 (internal to Docker network)
   - Ollama: http://localhost:11434 (internal)

## API Examples

See [API.md](API.md) for complete GraphQL schema documentation.

**Create a dossier config:**

```graphql
mutation {
  createDossierConfig(
    input: {
      title: "Tech News Daily"
      email: "you@example.com"
      feedUrls: [
        "https://news.ycombinator.com/rss"
        "https://techcrunch.com/feed/"
      ]
      articleCount: 15
      frequency: "daily"
      deliveryTime: "08:00"
      timezone: "America/New_York"
      tone: "professional"
      language: "english"
    }
  ) {
    id
    title
    active
  }
}
```

**Get config with delivery history:**

```graphql
query {
  dossierConfig(id: "1") {
    title
    email
    feedUrls
    frequency
    deliveryTime
    timezone
    tone
  }

  dossiers(configId: "1", limit: 10) {
    id
    subject
    sentAt
    articleCount
  }
}
```

**Create custom tone:**

```graphql
mutation {
  createTone(
    input: {
      name: "pirate"
      prompt: "Write like a pirate with 'arr' and 'matey'. Use nautical metaphors."
    }
  ) {
    id
    name
  }
}
```

## Project Structure

```
dossier/
├── server/
│   ├── cmd/
│   │   └── main.go                 # Entry point with scheduler initialization
│   └── internal/
│       ├── ai/
│       │   └── ai.go               # Ollama LLM integration
│       ├── database/
│       │   └── database.go         # PostgreSQL connection & schema migrations
│       ├── email/
│       │   └── email.go            # SMTP/TLS email delivery
│       ├── graphql/
│       │   ├── graphql.go          # Resolvers & schema implementation
│       │   └── schema.graphql      # GraphQL type definitions
│       ├── models/
│       │   └── models.go           # Domain models (DossierConfig, Article, Tone, etc.)
│       ├── rss/
│       │   └── rss.go              # RSS/Atom feed fetching & parsing
│       └── scheduler/
│           └── scheduler.go        # Automated delivery scheduler (1-minute ticker)
├── client/
│   ├── src/
│   │   ├── views/
│   │   │   ├── DossierConfigsView.vue   # Main config management
│   │   │   ├── DigestsView.vue          # Delivery history
│   │   │   └── ArticlesView.vue         # Article browser
│   │   ├── store/
│   │   │   └── index.js                 # Vuex state management
│   │   ├── styles/
│   │   │   ├── variables.css            # Design tokens
│   │   │   ├── reset.css                # CSS reset
│   │   │   ├── buttons.css              # Button styles
│   │   │   ├── forms.css                # Form styles
│   │   │   ├── components.css           # Component styles
│   │   │   ├── layout.css               # Layout utilities
│   │   │   ├── modals.css               # Modal dialogs
│   │   │   ├── utilities.css            # Utility classes
│   │   │   └── main.css                 # Main imports
│   │   ├── App.vue                      # Root component
│   │   └── main.js                      # Entry point
│   ├── index.html                       # HTML template
│   └── vite.config.js                   # Vite configuration
├── docker-compose.yml                   # Development environment (4 services)
├── Dockerfile                           # Production build
├── go.mod                               # Go dependencies
├── API.md                               # GraphQL API documentation
├── ARCHITECTURE.md                      # System architecture details
├── QUICKSTART.md                        # Quick setup guide
└── README.md                            # This file
```

## Configuration

### Environment Variables

**Database:**

- `DATABASE_URL`: PostgreSQL connection string (default: postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable)

**Server:**

- `PORT`: Server port (default: 8080)

**AI Service:**

- `OLLAMA_URL`: Ollama server URL (default: http://localhost:11434)
- `AI_MODEL`: Model name (default: llama3.2:3b)
- `AI_UNCENSORED_MODEL`: Uncensored model for mature tones (default: dolphin-mistral)

**Email Service (Required for delivery):**

- `SMTP_HOST`: SMTP server hostname (e.g., smtp.gmail.com)
- `SMTP_PORT`: SMTP server port (e.g., 587)
- `SMTP_USER`: SMTP username (your email address)
- `SMTP_PASS`: SMTP password (app-specific password for Gmail)
- `SMTP_FROM`: From address for outgoing emails

See [QUICKSTART.md](QUICKSTART.md) for detailed email configuration instructions.

### Building for Production

**Backend:**

```bash
go build -o bin/server ./server/cmd/main.go
./bin/server
```

**Frontend:**

```bash
cd client
npm run build
```

**Docker Production Build:**

```bash
docker build -t dossier .
docker run -p 8080:8080 --env-file .env dossier
```

## Security Notes

- **Single-User Design**: No authentication system, intended for personal use on private servers
- **Protect SMTP Credentials**: Use app-specific passwords, never your main account password
- **Use Environment Variables**: Never commit sensitive data (.env) to version control
- **Enable HTTPS in Production**: Secure all client-server communication with SSL/TLS
- **Local AI Processing**: All AI processing happens locally via Ollama - no external API calls or data leakage
- **Email Security**: SMTP connections use TLS encryption for secure email delivery
- **Docker Isolation**: Services run in isolated containers with internal networking
- **No JWT Tokens**: Single-user design eliminates authentication attack surface

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT

## Acknowledgments

- **GraphQL API**: Built with [graphql-go/graphql](https://github.com/graphql-go/graphql) for type-safe GraphQL server
- **RSS Processing**: Powered by [gofeed](https://github.com/mmcdole/gofeed) for reliable RSS/Atom feed parsing
- **Local AI**: [Ollama](https://ollama.ai/) for privacy-focused local language model inference
- **Frontend**: [Vue.js 3](https://vuejs.org/) with [Vite](https://vitejs.dev/) for fast development and hot module replacement
- **Database**: [PostgreSQL](https://postgresql.org/) for reliable data persistence with timezone support
- **Email**: Go's `net/smtp` package for SMTP/TLS email delivery
- **Containerization**: [Docker](https://docker.com/) for consistent development and deployment environments
