# Dossier GraphQL API Documentation

This document provides comprehensive information about the Dossier system's GraphQL API.

## Base URLs

- **Development**: `http://localhost:8080/graphql`
- **GraphQL Playground**: `http://localhost:8080/graphql/playground`

## Authentication

The Dossier system uses a **single-user design** with no authentication required

## Core Concepts

- **Dossier**: A configuration for automated email delivery with RSS feeds, schedule, and AI settings
- **Feed**: RSS/Atom feed sources that provide articles
- **Article**: Individual news items collected from feeds
- **Delivery**: Historical record of sent dossier emails

## GraphQL Schema

### Core Types

#### Dossier

```graphql
type Dossier {
  id: ID!
  name: String! # e.g., "Tech News", "Sports Updates"
  deliveryTime: String! # e.g., "08:00" (24-hour format)
  frequency: Frequency! # DAILY, WEEKLY, MONTHLY
  timezone: String # e.g., "America/New_York"
  tone: String! # AI summary tone
  language: String # Summary language (default: "english")
  specialInstructions: String # Custom AI prompts
  emailTo: String! # Delivery email address
  isActive: Boolean!
  feeds: [Feed!]! # Associated RSS feeds
  deliveries: [DossierDelivery!]! # Delivery history
  createdAt: String!
  updatedAt: String!
}
```

#### Feed

```graphql
type Feed {
  id: ID!
  url: String! # RSS/Atom feed URL
  title: String # Feed title from metadata
  description: String # Feed description
  lastFetchedAt: String # Last successful fetch time
  articles: [Article!]! # Articles from this feed
  dossiers: [Dossier!]! # Dossiers using this feed
  createdAt: String!
  updatedAt: String!
}
```

#### Article

```graphql
type Article {
  id: ID!
  feedId: ID! # Reference to source feed
  title: String! # Article headline
  link: String! # Original article URL (unique)
  description: String # Article summary/excerpt
  content: String # Full article content
  author: String # Article author
  publishedAt: String! # Original publication date
  feed: Feed! # Source feed relationship
  createdAt: String!
}
```

#### DossierDelivery

```graphql
type DossierDelivery {
  id: ID!
  dossierId: ID! # Reference to dossier
  deliveredAt: String! # Delivery timestamp with timezone
  status: DeliveryStatus! # SENT, FAILED
  emailContent: String # Generated HTML email content
  articleCount: Int! # Number of articles in delivery
  dossier: Dossier! # Dossier relationship
}
```

#### Enums

```graphql
enum Frequency {
  DAILY
  WEEKLY
  MONTHLY
}

enum DeliveryStatus {
  SENT
  FAILED
}
```

### Input Types

#### DossierInput

```graphql
input DossierInput {
  name: String!
  deliveryTime: String! # Format: "HH:MM" (24-hour)
  frequency: Frequency!
  timezone: String # IANA timezone (optional)
  tone: String! # AI tone selection
  language: String # Summary language (optional)
  specialInstructions: String # Custom AI instructions (optional)
  emailTo: String! # Delivery email address
  feedUrls: [String!]! # Initial RSS feed URLs
}
```

## Queries

### Get All Dossiers

```graphql
query {
  dossiers {
    id
    name
    deliveryTime
    frequency
    timezone
    tone
    language
    emailTo
    isActive
    feeds {
      id
      url
      title
    }
    deliveries(limit: 5) {
      deliveredAt
      status
      articleCount
    }
    createdAt
  }
}
```

**Returns:** All configured dossiers with their associated feeds and recent deliveries

### Get Specific Dossier

```graphql
query GetDossier($id: ID!) {
  dossier(id: $id) {
    id
    name
    deliveryTime
    frequency
    timezone
    tone
    language
    specialInstructions
    emailTo
    isActive
    feeds {
      id
      url
      title
      description
      lastFetchedAt
    }
    deliveries {
      id
      deliveredAt
      status
      articleCount
      emailContent
    }
    createdAt
    updatedAt
  }
}
```

**Parameters:**

- `id`: Dossier ID

**Returns:** Detailed dossier information including full delivery history

### Get All Feeds

```graphql
query {
  feeds {
    id
    url
    title
    description
    lastFetchedAt
    dossiers {
      id
      name
    }
    createdAt
    updatedAt
  }
}
```

**Returns:** All RSS feeds in the system with associated dossier information

### Get Articles

