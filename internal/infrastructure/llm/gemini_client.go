// oilan/internal/infrastructure/llm/gemini_client.go
package llm

import (
	"context"
	"fmt"
	"oilan/internal/app/services"
	"oilan/internal/domain"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.GenerativeModel
}

func NewGeminiClient() (services.LLMClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	model := client.GenerativeModel("gemini-1.5-pro-latest")
	return &GeminiClient{client: model}, nil
}

// GenerateResponse now sends the entire context in a single, clean request.
func (c *GeminiClient) GenerateResponse(ctx context.Context, history []domain.Message, systemPrompt string) (string, error) {
	// Set the system prompt for this specific request.
	c.client.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	// Start a new chat session.
	chat := c.client.StartChat()
	
	// Convert our entire history to Gemini's format.
	chat.History = make([]*genai.Content, 0, len(history))
	for _, msg := range history {
		var role string
		if msg.Role == domain.RoleUser {
			role = "user"
		} else {
			role = "model" // Gemini uses "model" for the AI's role
		}
		chat.History = append(chat.History, &genai.Content{
			Role:  role,
			Parts: []genai.Part{genai.Text(msg.Content)},
		})
	}
	
	// The prompt is the entire history. We send an empty message to get a response.
	resp, err := chat.SendMessage(ctx, genai.Text("")) // Send empty message to continue the conversation
	if err != nil {
		return "", fmt.Errorf("failed to send message to gemini: %w", err)
	}

	// Extract and return the text content from the response.
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(textPart), nil
		}
	}

	return "", fmt.Errorf("no text content found in Gemini response")
}