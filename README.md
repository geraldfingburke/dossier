# Dossier

An automated RSS digest system that sends personalized email summaries on your schedule. Built with Go, GraphQL, and Vue.js, Dossier transforms your RSS feeds into intelligent email digests using local AI processing.

## Features

- 📧 **Automated Email Delivery**: Scheduled dossiers sent directly to your inbox
- 🤖 **Local AI Processing**: Free AI summaries using Ollama (no OpenAI required)
- 📰 **Multi-Feed Support**: Combine articles from multiple RSS feeds
- 🎭 **Multiple Tones**: Professional, humorous, analytical, apocalyptic, and more
- 🌍 **Multi-language Support**: Generate dossiers in any language
- ⏰ **Flexible Scheduling**: Daily, weekly, or monthly delivery
- 🎯 **Custom Instructions**: Fine-tune AI behavior with special instructions
- � **Single-User Design**: No accounts needed, perfect for self-hosting
- 📱 **Modern UI**: Clean, responsive Vue.js interface
- 📊 **GraphQL API**: Flexible and efficient data fetching
- 🐳 **Docker Support**: Easy deployment with Docker Compose

## Tech Stack

### Backend

- **Go**: Server-side language
- **GraphQL**: API layer (graphql-go)
- **PostgreSQL**: Database with scheduled delivery tracking
- **Chi**: HTTP router
- **Ollama**: Local AI processing (free alternative to OpenAI)
- **SMTP**: Email delivery system
- **Scheduler**: Cron-like automated dossier generation

### Frontend

- **Vue.js 3**: Frontend framework with Composition API
- **Vite**: Build tool and dev server
- **GraphQL Request**: GraphQL client
- **VeeValidate**: Form validation

### AI & Processing

- **Ollama**: Local LLM inference (llama3.2:3b, dolphin-mistral)
- **3-Stage Pipeline**: Article selection, content extraction, summary generation
- **Multiple Models**: Supports various local AI models

## Quick Start

### Prerequisites

- Docker and Docker Compose
- SMTP email credentials (Gmail, Outlook, etc.)
- (Optional) OpenAI API key (uses free local LLM by default)

### 5-Minute Setup

1. **Clone and configure**

   ```bash
   git clone https://github.com/geraldfingburke/dossier.git
   cd dossier

   # Configure SMTP (interactive setup)
   ./setup-smtp.sh        # Linux/macOS
   .\setup-smtp.ps1       # Windows PowerShell
   ```

2. **Start the application**

   ```bash
   docker-compose up -d
   ```

3. **Access and configure**
   - Open http://localhost:5173
   - Click "Add New Dossier"
   - Configure your RSS feeds, schedule, and preferences
   - Test with the "Test Email" button

That's it! Your automated dossiers will be delivered on schedule.

For detailed setup instructions, see [QUICKSTART.md](QUICKSTART.md)

## How It Works

### 1. Configuration

- Create dossier configurations with RSS feeds, delivery preferences, and AI tone settings
- Set flexible schedules: daily, weekly, or monthly delivery
- Configure multiple dossiers for different topics (tech news, sports, etc.)

### 2. Automated Processing

- **Article Collection**: Fetches articles from configured RSS feeds
- **AI Selection**: Intelligently selects most relevant articles
- **Content Extraction**: Cleans and extracts factual information
- **Summary Generation**: Creates personalized summaries with your chosen tone

### 3. Email Delivery

- Generates HTML email with formatted summaries
- Includes links to original articles
- Sends at your scheduled time with timezone support
- Tracks delivery history to prevent duplicates

## Usage

### Dossier Management

- **Add New Dossier**: Click the + button to create a new configuration
- **Configure RSS Feeds**: Add multiple feeds per dossier
- **Set Delivery Schedule**: Choose frequency and time
- **Customize AI Behavior**: Select tone, language, and special instructions
- **Test Configuration**: Use "Test Email" button to verify setup

### Available Tones

- **Professional**: Standard business communication style
- **Humorous**: Witty and entertaining summaries
- **Analytical**: Data-driven insights and trends
- **Casual**: Relaxed, conversational tone
- **Apocalyptic/Doomsayer**: Dramatic, foreboding style with biblical references
- **Orc**: Warcraft-style blunt communication
- **Robot**: Mechanical, technical language
- **Southern Belle**: Polite, charming Southern style
- **Apologetic**: Sympathetic and reassuring
- **Sweary**: Adult language for mature audiences (requires uncensored model)

### Multi-language Support

Generate dossiers in any language: English, Spanish, French, German, Japanese, etc.

## Development

### Project Architecture

The system follows a clean architecture pattern:

- **Frontend**: Vue.js 3 SPA with dossier management UI
- **Backend**: Go GraphQL API with automated scheduler service
- **Database**: PostgreSQL with time-aware delivery tracking
- **AI Processing**: Local Ollama integration with multiple models
- **Email Service**: SMTP with HTML template rendering
- **Containerization**: Docker Compose for complete environment