```graphql
query GetArticles($limit: Int, $offset: Int, $feedId: ID) {
  articles(limit: $limit, offset: $offset, feedId: $feedId) {
    id
    title
    link
    description
    author
    publishedAt
    feed {
      title
      url
    }
    createdAt
  }
}
```

**Parameters:**

- `limit`: Number of articles to return (optional, default: 50)
- `offset`: Pagination offset (optional, default: 0)
- `feedId`: Filter by specific feed (optional)

**Returns:** Paginated list of articles, optionally filtered by feed

### Search Articles

```graphql
query SearchArticles($query: String!, $limit: Int) {
  searchArticles(query: $query, limit: $limit) {
    id
    title
    link
    description
    publishedAt
    feed {
      title
    }
  }
}
```

**Parameters:**

- `query`: Search term for article titles and descriptions
- `limit`: Maximum results (optional, default: 20)

**Returns:** Articles matching the search query
}

````

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
````

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

### Create Dossier

```graphql
mutation CreateDossier($input: DossierInput!) {
  createDossier(input: $input) {
    id
    name
    deliveryTime
    frequency
    timezone
    tone
    language
    emailTo
    isActive
    feeds {
      id
      url
      title
    }
    createdAt
  }
}
```

**Parameters:**

- `input`: DossierInput object with dossier configuration

**Example Variables:**

```json
{
  "input": {
    "name": "Tech News Daily",
    "deliveryTime": "08:00",
    "frequency": "DAILY",
    "timezone": "America/New_York",
    "tone": "professional",
    "language": "english",
    "emailTo": "user@example.com",
    "feedUrls": [
      "https://news.ycombinator.com/rss",
      "https://techcrunch.com/feed/"
    ]
  }
}
```

**Returns:** Created dossier with associated feeds

### Update Dossier

```graphql
mutation UpdateDossier($id: ID!, $input: DossierInput!) {
  updateDossier(id: $id, input: $input) {
    id
    name
    deliveryTime
    frequency
    tone
    emailTo
    isActive
    updatedAt
  }
}
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

**Parameters:**

- `feedIds`: Optional array of specific feed IDs to refresh (if omitted, refreshes all)

**Returns:** Refresh operation results with success count and any errors

## Error Handling

The API returns errors in the standard GraphQL error format:

```json
{
  "errors": [
    {
      "message": "Feed not found",
      "path": ["addFeedToDossier"]
    }
  ],
  "data": null
}
```

Common error messages:

- `Dossier not found`: Invalid dossier ID
- `Feed not found`: Invalid feed ID
- `Invalid RSS URL`: Feed URL is not accessible or valid
- `Email delivery failed`: SMTP configuration or network issue
- `AI processing failed`: Ollama service unavailable
- `Invalid time format`: Delivery time must be HH:MM format

## AI Tone Options

Available tone values for dossier configuration:

- `professional`: Standard business communication style
- `humorous`: Witty and entertaining summaries
- `analytical`: Data-driven insights and trends
- `casual`: Relaxed, conversational tone
- `apocalyptic`: Dramatic, foreboding style with biblical references
- `orc`: Warcraft-style blunt communication
- `robot`: Mechanical, technical language
- `southern_belle`: Polite, charming Southern style
- `apologetic`: Sympathetic and reassuring
- `sweary`: Adult language (requires uncensored model)

## Frequency Options

- `DAILY`: Delivers every day at specified time
- `WEEKLY`: Delivers once per week on current day
- `MONTHLY`: Delivers once per month on current date

## Timezone Support

Uses IANA timezone database format:

- `America/New_York`
- `Europe/London`
- `Asia/Tokyo`
- `UTC`

## Examples

### Complete Dossier Setup

```graphql
# 1. Create dossier with feeds
mutation {
  createDossier(
    input: {
      name: "Morning Tech News"
      deliveryTime: "08:00"
      frequency: DAILY
      timezone: "America/New_York"
      tone: "professional"
      language: "english"
      emailTo: "user@example.com"
      feedUrls: [
        "https://news.ycombinator.com/rss"
        "https://techcrunch.com/feed/"
      ]
    }
  ) {
    id
    name
    feeds {
      id
      title
    }
  }
}

# 2. Test email delivery
mutation {
  testDossier(id: "1") {
    success
    message
  }
}

# 3. Add more feeds later
mutation {
  addFeedToDossier(
    dossierId: "1"
    url: "https://www.theverge.com/rss/index.xml"
  ) {
    id
    title
  }
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
