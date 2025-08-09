// github.com/DauletBai/oilan.org/internal/infrastructure/llm/openai_client.go
package llm

import (
	"context"
	"github.com/DauletBai/oilan.org/internal/app/services"
	"github.com/DauletBai/oilan.org/internal/domain"
	"os"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient implements the LLMClient interface for the OpenAI API.
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new client for interacting with OpenAI.
func NewOpenAIClient() services.LLMClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return &OpenAIClient{client: client}
}

// GenerateResponse sends the conversation history to OpenAI and gets a response.
func (c *OpenAIClient) GenerateResponse(ctx context.Context, history []domain.Message, systemPrompt string) (string, error) {
	// 1. Convert our internal message format to the format OpenAI requires.
	messages := make([]openai.ChatCompletionMessage, 0, len(history)+1)

	// 2. Add the system prompt first. This sets the AI's personality and goals.
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})

	// 3. Add the rest of the conversation history.
	for _, msg := range history {
		var role string
		if msg.Role == domain.RoleUser {
			role = openai.ChatMessageRoleUser
		} else {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 4. Create the request to the API.
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o, // You can choose other models like gpt-4-turbo or gpt-3.5-turbo
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	// 5. Return the content of the AI's response.
	return resp.Choices[0].Message.Content, nil
}