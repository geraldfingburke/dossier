package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/geraldfingburke/dossier/server/internal/models"
)

// Config holds email configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	Username     string
	Password     string
	FromEmail    string
	FromName     string
}

// Service handles email operations
type Service struct {
	config Config
}

// NewService creates a new email service
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

// getEnvOrDefault gets environment variable or returns default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DossierEmail represents the email content for a dossier
type DossierEmail struct {
	To          string
	Subject     string
	HTMLBody    string
	TextBody    string
	DossierData DossierData
}

// DossierData contains the data for rendering dossier emails
type DossierData struct {
	Title           string
	Summary         string
	Articles        []ArticleData
	GeneratedAt     time.Time
	ArticleCount    int
	Tone            string
	Language        string
	Instructions    string
}

// ArticleData represents an article in the email
type ArticleData struct {
	Title       string
	Description string
	URL         string
	Source      string
	PublishedAt time.Time
}

// SendDossier sends a dossier email
func (s *Service) SendDossier(config *models.DossierConfig, summary string, articles []models.Article) error {
	log.Printf("Preparing to send dossier email: %s to %s", config.Title, config.Email)

	// Prepare article data
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

	// Prepare dossier data
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

	// Generate email content
	htmlBody, textBody, err := s.generateEmailContent(dossierData)
	if err != nil {
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	// Create email
	email := DossierEmail{
		To:          config.Email,
		Subject:     fmt.Sprintf("Dossier - %s", config.Title),
		HTMLBody:    htmlBody,
		TextBody:    textBody,
		DossierData: dossierData,
	}

	// Send email
	return s.sendEmail(email)
}

// generateEmailContent creates HTML and text versions of the email
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

// sendEmail sends the email using SMTP
func (s *Service) sendEmail(email DossierEmail) error {
	// Create message
	message := s.buildMIMEMessage(email)

	// SMTP authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)

	// Send email
	addr := s.config.SMTPHost + ":" + s.config.SMTPPort
	err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{email.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Successfully sent dossier email to %s", email.To)
	return nil
}

// buildMIMEMessage creates a MIME message with both HTML and text parts
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
`, s.config.FromName, s.config.FromEmail, email.To, email.Subject, boundary, boundary, email.TextBody, boundary, email.HTMLBody, boundary)

	return message
}

// extractDomain extracts domain name from URL
func extractDomain(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	}
	
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		domain := parts[0]
		// Remove www. if present
		if strings.HasPrefix(domain, "www.") {
			domain = domain[4:]
		}
		return domain
	}
	return url
}

// TestSMTPConnection tests the SMTP configuration
func (s *Service) TestSMTPConnection() error {
	log.Printf("Testing SMTP connection to %s:%s", s.config.SMTPHost, s.config.SMTPPort)
	
	addr := s.config.SMTPHost + ":" + s.config.SMTPPort
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)
	
	// Try to connect and authenticate
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	log.Printf("SMTP connection test successful")
	return nil
}