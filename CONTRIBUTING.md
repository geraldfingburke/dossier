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

### Backend Development

```bash
# Install Go dependencies
go mod download

# Run the server with hot reload
go run server/cmd/main.go

# Or use the Makefile
make run-server
```

### Frontend Development

```bash
cd client
npm install
npm run dev
```

### Database

You can use Docker for PostgreSQL:

```bash
docker run -d \
  --name dossier-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=dossier \
  -p 5432:5432 \
  postgres:15-alpine
```

## Project Structure

```
dossier/
├── server/           # Go backend
│   ├── cmd/          # Application entry points
│   └── internal/     # Internal packages
│       ├── ai/       # AI service
│       ├── auth/     # Authentication
│       ├── database/ # Database layer
│       ├── graphql/  # GraphQL API
│       ├── models/   # Data models
│       └── rss/      # RSS service
└── client/           # Vue.js frontend
    └── src/
        ├── components/ # Reusable components
        ├── views/      # Page views
        └── store/      # State management
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
