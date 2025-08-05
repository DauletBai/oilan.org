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

// GeminiClient implements the LLMClient interface for Google Gemini.
type GeminiClient struct {
	client *genai.GenerativeModel
}

// NewGeminiClient creates a new client for interacting with Gemini.
func NewGeminiClient() (services.LLMClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	
	// Initialize the specific model we want to use
	model := client.GenerativeModel("gemini-1.5-pro-latest")

	return &GeminiClient{client: model}, nil
}

// GenerateResponse sends the conversation history to Gemini.
func (c *GeminiClient) GenerateResponse(ctx context.Context, history []domain.Message, systemPrompt string) (string, error) {
	// Start a new chat session with the model
	chat := c.client.StartChat()
	chat.History = make([]*genai.Content, 0, len(history)+1)

	// Set the system prompt for the entire conversation
	c.client.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	// Convert our history to Gemini's format.
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

	// Get the last message from the history, which is the user's current prompt
	if len(history) == 0 {
		return "", fmt.Errorf("cannot generate response for empty history")
	}
	lastMessage := history[len(history)-1].Content

	// Send the message to the model
	resp, err := chat.SendMessage(ctx, genai.Text(lastMessage))
	if err != nil {
		return "", fmt.Errorf("failed to send message to gemini: %w", err)
	}

	// Extract and return the text content from the response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(textPart), nil
		}
	}

	return "", fmt.Errorf("no response content from Gemini")
}