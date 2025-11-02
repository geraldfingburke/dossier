# SMTP Setup Guide for Dossier

This guide will help you configure SMTP settings for email delivery in Dossier.

## Quick Setup Options

### Option 1: Gmail (Recommended for personal use)

1. **Enable 2-Factor Authentication** on your Gmail account
2. **Generate an App Password**:
   - Go to [Google Account Settings](https://myaccount.google.com/apppasswords)
   - Select "Mail" and generate a password
3. **Update your `.env` file**:
   ```bash
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-16-character-app-password
   SMTP_FROM_EMAIL=your-email@gmail.com
   SMTP_FROM_NAME=Dossier
   ```

### Option 2: Outlook/Hotmail

1. **Update your `.env` file**:
   ```bash
   SMTP_HOST=smtp-mail.outlook.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@outlook.com
   SMTP_PASSWORD=your-password
   SMTP_FROM_EMAIL=your-email@outlook.com
   SMTP_FROM_NAME=Dossier
   ```

### Option 3: Yahoo Mail

1. **Enable 2-Factor Authentication** on your Yahoo account
2. **Generate an App Password**:
   - Go to Yahoo Account Security settings
   - Generate an app password for "Mail"
3. **Update your `.env` file**:
   ```bash
   SMTP_HOST=smtp.mail.yahoo.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@yahoo.com
   SMTP_PASSWORD=your-app-password
   SMTP_FROM_EMAIL=your-email@yahoo.com
   SMTP_FROM_NAME=Dossier
   ```

### Option 4: Local Testing with MailHog (No real emails sent)

Perfect for testing without sending actual emails:

1. **Uncomment MailHog service** in `docker-compose.yml`:

   ```yaml
   mailhog:
     image: mailhog/mailhog:latest
     ports:
       - "1025:1025" # SMTP server
       - "8025:8025" # Web interface
     restart: unless-stopped
   ```

2. **Update your `.env` file**:

   ```bash
   SMTP_HOST=mailhog
   SMTP_PORT=1025
   SMTP_USERNAME=
   SMTP_PASSWORD=
   SMTP_FROM_EMAIL=dossier@localhost
   SMTP_FROM_NAME=Dossier
   ```

3. **View emails** at http://localhost:8025

## Setup Steps

1. **Copy the example environment file**:

   ```bash
   cp .env.example .env
   ```

2. **Edit the `.env` file** with your preferred SMTP configuration from above

3. **Restart the application**:

   ```bash
   docker-compose down
   docker-compose up -d
   ```

4. **Test the email connection**:
   - Go to http://localhost:5173
   - Click "Test Email" button in the header
   - Or create a dossier configuration and click "Test Dossier"

## Troubleshooting

### Common Issues

1. **"Connection refused" error**:

   - Check that your SMTP host and port are correct
   - Ensure you're not behind a firewall blocking SMTP ports

2. **Authentication failed**:

   - For Gmail/Yahoo: Make sure you're using an App Password, not your regular password
   - Double-check your username and password

3. **TLS/SSL errors**:
   - Most modern email providers require TLS on port 587
   - Port 25 is often blocked by ISPs

### Testing with MailHog

If you want to test without sending real emails:

1. Enable MailHog service in docker-compose.yml
2. Set SMTP_HOST=mailhog and SMTP_PORT=1025
3. View caught emails at http://localhost:8025

## Security Notes

- **Never commit your `.env` file** with real credentials to version control
- **Use App Passwords** instead of your main email password when possible
- **Consider using a dedicated email account** for your dossier service
- **In production**, use environment variables or a secure secrets management system

## Need Help?

- Check the application logs: `docker-compose logs backend`
- Test SMTP connection using the "Test Email" button in the UI
- Make sure your email provider allows SMTP access (some require enabling "Less secure app access")
