# Dossier

An AI-assisted RSS digest application built with Go, GraphQL, and Vue.js. Dossier helps you stay on top of your favorite RSS feeds by automatically fetching articles and generating AI-powered daily summaries.

## Features

- 📰 **RSS Feed Management**: Add and manage multiple RSS feeds
- 🤖 **AI-Powered Summaries**: Automatic daily digests using OpenAI
- 📱 **Modern UI**: Clean, responsive Vue.js interface
- 🔒 **User Authentication**: Secure JWT-based authentication
- 📊 **GraphQL API**: Flexible and efficient data fetching
- 🐳 **Docker Support**: Easy deployment with Docker Compose

## Tech Stack

### Backend
- **Go**: Server-side language
- **GraphQL**: API layer (graphql-go)
- **PostgreSQL**: Database
- **Chi**: HTTP router
- **OpenAI API**: Content summarization

### Frontend
- **Vue.js 3**: Frontend framework with Composition API
- **Vite**: Build tool and dev server
- **GraphQL Request**: GraphQL client

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Node.js 20 or higher
- PostgreSQL 15 or higher
- OpenAI API key (optional, works with mock summaries without it)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/geraldfingburke/dossier.git
   cd dossier
   ```

2. **Set up environment variables**
   ```bash
   # Optional: Create a .env file
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable"
   export OPENAI_API_KEY="your-openai-api-key"  # Optional
   export PORT="8080"
   ```

3. **Start PostgreSQL**
   ```bash
   # Using Docker
   docker run -d \
     --name dossier-postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=dossier \
     -p 5432:5432 \
     postgres:15-alpine
   ```

4. **Run the backend**
   ```bash
   # Install dependencies
   go mod download
   
   # Build and run
   go run server/cmd/main.go
   ```
   
   The server will start on http://localhost:8080
   - GraphQL endpoint: http://localhost:8080/graphql
   - GraphQL Playground: http://localhost:8080/graphql/playground

5. **Run the frontend**
   ```bash
   cd client
   npm install
   npm run dev
   ```
   
   The frontend will start on http://localhost:5173

### Using Docker Compose

The easiest way to run the entire stack:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

This will start:
- PostgreSQL on port 5432
- Backend server on port 8080
- Frontend dev server on port 5173

## Usage

1. **Register/Login**: Create an account or login to access the application

2. **Add RSS Feeds**: Navigate to the "Feeds" tab and add your favorite RSS feeds
   - Examples:
     - Hacker News: https://news.ycombinator.com/rss
     - TechCrunch: https://techcrunch.com/feed/
     - The Verge: https://www.theverge.com/rss/index.xml

3. **View Articles**: Click "Articles" to see all fetched articles from your feeds

4. **Generate Digests**: Go to "Digests" and click "Generate New Digest" to create an AI summary of recent articles

5. **Automatic Daily Digests**: The server automatically generates digests every 24 hours for all users

## API Examples

### GraphQL Mutations

**Register a new user:**
```graphql
mutation {
  register(
    email: "user@example.com"
    password: "password123"
    name: "John Doe"
  ) {
    token
    user {
      id
      email
      name
    }
  }
}
```

**Add an RSS feed:**
```graphql
mutation {
  addFeed(url: "https://news.ycombinator.com/rss") {
    id
    title
    description
  }
}
```

**Generate a digest:**
```graphql
mutation {
  generateDigest {
    id
    date
    summary
    articles {
      title
      link
    }
  }
}
```

### GraphQL Queries

**Get your feeds:**
```graphql
query {
  feeds {
    id
    url
    title
    active
  }
}
```

**Get recent articles:**
```graphql
query {
  articles(limit: 20) {
    id
    title
    link
    description
    publishedAt
  }
}
```

**Get digests:**
```graphql
query {
  digests(limit: 10) {
    id
    date
    summary
    articles {
      title
      link
    }
  }
}
```

## Project Structure

```
dossier/
├── server/
│   ├── cmd/
│   │   └── main.go              # Server entry point
│   └── internal/
│       ├── ai/                   # AI service (OpenAI integration)
│       ├── auth/                 # Authentication logic
│       ├── database/             # Database connection & migrations
│       ├── graphql/              # GraphQL schema & resolvers
│       ├── models/               # Data models
│       └── rss/                  # RSS feed fetching & parsing
├── client/
│   ├── src/
│   │   ├── components/           # Vue components
│   │   ├── views/                # Page views
│   │   ├── store/                # State management
│   │   ├── App.vue               # Root component
│   │   └── main.js               # Entry point
│   ├── index.html                # HTML template
│   └── vite.config.js            # Vite configuration
├── docker-compose.yml            # Docker Compose configuration
├── Dockerfile                    # Multi-stage Docker build
├── go.mod                        # Go dependencies
└── README.md                     # This file
```

## Configuration

### Environment Variables

- `DATABASE_URL`: PostgreSQL connection string (default: postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable)
- `PORT`: Server port (default: 8080)
- `OPENAI_API_KEY`: OpenAI API key for AI summaries (optional, uses mock summaries without it)
- `JWT_SECRET`: Secret for JWT token signing (default: development-secret-key-change-in-production)

## Development

### Building the Backend

```bash
go build -o bin/server ./server/cmd/main.go
./bin/server
```

### Building the Frontend

```bash
cd client
npm run build
```

The built files will be in `client/dist/`

## Security Notes

- Change the JWT secret in production
- Use environment variables for sensitive data
- Never commit API keys or passwords
- Use HTTPS in production
- The OpenAI API key is optional - the app works with mock summaries without it

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT

## Acknowledgments

- Built with [gqlgen](https://gqlgen.com/) for GraphQL
- RSS parsing by [gofeed](https://github.com/mmcdole/gofeed)
- AI summaries powered by [OpenAI](https://openai.com/)
- Frontend built with [Vue.js](https://vuejs.org/) and [Vite](https://vitejs.dev/)

