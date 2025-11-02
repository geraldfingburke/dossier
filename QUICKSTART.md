# Quick Start Guide

This guide will help you get Dossier up and running in under 5 minutes using Docker Compose.

## Prerequisites

- Docker and Docker Compose installed
- SMTP email credentials (Gmail, Outlook, etc.) for email delivery
- (Optional) OpenAI API key for AI-powered summaries (uses free local LLM by default)

## Steps

### 1. Clone the Repository

```bash
git clone https://github.com/geraldfingburke/dossier.git
cd dossier
```

### 2. Configure SMTP for Email Delivery

Dossier sends automated email dossiers, so you need to configure SMTP. Use one of the provided setup scripts:

**For Linux/macOS:**

```bash
chmod +x setup-smtp.sh
./setup-smtp.sh
```

**For Windows PowerShell:**

```powershell
.\setup-smtp.ps1
```

Or manually create a `.env` file from the example:

```bash
cp .env.example .env
# Edit .env and configure your SMTP settings
```

### 3. Start the Application

```bash
docker-compose up -d
```

This will start:

- PostgreSQL database on port 5432
- Ollama (local AI) service for free AI summaries
- Backend server on http://localhost:8080
- Frontend on http://localhost:5173

### 4. Access the Application

Open your browser and navigate to:

```
http://localhost:5173
```

### 5. Create Your First Dossier Configuration

1. Click "Add New Dossier"
2. Configure your dossier:
   - **Title**: Name for your dossier (e.g., "Tech News Daily")
   - **Email**: Where to send the dossier
   - **RSS Feeds**: Add feeds you want to monitor
   - **Article Count**: How many articles to include (1-50)
   - **Frequency**: Daily, Weekly, or Monthly
   - **Delivery Time**: When to send the dossier
   - **Tone**: Professional, Humorous, Analytical, etc.
   - **Language**: English, Spanish, French, etc.
3. Click "Save Dossier"

### 6. Test Your Configuration

1. Click the "Test Email" button on your dossier
2. Check your email to confirm delivery works
3. The dossier will be automatically generated and sent at your scheduled time

That's it! You're now using Dossier for automated RSS email summaries.

## Useful Commands

### View Logs

```bash
docker-compose logs -f backend    # Backend logs
docker-compose logs -f ollama     # AI service logs
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

### Test Email Configuration

```bash
# Use the setup scripts to reconfigure SMTP
./setup-smtp.sh      # Linux/macOS
.\setup-smtp.ps1     # Windows PowerShell
```

## Troubleshooting

### Email Not Sending

1. **Check SMTP Configuration**: Run the setup script again
2. **Gmail Users**: Make sure you're using an App Password, not your regular password
   - Generate at: https://myaccount.google.com/apppasswords
3. **Check Logs**: `docker-compose logs backend` for error messages
4. **Test Email Button**: Use the test button in the web interface

### Scheduler Not Running

1. **Check Backend Logs**: `docker-compose logs backend`
2. **Verify Timezone**: Make sure your delivery time accounts for timezone settings
3. **Check Frequency**: Daily dossiers only send once per day

### Port Already in Use

If you get an error that a port is already in use, you can change the ports in `docker-compose.yml`:

```yaml
services:
  backend:
    ports:
      - "8081:8080" # Change first number to any available port

  frontend:
    ports:
      - "3000:5173" # Change first number to any available port
```

### Database Connection Issues

If the backend can't connect to the database, try:

```bash
docker-compose down -v  # Remove volumes
docker-compose up -d
```

### AI Service Issues

If Ollama (AI service) isn't working:

```bash
docker-compose logs ollama        # Check AI service logs
docker-compose restart ollama     # Restart AI service
```

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

- **Create Multiple Dossiers**: Set up different dossiers for different topics (tech, news, etc.)
- **Try Different Tones**: Experiment with humorous, analytical, or other creative tones
- **Schedule Optimization**: Fine-tune delivery times and frequencies
- **Advanced Configuration**: Use special instructions to customize AI behavior
- **Multiple Languages**: Configure dossiers in different languages
- **Explore API**: Check out the GraphQL API at http://localhost:8080/graphql/playground

## Features

- ✅ **Automated Email Delivery**: Scheduled dossiers sent directly to your inbox
- ✅ **Local AI Processing**: Free AI summaries using Ollama (no OpenAI required)
- ✅ **Multiple Tones**: Professional, humorous, analytical, apocalyptic, and more
- ✅ **Multi-language Support**: Generate dossiers in any language
- ✅ **Flexible Scheduling**: Daily, weekly, or monthly delivery
- ✅ **Custom Instructions**: Fine-tune AI behavior with special instructions
- ✅ **Single-User Design**: No accounts needed, perfect for self-hosting

For more information, see the main [README.md](README.md) and [SMTP_SETUP.md](SMTP_SETUP.md).
