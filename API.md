# API Documentation

This document provides detailed information about the Dossier GraphQL API.

## Base URL

- Development: `http://localhost:8080/graphql`
- GraphQL Playground: `http://localhost:8080/graphql/playground`

## Authentication

All authenticated requests must include a JWT token in the Authorization header:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

You receive a token after successful login or registration.

## GraphQL Schema

### Types

#### User
```graphql
type User {
  id: ID!
  email: String!
  name: String!
  feeds: [Feed!]!
  digests: [Digest!]!
  createdAt: String!
}
```

#### Feed
```graphql
type Feed {
  id: ID!
  url: String!
  title: String
  description: String
  active: Boolean!
  articles: [Article!]!
  createdAt: String!
  updatedAt: String!
}
```

#### Article
```graphql
type Article {
  id: ID!
  feedId: ID!
  title: String!
  link: String!
  description: String
  content: String
  author: String
  publishedAt: String!
  createdAt: String!
}
```

#### Digest
```graphql
type Digest {
  id: ID!
  userId: ID!
  date: String!
  summary: String!
  articles: [Article!]!
  createdAt: String!
}
```

#### AuthPayload
```graphql
type AuthPayload {
  token: String!
  user: User!
}
```

## Queries

### Get Current User

```graphql
query {
  me {
    id
    email
    name
    createdAt
  }
}
```

**Authentication:** Required

**Returns:** The currently authenticated user

### Get All Feeds

```graphql
query {
  feeds {
    id
    url
    title
    description
    active
    createdAt
    updatedAt
  }
}
```

**Authentication:** Required

**Returns:** List of all feeds for the current user

### Get Single Feed

```graphql
query {
  feed(id: "1") {
    id
    url
    title
    description
    active
  }
}
```

**Authentication:** Required

**Parameters:**
- `id` (ID!): Feed ID

### Get Articles

```graphql
query {
  articles(limit: 20, offset: 0) {
    id
    feedId
    title
    link
    description
    content
    author
    publishedAt
    createdAt
  }
}
```

**Authentication:** Required

**Parameters:**
- `limit` (Int): Maximum number of articles to return (default: 50)
- `offset` (Int): Number of articles to skip (default: 0)

### Get Single Article

```graphql
query {
  article(id: "1") {
    id
    title
    link
    description
    content
    author
    publishedAt
  }
}
```

**Authentication:** Required

**Parameters:**
- `id` (ID!): Article ID

### Get Digests

```graphql
query {
  digests(limit: 10) {
    id
    userId
    date
    summary
    createdAt
    articles {
      id
      title
      link
    }
  }
}
```

**Authentication:** Required

**Parameters:**
- `limit` (Int): Maximum number of digests to return (default: 10)

### Get Single Digest

```graphql
query {
  digest(id: "1") {
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

**Authentication:** Required

**Parameters:**
- `id` (ID!): Digest ID

### Get Latest Digest

```graphql
query {
  latestDigest {
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

**Authentication:** Required

**Returns:** The most recent digest for the current user

## Mutations

### Register

```graphql
mutation {
  register(
    email: "user@example.com"
    password: "securepassword"
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

**Authentication:** Not required

**Parameters:**
- `email` (String!): User email address
- `password` (String!): User password
- `name` (String!): User full name

**Returns:** Auth token and user object

### Login

```graphql
mutation {
  login(
    email: "user@example.com"
    password: "securepassword"
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

**Authentication:** Not required

**Parameters:**
- `email` (String!): User email
- `password` (String!): User password

**Returns:** Auth token and user object

### Add Feed

```graphql
mutation {
  addFeed(url: "https://news.ycombinator.com/rss") {
    id
    url
    title
    description
    active
  }
}
```

**Authentication:** Required

**Parameters:**
- `url` (String!): RSS feed URL

**Returns:** The created feed

**Note:** The feed is validated before being added. Articles are fetched asynchronously.

### Update Feed

```graphql
mutation {
  updateFeed(
    id: "1"
    title: "My Custom Title"
    description: "My custom description"
    active: true
  ) {
    id
    title
    description
    active
  }
}
```

**Authentication:** Required

**Parameters:**
- `id` (ID!): Feed ID
- `title` (String): New title (optional)
- `description` (String): New description (optional)
- `active` (Boolean): Active status (optional)

**Returns:** The updated feed

### Delete Feed

```graphql
mutation {
  deleteFeed(id: "1")
}
```

**Authentication:** Required

**Parameters:**
- `id` (ID!): Feed ID to delete

**Returns:** Boolean indicating success

### Refresh Single Feed

```graphql
mutation {
  refreshFeed(id: "1") {
    id
    title
    updatedAt
  }
}
```

**Authentication:** Required

**Parameters:**
- `id` (ID!): Feed ID to refresh

**Returns:** The refreshed feed

### Refresh All Feeds

```graphql
mutation {
  refreshAllFeeds
}
```

**Authentication:** Required

**Returns:** Boolean indicating the operation was started

**Note:** This operation runs asynchronously in the background.

### Generate Digest

```graphql
mutation {
  generateDigest {
    id
    date
    summary
    articles {
      id
      title
      link
    }
  }
}
```

**Authentication:** Required

**Returns:** The generated digest

**Note:** This operation may take a few seconds to complete, especially if using the OpenAI API.

## Error Handling

The API returns errors in the standard GraphQL error format:

```json
{
  "errors": [
    {
      "message": "unauthorized",
      "path": ["me"]
    }
  ],
  "data": null
}
```

Common error messages:
- `unauthorized`: Authentication required or invalid token
- `invalid credentials`: Login failed
- `user already exists`: Email already registered
- `invalid feed URL`: Feed URL is not valid or inaccessible

## Rate Limiting

Currently, there is no rate limiting implemented. This may be added in future versions.

## Pagination

For queries that return lists (articles, digests), pagination is supported using `limit` and `offset` parameters:

```graphql
query {
  articles(limit: 10, offset: 20) {
    id
    title
  }
}
```

## Best Practices

1. **Always use HTTPS in production**
2. **Store tokens securely** (e.g., httpOnly cookies or secure storage)
3. **Handle errors gracefully** in your client application
4. **Use pagination** for large result sets
5. **Request only needed fields** to reduce response size
6. **Refresh tokens** before they expire (24-hour lifetime)

## Examples

### Complete User Flow

```graphql
# 1. Register
mutation {
  register(
    email: "new@example.com"
    password: "password123"
    name: "New User"
  ) {
    token
    user { id email }
  }
}

# 2. Add a feed (use token from step 1)
mutation {
  addFeed(url: "https://news.ycombinator.com/rss") {
    id
    title
  }
}

# 3. Refresh feeds
mutation {
  refreshAllFeeds
}

# 4. Get articles
query {
  articles(limit: 10) {
    title
    link
  }
}

# 5. Generate digest
mutation {
  generateDigest {
    summary
  }
}
```

## Support

For issues or questions about the API:
- Open an issue on GitHub
- Check the main [README.md](README.md) for more information
- See [QUICKSTART.md](QUICKSTART.md) for getting started
