# Dossier GraphQL API Documentation

This document provides comprehensive information about the Dossier system's GraphQL API.

## Base URLs

- **Development**: `http://localhost:8080/graphql`
- **Production**: Configure via environment variables

## Authentication

The Dossier system uses a **single-user design** with no authentication required. All mutations and queries are accessible without tokens or credentials.

## Core Concepts

- **DossierConfig**: A configuration for automated email delivery with RSS feeds, schedule, and AI settings
- **Dossier**: A historical record of a generated and sent digest email
- **Article**: Individual news items collected from RSS feeds
- **Tone**: AI personality/style settings for content generation (system defaults + custom)

## GraphQL Schema

### Core Types

#### DossierConfig

```graphql
type DossierConfig {
  id: ID!
  title: String! # Display name (e.g., "Tech News Daily")
  email: String! # Delivery email address
  feedUrls: [String!]! # RSS/Atom feed URLs
  articleCount: Int! # Number of articles to include per digest
  frequency: String! # "daily", "weekly", or "monthly"
  deliveryTime: String! # HH:MM format (24-hour)
  timezone: String! # IANA timezone (e.g., "America/New_York")
  tone: String # AI tone name (references Tone.name)
  language: String # Summary language (e.g., "english", "spanish")
  specialInstructions: String # Custom AI instructions
  active: Boolean! # Whether scheduler processes this config
  createdAt: String!
}
```

#### Dossier

```graphql
type Dossier {
  id: ID!
  configId: ID! # Reference to DossierConfig
  subject: String! # Email subject line
  content: String! # Generated HTML email content
  sentAt: String! # Timestamp when email was sent
}
```

#### Article

```graphql
type Article {
  id: ID!
  title: String! # Article headline
  link: String! # Original article URL (unique)
  description: String # Article summary/excerpt
  content: String # Full article content
  author: String # Article author
  publishedAt: String! # Original publication date
  createdAt: String! # When article was fetched
}
```

#### Tone

```graphql
type Tone {
  id: ID!
  name: String! # Unique tone identifier (e.g., "professional")
  prompt: String! # AI system prompt for this tone
  isSystemDefault: Boolean! # Whether this is a built-in tone
  createdAt: String!
  updatedAt: String!
}
```

#### SchedulerStatus

```graphql
type SchedulerStatus {
  running: Boolean! # Whether scheduler is active
  nextCheck: String # Timestamp of next scheduler tick
  activeDossiers: Int! # Count of active DossierConfigs
}
```

### Input Types

#### DossierConfigInput

```graphql
input DossierConfigInput {
  title: String!
  email: String!
  feedUrls: [String!]!
  articleCount: Int!
  frequency: String! # "daily", "weekly", or "monthly"
  deliveryTime: String! # HH:MM format (24-hour)
  timezone: String! # IANA timezone
  tone: String # Tone name (optional)
  language: String # Summary language (optional)
  specialInstructions: String # Custom AI instructions (optional)
}
```

#### ToneInput

```graphql
input ToneInput {
  name: String! # Unique tone identifier
  prompt: String! # AI system prompt
}
```

## Queries

### Get All Dossier Configs

```graphql
query {
  dossierConfigs {
    id
    title
    email
    feedUrls
    articleCount
    frequency
    deliveryTime
    timezone
    tone
    language
    specialInstructions
    active
    createdAt
  }
}
```

**Returns:** All dossier configurations

### Get Single Dossier Config

```graphql
query GetDossierConfig($id: ID!) {
  dossierConfig(id: $id) {
    id
    title
    email
    feedUrls
    articleCount
    frequency
    deliveryTime
    timezone
    tone
    language
    specialInstructions
    active
    createdAt
  }
}
```

**Parameters:**

- `id`: DossierConfig ID

**Returns:** Specific dossier configuration or null if not found

### Get Dossier History

