# Contributing to Dossier

Thank you for your interest in contributing to Dossier! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/dossier.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit with a clear message: `git commit -m "Add feature: description"`
7. Push to your fork: `git push origin feature/your-feature-name`
8. Create a Pull Request

## Development Setup

### Quick Start (Recommended)

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/dossier.git
cd dossier

# Start all services with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f
```

This starts:

- PostgreSQL on port 5432
- Ollama AI service on port 11434
- Go backend on port 8080
- Vue.js frontend on port 5173

### Manual Development Setup

#### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker (for Ollama)

#### Backend Development

```bash
# Install Go dependencies
go mod download

# Set environment variables
cp .env.example .env
# Edit .env with your configuration

# Run the server with scheduler
go run server/cmd/main.go
```

#### Frontend Development

```bash
cd client
npm install
npm run dev
```

#### AI Service Setup

```bash
# Start Ollama in Docker
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama

# Pull required models
docker exec -it ollama ollama pull llama3.2:3b
docker exec -it ollama ollama pull dolphin-mistral
```

## Project Structure

```
dossier/
├── server/           # Go backend with scheduler
│   ├── cmd/          # Application entry points
│   │   └── main.go   # Server with integrated scheduler
│   └── internal/     # Internal packages
│       ├── ai/       # Local AI service (Ollama integration)
│       ├── database/ # Database layer and migrations
│       ├── graphql/  # GraphQL API and resolvers
│       ├── models/   # Data models (dossiers, feeds, articles)
│       └── rss/      # RSS feed processing service
├── client/           # Vue.js 3 frontend
│   └── src/
│       ├── views/    # Main application pages
│       │   ├── DossiersView.vue  # Dossier management
│       │   ├── FeedsView.vue     # Feed management
│       │   └── ArticlesView.vue  # Article browsing
│       ├── store/    # Vuex state management
│       └── App.vue   # Root component
├── docker-compose.yml    # Complete development environment
├── Dockerfile           # Production container build
├── QUICKSTART.md        # Quick setup guide
└── SMTP_SETUP.md        # Email configuration guide
```

## Code Style

### Go

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions small and focused
- Handle errors explicitly

### JavaScript/Vue

- Use ES6+ features
- Follow Vue.js style guide
- Use Composition API for new components
- Keep components small and reusable
- Use meaningful prop and event names

## Testing

### Backend Tests

```bash
go test ./...
```

### Frontend Tests

Currently, the project doesn't have tests. Contributions to add testing infrastructure are welcome!

## Adding New Features

When adding new features:

1. Check if there's an existing issue or create one
2. Discuss the feature before implementing
3. Follow the existing code structure
4. Update documentation
5. Add tests if possible
6. Update the README if needed

## GraphQL Schema Changes

When modifying the GraphQL schema:

1. Update `server/internal/graphql/schema.graphql`
2. Update the resolver implementation in `server/internal/graphql/graphql.go`
3. Update frontend queries/mutations if needed
4. Document the changes

## Database Migrations

For database schema changes:

1. Update the schema in `server/internal/database/database.go`
2. Consider backward compatibility
3. Test migrations on a clean database
4. Document any manual migration steps needed

## Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Write clear commit messages
- Update documentation
- Ensure the code builds successfully
- Test your changes thoroughly
- Link to related issues

## Code Review Process

All submissions require review. We'll:

- Check code quality and style
- Verify tests pass
- Ensure documentation is updated
- Test the feature/fix

## Reporting Bugs

When reporting bugs, include:

- Description of the issue
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, Node version)
- Relevant logs or error messages

## Feature Requests

For feature requests:

- Describe the feature clearly
- Explain the use case
- Consider implementation complexity
- Discuss alternatives

## Questions?

Feel free to open an issue for questions or join discussions.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
