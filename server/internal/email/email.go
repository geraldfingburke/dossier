// Package email provides SMTP email delivery services with TLS encryption support.
// It handles dossier email composition, HTML/text template rendering, and secure
// email transmission via SMTP servers using either STARTTLS (port 587) or direct
// TLS (port 465) connections.
//
// Key Features:
//   - Multi-part MIME emails (HTML + plain text)
//   - Beautiful responsive HTML templates
//   - TLS encryption for secure transmission
//   - Support for both STARTTLS and direct TLS
//   - Environment-based configuration
//   - Connection testing capabilities
package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/geraldfingburke/dossier/server/internal/models"
)

// ============================================================================
// TYPES AND CONFIGURATION
// ============================================================================

// Config holds SMTP server configuration for email delivery.
// All fields are populated from environment variables with sensible defaults.
type Config struct {
	SMTPHost     string // SMTP server hostname (e.g., "smtp.gmail.com")
	SMTPPort     string // SMTP server port ("587" for STARTTLS, "465" for direct TLS)
	Username     string // SMTP authentication username (usually email address)
	Password     string // SMTP authentication password or app-specific password
	FromEmail    string // Sender email address
	FromName     string // Display name for sender
}

// Service handles all email operations including template rendering and SMTP delivery.
type Service struct {
	config Config
}

// DossierEmail represents a complete email ready for delivery.
// Contains both HTML and plain text versions for maximum compatibility.
type DossierEmail struct {
	To          string      // Recipient email address
	Subject     string      // Email subject line
	HTMLBody    string      // HTML version of email body
	TextBody    string      // Plain text version of email body
	DossierData DossierData // Structured data for template rendering
}

// DossierData contains structured information for rendering dossier email templates.
// Used by both HTML and text templates to generate consistent content.
type DossierData struct {
	Title           string        // Dossier configuration title
	Summary         string        // AI-generated summary (HTML format)
	Articles        []ArticleData // List of articles included in dossier
	GeneratedAt     time.Time     // When the dossier was generated
	ArticleCount    int           // Number of articles included
	Tone            string        // AI tone used for summary
	Language        string        // Language of the summary
	Instructions    string        // Special instructions applied (if any)
}

// ArticleData represents a single article in the email template.
type ArticleData struct {
	Title       string    // Article headline
	Description string    // Article summary/excerpt
	URL         string    // Full article URL
	Source      string    // Domain name of source (extracted from URL)
	PublishedAt time.Time // Original publication date
}

// ============================================================================
// SERVICE INITIALIZATION
// ============================================================================

// NewService creates a new email service instance with configuration loaded
// from environment variables.
//
// Environment Variables:
//   - SMTP_HOST: SMTP server hostname (default: "localhost")
//   - SMTP_PORT: SMTP server port (default: "587")
//   - SMTP_USERNAME: Authentication username (default: "")
//   - SMTP_PASSWORD: Authentication password (default: "")
//   - SMTP_FROM_EMAIL: Sender email address (default: "dossier@localhost")
//   - SMTP_FROM_NAME: Sender display name (default: "Dossier")
//
// Port Selection Guide:
//   - 587: Use STARTTLS (upgrade plain connection to TLS)
//   - 465: Use direct TLS (TLS from connection start)
//   - 25: Plain SMTP (not recommended, no encryption)
//
// Returns:
//   - *Service: Configured email service ready for use
//
// Example:
//   emailService := NewService()
//   err := emailService.SendDossier(config, summary, articles)
func NewService() *Service {
	config := Config{
		SMTPHost:  getEnvOrDefault("SMTP_HOST", "localhost"),
		SMTPPort:  getEnvOrDefault("SMTP_PORT", "587"),
		Username:  getEnvOrDefault("SMTP_USERNAME", ""),
		Password:  getEnvOrDefault("SMTP_PASSWORD", ""),
		FromEmail: getEnvOrDefault("SMTP_FROM_EMAIL", "dossier@localhost"),
		FromName:  getEnvOrDefault("SMTP_FROM_NAME", "Dossier"),
	}

	return &Service{config: config}
}

// getEnvOrDefault retrieves an environment variable value or returns a default.
//
// Parameters:
//   - key: Environment variable name
//   - defaultValue: Fallback value if variable is not set or empty
//
// Returns:
//   - string: Environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ============================================================================
// PUBLIC API - EMAIL DELIVERY
// ============================================================================