```graphql
query GetDossiers($configId: ID, $limit: Int) {
  dossiers(configId: $configId, limit: $limit) {
    id
    configId
    subject
    content
    sentAt
  }
}
```

**Parameters:**

- `configId`: Filter by specific DossierConfig (optional)
- `limit`: Maximum number of dossiers to return (optional)

**Returns:** Historical records of generated and sent dossiers

### Get Scheduler Status

```graphql
query {
  schedulerStatus {
    running
    nextCheck
    activeDossiers
  }
}
```

**Returns:** Current scheduler state and statistics

### Get All Tones

```graphql
query {
  tones {
    id
    name
    prompt
    isSystemDefault
    createdAt
    updatedAt
  }
}
```

**Returns:** All available AI tones (system defaults + custom)

### Get Single Tone

```graphql
query GetTone($id: ID!) {
  tone(id: $id) {
    id
    name
    prompt
    isSystemDefault
    createdAt
    updatedAt
  }
}
```

**Parameters:**

- `id`: Tone ID

**Returns:** Specific tone or null if not found

## Mutations

### Create Dossier Config

```graphql
mutation CreateDossierConfig($input: DossierConfigInput!) {
  createDossierConfig(input: $input) {
    id
    title
    email
    feedUrls
    frequency
    deliveryTime
    timezone
    tone
    active
    createdAt
  }
}
```

**Parameters:**

- `input`: DossierConfigInput object with configuration

**Example Variables:**

```json
{
  "input": {
    "title": "Tech News Daily",
    "email": "user@example.com",
    "feedUrls": [
      "https://news.ycombinator.com/rss",
      "https://techcrunch.com/feed/"
    ],
    "articleCount": 15,
    "frequency": "daily",
    "deliveryTime": "08:00",
    "timezone": "America/New_York",
    "tone": "professional",
    "language": "english"
  }
}
```

**Returns:** Created dossier configuration

### Update Dossier Config

```graphql
mutation UpdateDossierConfig($id: ID!, $input: DossierConfigInput!) {
  updateDossierConfig(id: $id, input: $input) {
    id
    title
    email
    frequency
    deliveryTime
    active
  }
}
```

**Parameters:**

- `id`: DossierConfig ID to update
- `input`: DossierConfigInput with new values

**Returns:** Updated dossier configuration

### Delete Dossier Config

```graphql
mutation DeleteDossierConfig($id: ID!) {
  deleteDossierConfig(id: $id)
}
```

**Parameters:**

- `id`: DossierConfig ID to delete

**Returns:** Boolean indicating success

### Toggle Dossier Config Active State

```graphql
mutation ToggleDossierConfig($id: ID!, $active: Boolean!) {
  toggleDossierConfig(id: $id, active: $active) {
    id
    active
  }
}
```

**Parameters:**

- `id`: DossierConfig ID
- `active`: New active state

**Returns:** Updated dossier configuration

### Generate and Send Dossier (Manual Trigger)

```graphql
mutation GenerateAndSendDossier($configId: ID!) {
  generateAndSendDossier(configId: $configId) {
    id
    configId
    subject
    content
    sentAt
  }
}
```

**Parameters:**

- `configId`: DossierConfig ID to generate and send

**Returns:** Generated dossier with email content

**Note:** This manually triggers dossier generation, bypassing the scheduler

### Send Test Email

```graphql
mutation SendTestEmail($configId: ID!) {
  sendTestEmail(configId: $configId)
}
```

**Parameters:**

- `configId`: DossierConfig ID to test

**Returns:** Boolean indicating if test email was sent successfully

**Note:** Generates a test dossier and sends it to the configured email address

### Test Email Connection

```graphql
mutation TestEmailConnection(
  $email: String!
  $smtpHost: String!
  $smtpPort: Int!
  $smtpUser: String!
  $smtpPass: String!
) {
  testEmailConnection(
    email: $email
    smtpHost: $smtpHost
    smtpPort: $smtpPort
    smtpUser: $smtpUser
    smtpPass: $smtpPass
  )
}
```

