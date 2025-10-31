# Quick Start Guide

This guide will help you get Dossier up and running in under 5 minutes using Docker Compose.

## Prerequisites

- Docker and Docker Compose installed
- (Optional) OpenAI API key for AI-powered summaries

## Steps

### 1. Clone the Repository

```bash
git clone https://github.com/geraldfingburke/dossier.git
cd dossier
```

### 2. Configure Environment (Optional)

If you have an OpenAI API key, create a `.env` file:

```bash
cp .env.example .env
# Edit .env and add your OPENAI_API_KEY
```

If you don't have an OpenAI API key, the app will work with mock summaries.

### 3. Start the Application

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Backend server on http://localhost:8080
- Frontend on http://localhost:5173

### 4. Access the Application

Open your browser and navigate to:
```
http://localhost:5173
```

### 5. Create an Account

1. Click "Register" on the login page
2. Enter your email, name, and password
3. You'll be automatically logged in

### 6. Add Your First RSS Feed

1. Go to the "Feeds" tab
2. Enter an RSS feed URL (e.g., `https://news.ycombinator.com/rss`)
3. Click "Add Feed"
4. Click "Refresh All Feeds" to fetch articles

### 7. Generate Your First Digest

1. Go to the "Digests" tab
2. Click "Generate New Digest"
3. Wait a few seconds for the AI to summarize your articles

That's it! You're now using Dossier.

## Useful Commands

### View Logs
```bash
docker-compose logs -f
```

### Stop the Application
```bash
docker-compose down
```

### Rebuild After Changes
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Troubleshooting

### Port Already in Use

If you get an error that a port is already in use, you can change the ports in `docker-compose.yml`:

```yaml
services:
  backend:
    ports:
      - "8081:8080"  # Change first number to any available port
  
  frontend:
    ports:
      - "3000:5173"  # Change first number to any available port
```

### Database Connection Issues

If the backend can't connect to the database, try:

```bash
docker-compose down -v  # Remove volumes
docker-compose up -d
```

### Frontend Not Loading

Make sure you're accessing the correct URL (http://localhost:5173 by default) and that the backend is running on port 8080.

## Popular RSS Feed Examples

Here are some popular feeds to get you started:

**Tech News:**
- Hacker News: `https://news.ycombinator.com/rss`
- TechCrunch: `https://techcrunch.com/feed/`
- The Verge: `https://www.theverge.com/rss/index.xml`
- Ars Technica: `https://feeds.arstechnica.com/arstechnica/index`

**Programming:**
- DEV.to: `https://dev.to/feed`
- CSS-Tricks: `https://css-tricks.com/feed/`
- Smashing Magazine: `https://www.smashingmagazine.com/feed/`

**General News:**
- BBC News: `http://feeds.bbci.co.uk/news/rss.xml`
- NPR: `https://feeds.npr.org/1001/rss.xml`
- Reuters: `https://www.reutersagency.com/feed/`

## Next Steps

- Explore the GraphQL API at http://localhost:8080/graphql/playground
- Add more RSS feeds
- Set up automatic daily digests (runs automatically every 24 hours)
- Customize the application to your needs

For more information, see the main [README.md](README.md).
