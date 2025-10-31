package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/geraldfingburke/dossier/server/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

// Service handles AI operations
type Service struct {
	client *openai.Client
}

// NewService creates a new AI service
func NewService() *Service {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// Allow running without API key for development
		return &Service{client: nil}
	}
	
	client := openai.NewClient(apiKey)
	return &Service{client: client}
}

// SummarizeArticles generates a summary of multiple articles
func (s *Service) SummarizeArticles(ctx context.Context, articles []models.Article) (string, error) {
	if s.client == nil {
		// Return a mock summary if no API key is configured
		return s.mockSummary(articles), nil
	}

	// Prepare the content for summarization
	var content strings.Builder
	content.WriteString("Please provide a concise daily digest summary of the following articles:\n\n")
	
	for i, article := range articles {
		content.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, article.Title))
		if article.Description != "" {
			content.WriteString(fmt.Sprintf("   %s\n", article.Description))
		}
		content.WriteString(fmt.Sprintf("   Link: %s\n\n", article.Link))
		
		// Limit to prevent token overflow
		if i >= 20 {
			content.WriteString("... and more articles\n")
			break
		}
	}

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that creates concise, informative daily digest summaries of RSS feed articles. " +
					"Group related articles together and highlight key themes and important news. " +
					"Keep the summary engaging and easy to read.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content.String(),
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
	})

	if err != nil {
		return "", fmt.Errorf("error calling OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// mockSummary provides a mock summary for development without API key
func (s *Service) mockSummary(articles []models.Article) string {
	var summary strings.Builder
	summary.WriteString("📰 Daily Digest Summary\n\n")
	summary.WriteString(fmt.Sprintf("Today's digest includes %d articles from your RSS feeds.\n\n", len(articles)))
	
	summary.WriteString("**Highlights:**\n")
	for i, article := range articles {
		if i >= 5 {
			summary.WriteString(fmt.Sprintf("\n... and %d more articles\n", len(articles)-5))
			break
		}
		summary.WriteString(fmt.Sprintf("• %s\n", article.Title))
	}
	
	summary.WriteString("\n*Note: This is a mock summary. Configure OPENAI_API_KEY for AI-powered summaries.*")
	
	return summary.String()
}

// SummarizeArticle generates a summary of a single article
func (s *Service) SummarizeArticle(ctx context.Context, article models.Article) (string, error) {
	if s.client == nil {
		return fmt.Sprintf("Summary of: %s\n\n%s", article.Title, article.Description), nil
	}

	content := article.Content
	if content == "" {
		content = article.Description
	}

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that summarizes articles concisely.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Summarize this article:\n\nTitle: %s\n\n%s", article.Title, content),
			},
		},
		MaxTokens:   300,
		Temperature: 0.7,
	})

	if err != nil {
		return "", fmt.Errorf("error calling OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
