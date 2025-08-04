// oilan/internal/infrastructure/llm/mock_client.go
package llm

import (
	"context"
	"oilan/internal/app/services"
	"oilan/internal/domain"
	"time"
)

// MockLLMClient is a dummy implementation of the LLMClient for testing.
type MockLLMClient struct{}

// NewMockLLMClient creates a new mock LLM client.
func NewMockLLMClient() services.LLMClient {
	return &MockLLMClient{}
}

// GenerateResponse simulates a response from an LLM.
func (c *MockLLMClient) GenerateResponse(ctx context.Context, history []domain.Message, prompt string) (string, error) {
	// Simulate a network delay
	time.Sleep(1 * time.Second)
	
	// Return a fixed, pre-programmed response.
	return "This is a mock response from the AI. The real LLM is not connected yet.", nil
}