// SendDossier composes and sends a complete dossier email.
// This is the main entry point for email delivery operations.
//
// Process Flow:
//   1. Transform articles into email-friendly data structures
//   2. Extract domain names from URLs for source attribution
//   3. Compile dossier data with metadata
//   4. Generate HTML and plain text versions via templates
//   5. Build MIME multi-part message
//   6. Send via SMTP with TLS encryption
//
// The resulting email includes:
//   - Styled HTML version with responsive design
//   - Plain text fallback for simple email clients
//   - Article previews with links
//   - AI-generated summary
//   - Metadata (generation time, article count, tone, etc.)
//
// Parameters:
//   - config: Dossier configuration (recipient, title, preferences)
//   - summary: AI-generated HTML summary of articles
//   - articles: List of articles to include in email
//
// Returns:
//   - error: Template rendering or SMTP delivery failure
//
// Example:
//   err := emailService.SendDossier(dossierConfig, aiSummary, articles)
//   if err != nil {
//       log.Printf("Failed to send dossier: %v", err)
//   }
func (s *Service) SendDossier(config *models.DossierConfig, summary string, articles []models.Article) error {
	log.Printf("Preparing to send dossier email: %s to %s", config.Title, config.Email)

	// Transform articles into email data structures
	articleData := make([]ArticleData, len(articles))
	for i, article := range articles {
		articleData[i] = ArticleData{
			Title:       article.Title,
			Description: article.Description,
			URL:         article.Link,
			Source:      extractDomain(article.Link),
			PublishedAt: article.PublishedAt,
		}
	}

	// Compile dossier data for template rendering
	dossierData := DossierData{
		Title:           config.Title,
		Summary:         summary,
		Articles:        articleData,
		GeneratedAt:     time.Now(),
		ArticleCount:    len(articles),
		Tone:            config.Tone,
		Language:        config.Language,
		Instructions:    config.SpecialInstructions,
	}

	// Generate HTML and text email content
	htmlBody, textBody, err := s.generateEmailContent(dossierData)
	if err != nil {
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	// Create email structure
	email := DossierEmail{
		To:          config.Email,
		Subject:     fmt.Sprintf("Dossier - %s", config.Title),
		HTMLBody:    htmlBody,
		TextBody:    textBody,
		DossierData: dossierData,
	}

	// Send via SMTP
	return s.sendEmail(email)
}

// TestSMTPConnection validates SMTP configuration by attempting authentication.
// This is useful for configuration verification before sending actual emails.
//
// Connection Strategy:
//   - Port 587: Use STARTTLS (connect plain, upgrade to TLS)
//   - Port 465: Use direct TLS (TLS from start)
//   - Other ports: Attempt STARTTLS as fallback
//
// Validation Steps:
//   1. Connect to SMTP server
//   2. Establish TLS encryption
//   3. Authenticate with credentials
//   4. Close connection
//
// Returns:
//   - error: Connection, TLS, or authentication failure
//
// Example:
//   if err := emailService.TestSMTPConnection(); err != nil {
//       log.Fatal("SMTP configuration invalid:", err)
//   }
func (s *Service) TestSMTPConnection() error {
	log.Printf("Testing SMTP connection to %s:%s", s.config.SMTPHost, s.config.SMTPPort)
	
	addr := s.config.SMTPHost + ":" + s.config.SMTPPort
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)

	// Select connection method based on port
	if s.config.SMTPPort == "587" {
		return s.testWithSTARTTLS(auth, addr)
	}
	return s.testWithDirectTLS(auth, addr)
}

// ============================================================================
// TEMPLATE RENDERING
// ============================================================================