### Local Development

1. **Prerequisites**

   - Docker and Docker Compose
   - Git

2. **Clone and Start**

   ```bash
   git clone https://github.com/yourusername/dossier.git
   cd dossier
   docker-compose up -d
   ```

3. **Access Services**
   - Frontend: http://localhost:5173
   - GraphQL API: http://localhost:8080/graphql
   - GraphQL Playground: http://localhost:8080/graphql/playground

## API Examples

See [API.md](API.md) for complete GraphQL schema documentation.

**Create a dossier:**

```graphql
mutation {
  createDossier(input: {
    name: "Tech News",
    deliveryTime: "08:00",
    frequency: DAILY,
    tone: "professional",
    emailTo: "you@example.com"
  }) {
    id
    name
  }
```

**Add RSS feed to dossier:**

```graphql
mutation {
  addFeedToDossier(dossierId: "123", url: "https://news.ycombinator.com/rss") {
    id
    title
  }
}
```

**Get dossier with delivery history:**

```graphql
query {
  dossier(id: "123") {
    name
    deliveryTime
    frequency
    feeds {
      url
      title
    }
    deliveries {
      deliveredAt
      status
    }
  }
}
```

## Project Structure

```
dossier/
├── server/
│   ├── cmd/
│   │   └── main.go              # Server entry point with scheduler
│   └── internal/
│       ├── ai/                   # Local AI service (Ollama integration)
│       │   └── ai.go            # 3-stage preprocessing pipeline
│       ├── database/             # Database connection & migrations
│       ├── graphql/              # GraphQL schema & resolvers
│       │   ├── graphql.go       # Resolver implementations
│       │   └── schema.graphql   # GraphQL schema definition
│       ├── models/               # Data models & database structs
│       └── rss/                  # RSS feed fetching & parsing
├── client/
│   ├── src/
│   │   ├── views/                # Main application pages
│   │   │   ├── DossiersView.vue  # Dossier management
│   │   │   ├── FeedsView.vue     # Feed management
│   │   │   └── ArticlesView.vue  # Article browsing
│   │   ├── store/                # Vuex state management
│   │   ├── App.vue               # Root component
│   │   └── main.js               # Entry point
│   ├── index.html                # HTML template
│   └── vite.config.js            # Vite configuration
├── docker-compose.yml            # Complete development environment
├── Dockerfile                    # Multi-stage production build
├── go.mod                        # Go dependencies
├── QUICKSTART.md                 # Quick setup guide
├── SMTP_SETUP.md                 # Email configuration guide
└── README.md                     # This file
```

## Configuration

### Environment Variables

**Database:**

- `DATABASE_URL`: PostgreSQL connection string (default: postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable)

**Server:**

- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: Secret for JWT token signing (default: development-secret-key-change-in-production)

**AI Service:**

- `OLLAMA_URL`: Ollama server URL (default: http://localhost:11434)
- `AI_MODEL`: Model name (default: llama3.2:3b)
- `AI_UNCENSORED_MODEL`: Uncensored model for mature tones (default: dolphin-mistral)

**Email Service (Required for delivery):**

- `SMTP_HOST`: SMTP server hostname (e.g., smtp.gmail.com)
- `SMTP_PORT`: SMTP server port (e.g., 587)
- `SMTP_USERNAME`: SMTP username (your email address)
- `SMTP_PASSWORD`: SMTP password (app-specific password for Gmail)
- `FROM_EMAIL`: From address for outgoing emails
- `FROM_NAME`: Display name for sender (default: "Dossier Service")

See [SMTP_SETUP.md](SMTP_SETUP.md) for detailed email configuration instructions.

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

- **Change JWT secret in production**: Use a secure random string for `JWT_SECRET`
- **Protect SMTP credentials**: Use app-specific passwords, never your main account password
- **Use environment variables**: Never commit sensitive data to version control
- **Enable HTTPS in production**: Secure all client-server communication
- **Local AI processing**: All AI processing happens locally via Ollama - no external API calls
- **Email security**: SMTP connections use TLS encryption for secure email delivery

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT

## Acknowledgments

- **GraphQL API**: Built with [gqlgen](https://gqlgen.com/) for type-safe GraphQL
- **RSS Processing**: Powered by [gofeed](https://github.com/mmcdole/gofeed) for reliable feed parsing
- **Local AI**: [Ollama](https://ollama.ai/) for privacy-focused local language models
- **Frontend**: [Vue.js 3](https://vuejs.org/) with [Vite](https://vitejs.dev/) for fast development
- **Email Templates**: HTML email rendering with Go templates
- **Containerization**: [Docker](https://docker.com/) for consistent deployment
- **Database**: [PostgreSQL](https://postgresql.org/) for reliable data persistence
- **Time Handling**: Timezone-aware scheduling with Go's time package