**Parameters:**

- `email`: Test recipient email
- `smtpHost`: SMTP server hostname
- `smtpPort`: SMTP server port
- `smtpUser`: SMTP username
- `smtpPass`: SMTP password

**Returns:** Boolean indicating if SMTP connection and email send succeeded

**Note:** Tests raw SMTP credentials without creating a dossier config

### Create Custom Tone

```graphql
mutation CreateTone($input: ToneInput!) {
  createTone(input: $input) {
    id
    name
    prompt
    isSystemDefault
    createdAt
  }
}
```

**Parameters:**

- `input`: ToneInput object with name and prompt

**Example Variables:**

```json
{
  "input": {
    "name": "poetic",
    "prompt": "Write in a lyrical, poetic style with metaphors and elegant prose"
  }
}
```

**Returns:** Created tone

### Update Tone

```graphql
mutation UpdateTone($id: ID!, $input: ToneInput!) {
  updateTone(id: $id, input: $input) {
    id
    name
    prompt
    updatedAt
  }
}
```

**Parameters:**

- `id`: Tone ID to update
- `input`: ToneInput with new values

**Returns:** Updated tone

**Note:** Cannot update system default tones

### Delete Tone

```graphql
mutation DeleteTone($id: ID!) {
  deleteTone(id: $id)
}
```

**Parameters:**

- `id`: Tone ID to delete

**Returns:** Boolean indicating success

**Note:** Cannot delete system default tones

## Error Handling

The API returns errors in the standard GraphQL error format:

```json
{
  "errors": [
    {
      "message": "dossier config not found",
      "path": ["dossierConfig"]
    }
  ],
  "data": null
}
```

Common error messages:

- `dossier config not found`: Invalid DossierConfig ID
- `tone not found`: Invalid Tone ID
- `cannot delete system default tone`: Attempted to delete built-in tone
- `cannot update system default tone`: Attempted to modify built-in tone
- `invalid email format`: Email address is malformed
- `invalid time format`: Delivery time must be HH:MM format (24-hour)
- `invalid timezone`: Timezone is not a valid IANA timezone
- `invalid frequency`: Frequency must be "daily", "weekly", or "monthly"
- `failed to fetch RSS feed`: One or more feed URLs are inaccessible
- `AI generation failed`: Ollama service error or model unavailable
- `email delivery failed`: SMTP configuration issue or network error

## System Default Tones

The following tones are provided by default and cannot be modified or deleted:

1. **professional**: Standard business communication style
2. **humorous**: Witty and entertaining summaries
3. **analytical**: Data-driven insights and trends analysis
4. **casual**: Relaxed, conversational tone
5. **apocalyptic**: Dramatic, foreboding style with biblical references
6. **orc**: Warcraft-style blunt communication ("Me Grognak!")
7. **robot**: Mechanical, technical language ("EXECUTING SUMMARY PROTOCOL")
8. **southern_belle**: Polite, charming Southern style ("Well, bless your heart")
9. **apologetic**: Sympathetic and reassuring ("I'm so sorry to report...")
10. **sweary**: Adult language (requires uncensored LLM model)

Custom tones can be created via the `createTone` mutation.

## Frequency Options

- `daily`: Delivers every day at the specified time
- `weekly`: Delivers once per week (same day of week)
- `monthly`: Delivers once per month (same day of month)

**Note:** Frequencies are case-insensitive strings, not enums

## Timezone Support

Uses IANA timezone database format. Examples:

- `America/New_York` - Eastern Time (US)
- `America/Chicago` - Central Time (US)
- `America/Denver` - Mountain Time (US)
- `America/Los_Angeles` - Pacific Time (US)
- `Europe/London` - UK
- `Europe/Paris` - Central European Time
- `Asia/Tokyo` - Japan
- `UTC` - Coordinated Universal Time

## AI Generation