// generateEmailContent creates both HTML and plain text versions of the email
// using Go templates. Both versions contain the same information but with
// appropriate formatting for their medium.
//
// HTML Version Features:
//   - Responsive design (mobile-friendly)
//   - Gradient header with branding
//   - Styled article cards with hover effects
//   - Inline CSS for maximum email client compatibility
//   - Embedded links in formatted text
//
// Plain Text Version Features:
//   - Clean ASCII formatting
//   - Section separators for readability
//   - Numbered article list
//   - All information from HTML version
//
// Template Functions:
//   - title: Capitalizes first letter of each word
//   - nl2br: Converts newlines to <br> tags for HTML
//   - add: Addition for template math (e.g., array indexing)
//
// Parameters:
//   - data: Structured dossier data for template rendering
//
// Returns:
//   - string: HTML version of email
//   - string: Plain text version of email
//   - error: Template parsing or execution failure
func (s *Service) generateEmailContent(data DossierData) (string, string, error) {
	// HTML template
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            max-width: 800px; 
            margin: 0 auto; 
            padding: 20px; 
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
            color: white; 
            padding: 30px; 
            border-radius: 10px; 
            margin-bottom: 30px; 
            text-align: center; 
        }
        .header h1 { margin: 0; font-size: 2em; }
        .meta { 
            background: #f8f9fa; 
            padding: 15px; 
            border-radius: 8px; 
            margin-bottom: 25px; 
            font-size: 0.9em; 
            color: #666; 
        }
        .summary { 
            background: white; 
            border-left: 4px solid #667eea; 
            padding: 20px; 
            margin-bottom: 30px; 
            border-radius: 0 8px 8px 0; 
        }
        .articles { margin-bottom: 30px; }
        .article { 
            border: 1px solid #e9ecef; 
            border-radius: 8px; 
            padding: 20px; 
            margin-bottom: 15px; 
            transition: box-shadow 0.2s; 
        }
        .article:hover { box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .article-title { 
            font-size: 1.2em; 
            font-weight: 600; 
            margin-bottom: 8px; 
        }
        .article-title a { color: #333; text-decoration: none; }
        .article-title a:hover { color: #667eea; }
        .article-meta { 
            font-size: 0.85em; 
            color: #666; 
            margin-bottom: 10px; 
        }
        .article-description { color: #555; }
        .footer { 
            text-align: center; 
            padding: 20px; 
            border-top: 1px solid #e9ecef; 
            margin-top: 30px; 
            font-size: 0.9em; 
            color: #666; 
        }
        .btn { 
            display: inline-block; 
            background: #667eea; 
            color: white; 
            padding: 10px 20px; 
            text-decoration: none; 
            border-radius: 5px; 
            font-weight: 500; 
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>📰 {{.Title}}</h1>
        <p>Your personalized news dossier</p>
    </div>

    <div class="meta">
        <strong>Generated:</strong> {{.GeneratedAt.Format "Monday, January 2, 2006 at 3:04 PM"}} | 
        <strong>Articles:</strong> {{.ArticleCount}} | 
        <strong>Style:</strong> {{.Tone | title}} {{.Language}}
        {{if .Instructions}}<br><strong>Special Instructions:</strong> {{.Instructions}}{{end}}
    </div>

    <div class="summary">
        <h2>🔍 Executive Summary</h2>
        {{.Summary | nl2br}}
    </div>

    <div class="articles">
        <h2>📖 Articles</h2>
        {{range $index, $article := .Articles}}
        <div class="article">
            <div class="article-title">
                <a href="{{$article.URL}}" target="_blank">{{$article.Title}}</a>
            </div>
            <div class="article-meta">
                <strong>Source:</strong> {{$article.Source}} | 
                <strong>Published:</strong> {{$article.PublishedAt.Format "Jan 2, 2006"}}
            </div>
            {{if $article.Description}}
            <div class="article-description">
                {{$article.Description}}
            </div>
            {{end}}
        </div>
        {{end}}
    </div>

    <div class="footer">
        <p>This dossier was automatically generated by <strong>Dossier</strong></p>
        <p>Delivered with ❤️ from your personal news automation system</p>
    </div>
</body>
</html>`

	// Text template
	textTemplate := `
{{.Title}}
==============================================

Generated: {{.GeneratedAt.Format "Monday, January 2, 2006 at 3:04 PM"}}
Articles: {{.ArticleCount}} | Style: {{.Tone | title}} {{.Language}}
{{if .Instructions}}Special Instructions: {{.Instructions}}{{end}}

EXECUTIVE SUMMARY
----------------------------------------------
{{.Summary}}

ARTICLES
----------------------------------------------
{{range $index, $article := .Articles}}
{{add $index 1}}. {{$article.Title}}
   Source: {{$article.Source}} | Published: {{$article.PublishedAt.Format "Jan 2, 2006"}}
   {{if $article.Description}}{{$article.Description}}{{end}}
   Read more: {{$article.URL}}

{{end}}

----------------------------------------------
This dossier was automatically generated by Dossier
Delivered from your personal news automation system
`

	// Create template functions
	funcMap := template.FuncMap{
		"title": strings.Title,
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(text, "\n", "<br>"))
		},
		"add": func(a, b int) int {
			return a + b
		},
	}

	// Parse and execute HTML template
	htmlTmpl, err := template.New("html").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := htmlTmpl.Execute(&htmlBuf, data); err != nil {
		return "", "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	// Parse and execute text template
	textTmpl, err := template.New("text").Funcs(funcMap).Parse(textTemplate)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse text template: %w", err)
	}

	var textBuf bytes.Buffer
	if err := textTmpl.Execute(&textBuf, data); err != nil {
		return "", "", fmt.Errorf("failed to execute text template: %w", err)
	}

	return htmlBuf.String(), textBuf.String(), nil
}

// ============================================================================
// SMTP OPERATIONS
// ============================================================================

// sendEmail orchestrates the complete email sending process.
//
// Steps:
//   1. Build MIME multi-part message (HTML + text)
//   2. Select appropriate TLS method (STARTTLS or direct)
//   3. Authenticate with SMTP server
//   4. Transmit message
//
// Parameters:
//   - email: Complete email with HTML and text bodies
//
// Returns:
//   - error: SMTP connection or delivery failure
func (s *Service) sendEmail(email DossierEmail) error {
	message := s.buildMIMEMessage(email)

	err := s.sendSMTPWithTLS(s.config.FromEmail, []string{email.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Successfully sent dossier email to %s", email.To)
	return nil
}

// buildMIMEMessage creates a properly formatted MIME multi-part message
// containing both HTML and plain text versions.
//
// MIME Structure:
//   - multipart/alternative: Email clients choose best format
//   - text/plain: First alternative (fallback)
//   - text/html: Second alternative (preferred)
//
// Email Client Behavior:
//   - Modern clients: Display HTML version
//   - Basic clients: Display plain text version
//   - Accessibility tools: May prefer plain text
//
// Parameters:
//   - email: Email data with both HTML and text bodies
//
// Returns:
//   - string: Complete RFC-compliant MIME message
func (s *Service) buildMIMEMessage(email DossierEmail) string {
	boundary := "boundary-dossier-" + fmt.Sprintf("%d", time.Now().Unix())

	message := fmt.Sprintf(`From: %s <%s>
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="%s"

--%s
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 7bit

%s

--%s
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: 7bit

%s

--%s--
`, s.config.FromName, s.config.FromEmail, email.To, email.Subject, 
	boundary, boundary, email.TextBody, boundary, email.HTMLBody, boundary)

	return message
}

// sendSMTPWithTLS sends an email using the appropriate TLS method based on port.
//
// Port-Based Strategy:
//   - 587: STARTTLS (RFC 3207) - Upgrade plain connection
//   - 465: Direct TLS (SMTPS) - TLS from connection start
//   - Other: Attempt STARTTLS as safest fallback
//
// Parameters:
//   - from: Sender email address
//   - to: List of recipient email addresses
//   - msg: Complete RFC-compliant email message
//
// Returns:
//   - error: Connection, authentication, or transmission failure
func (s *Service) sendSMTPWithTLS(from string, to []string, msg []byte) error {
	addr := s.config.SMTPHost + ":" + s.config.SMTPPort
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)

	if s.config.SMTPPort == "587" {
		return s.sendWithSTARTTLS(from, to, msg, auth, addr)
	}
	return s.sendWithDirectTLS(from, to, msg, auth, addr)
}

// ============================================================================
// TLS CONNECTION METHODS
// ============================================================================

// sendWithSTARTTLS sends email using STARTTLS protocol (RFC 3207).
// This is the modern standard for secure SMTP on port 587.
//
// STARTTLS Protocol Flow:
//   1. Connect to SMTP server over plain TCP
//   2. Issue EHLO command to discover capabilities
//   3. Send STARTTLS command to upgrade connection
//   4. Perform TLS handshake
//   5. Authenticate over encrypted connection
//   6. Send email data
//
// Security Features:
//   - Certificate validation (InsecureSkipVerify: false)
//   - Server name verification (SNI)
//   - Prevents downgrade attacks
//
// Parameters:
//   - from: Sender email address
//   - to: List of recipient addresses
//   - msg: Complete email message
//   - auth: SMTP authentication credentials
//   - addr: Server address (host:port)
//
// Returns:
//   - error: Connection, TLS, authentication, or transmission failure
func (s *Service) sendWithSTARTTLS(from string, to []string, msg []byte, auth smtp.Auth, addr string) error {
	// Establish plain TCP connection
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Quit()

	// Configure TLS with certificate validation
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	}

	// Upgrade connection to TLS
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate over encrypted connection
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Send the message
	return s.sendMessage(client, from, to, msg)
}

// sendWithDirectTLS sends email using direct TLS (SMTPS).
// This is the traditional secure SMTP on port 465.
//
// Direct TLS Protocol Flow:
//   1. Establish TLS connection immediately
//   2. Perform TLS handshake before any SMTP commands
//   3. Create SMTP client over TLS connection
//   4. Authenticate over encrypted connection
//   5. Send email data
//
// Differences from STARTTLS:
//   - TLS from connection start (no plaintext phase)
//   - Used on port 465 by convention
//   - Older but still widely supported
//
// Security Features:
//   - Certificate validation
//   - Server name verification
//   - No plaintext exposure
//
// Parameters:
//   - from: Sender email address
//   - to: List of recipient addresses
//   - msg: Complete email message
//   - auth: SMTP authentication credentials
//   - addr: Server address (host:port)
//
// Returns:
//   - error: Connection, TLS, authentication, or transmission failure
func (s *Service) sendWithDirectTLS(from string, to []string, msg []byte, auth smtp.Auth, addr string) error {
	// Configure TLS with certificate validation
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	}

	// Establish TLS connection
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with TLS: %w", err)
	}
	defer conn.Close()

	// Create SMTP client over TLS connection
	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Send the message
	return s.sendMessage(client, from, to, msg)
}

// testWithSTARTTLS tests SMTP connectivity using STARTTLS protocol.
// Used for configuration validation before sending actual emails.
//
// Parameters:
//   - auth: SMTP authentication credentials
//   - addr: Server address (host:port)
//
// Returns:
//   - error: Connection, TLS, or authentication failure
func (s *Service) testWithSTARTTLS(auth smtp.Auth, addr string) error {
	log.Printf("Testing SMTP connection with STARTTLS")
	
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Quit()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	}

	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	log.Printf("SMTP connection test successful with STARTTLS")
	return nil
}

// testWithDirectTLS tests SMTP connectivity using direct TLS.
// Used for configuration validation before sending actual emails.
//
// Parameters:
//   - auth: SMTP authentication credentials
//   - addr: Server address (host:port)
//
// Returns:
//   - error: Connection, TLS, or authentication failure
func (s *Service) testWithDirectTLS(auth smtp.Auth, addr string) error {
	log.Printf("Testing SMTP connection with direct TLS")
	
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with TLS: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	log.Printf("SMTP connection test successful with direct TLS")
	return nil
}

// ============================================================================
// SMTP MESSAGE TRANSMISSION
// ============================================================================

// sendMessage handles the actual SMTP message transmission protocol.
// This implements the core SMTP commands for email delivery.
//
// SMTP Protocol Sequence:
//   1. MAIL FROM: Specify sender
//   2. RCPT TO: Specify each recipient (can be multiple)
//   3. DATA: Send message content
//   4. Quit: Close connection gracefully
//
// Error Handling:
//   - Validates each step before proceeding
//   - Returns detailed error context
//   - Ensures cleanup even on failure
//
// Parameters:
//   - client: Authenticated SMTP client (already connected and encrypted)
//   - from: Sender email address
//   - to: List of recipient addresses
//   - msg: Complete RFC-compliant email message
//
// Returns:
//   - error: SMTP protocol error at any step
func (s *Service) sendMessage(client *smtp.Client, from string, to []string, msg []byte) error {
	// Set sender (MAIL FROM)
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients (RCPT TO for each)
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	if _, err := writer.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// extractDomain extracts a clean domain name from a URL for display purposes.
// This provides user-friendly source attribution in emails.
//
// Processing Steps:
//   1. Remove protocol (http:// or https://)
//   2. Extract hostname (everything before first /)
//   3. Remove www. prefix if present
//
// Examples:
//   - "https://www.example.com/article" → "example.com"
//   - "http://blog.site.org/post/123" → "blog.site.org"
//   - "example.com" → "example.com"
//
// Parameters:
//   - url: Full article URL
//
// Returns:
//   - string: Clean domain name for display
func extractDomain(url string) string {
	// Remove protocol
	if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	}
	
	// Extract hostname
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		domain := parts[0]
		// Remove www. prefix
		if strings.HasPrefix(domain, "www.") {
			domain = domain[4:]
		}
		return domain
	}
	return url
}