Dossier uses **Ollama** for local LLM-based content generation:

- **Model**: llama3.2:3b (default) or dolphin-mistral (uncensored, for sweary tone)
- **Temperature**: 0.7 (balanced creativity)
- **Max Tokens**: 2000 (for comprehensive summaries)
- **Endpoint**: Configurable via `OLLAMA_URL` environment variable

The AI summarization pipeline:

1. Fetch articles from configured RSS feeds
2. Filter to `articleCount` most recent articles
3. Format articles with title, description, link
4. Apply tone-specific system prompt
5. Apply language and special instructions
6. Generate markdown-formatted summary
7. Convert to HTML email template

## Scheduler Behavior

The scheduler runs continuously with these characteristics:

- **Granularity**: 1-minute ticker
- **Timezone-Aware**: Converts delivery times to UTC for comparison
- **Frequency Rules**:
  - Daily: Generates if current time matches delivery time
  - Weekly: Generates if current day matches last generation day + 7 days
  - Monthly: Generates if current day matches last generation day + 1 month
- **Duplicate Prevention**: Tracks last generation time per config
- **Concurrency**: Processes each dossier in separate goroutine
- **Error Resilience**: Individual failures don't stop scheduler

## Email Delivery

Email delivery uses SMTP with TLS encryption:

- **Configuration**: Via environment variables
  - `SMTP_HOST`: SMTP server hostname
  - `SMTP_PORT`: SMTP server port (typically 587)
  - `SMTP_USER`: SMTP username
  - `SMTP_PASS`: SMTP password
  - `SMTP_FROM`: Sender email address
- **Format**: HTML emails with inline CSS
- **Template**: Professional layout with article cards

## RSS Feed Support

Supported feed formats:

- RSS 1.0
- RSS 2.0
- Atom 1.0

The RSS service handles:

- Multi-feed aggregation
- Duplicate detection (by link URL)
- Missing field handling (graceful degradation)
- Feed validation and error recovery
- Date parsing from multiple formats

## Complete Workflow Example

```graphql
# 1. Check available tones
query {
  tones {
    name
    prompt
  }
}

# 2. Create dossier configuration
mutation {
  createDossierConfig(
    input: {
      title: "Morning Tech Digest"
      email: "user@example.com"
      feedUrls: [
        "https://news.ycombinator.com/rss"
        "https://techcrunch.com/feed/"
        "https://www.theverge.com/rss/index.xml"
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

# 3. Test the configuration
mutation {
  sendTestEmail(configId: "1")
}

# 4. Check scheduler status
query {
  schedulerStatus {
    running
    activeDossiers
  }
}

# 5. View dossier history
query {
  dossiers(configId: "1", limit: 10) {
    id
    subject
    sentAt
  }
}

# 6. Manually trigger generation
mutation {
  generateAndSendDossier(configId: "1") {
    id
    subject
    content
    sentAt
  }
}

# 7. Toggle active state
mutation {
  toggleDossierConfig(id: "1", active: false) {
    id
    active
  }
}
```

## Environment Variables

Required for full functionality:

```bash
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/dossier

# SMTP Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
SMTP_FROM=your-email@gmail.com

# Ollama AI
OLLAMA_URL=http://localhost:11434

# Server
PORT=8080
```

## Database Schema

Key tables:

- `dossier_configs`: Dossier configuration and scheduling
- `dossiers`: Historical records of sent digests
- `articles`: Cached RSS articles
- `tones`: AI tone definitions

See [ARCHITECTURE.md](ARCHITECTURE.md) for complete schema details.

## Support

For issues or questions:

- Open an issue on GitHub: [geraldfingburke/dossier](https://github.com/geraldfingburke/dossier)
- Check [README.md](README.md) for project overview
- See [QUICKSTART.md](QUICKSTART.md) for setup instructions
- Review [ARCHITECTURE.md](ARCHITECTURE.md) for technical details
