// oilan/internal/app/services/chat_service.go
package services

import (
	"context"
	"errors"
	"fmt"
	"oilan/internal/domain"
	"oilan/internal/domain/repository"
	"os"
	"time"
)

// LLMClient defines the interface for an external Large Language Model.
type LLMClient interface {
	GenerateResponse(ctx context.Context, history []domain.Message, systemPrompt string) (string, error)
}

// ChatService provides methods for chat-related operations.
type ChatService struct {
	dialogRepo repository.DialogRepository
	llmClient  LLMClient
	systemPrompt string
}

// NewChatService creates a new ChatService.
func NewChatService(dialogRepo repository.DialogRepository, llmClient LLMClient) (*ChatService, error) {
	// Read the system prompt from the file system upon initialization.
	promptBytes, err := os.ReadFile("configs/prompt_therapist.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read system prompt: %w", err)
	}

	return &ChatService{
		dialogRepo: dialogRepo,
		llmClient:  llmClient,
		systemPrompt: string(promptBytes),
	}, nil
}

// StartNewDialog creates a new dialog for a user.
func (s *ChatService) StartNewDialog(ctx context.Context, userID int64, title string) (*domain.Dialog, error) {
	if title == "" {
		title = "New Chat"
	}

	dialog := &domain.Dialog{
		UserID: userID,
		Title:  title,
	}

	err := s.dialogRepo.Save(ctx, dialog)
	if err != nil {
		return nil, fmt.Errorf("could not save new dialog: %w", err)
	}

	return dialog, nil
}

// PostMessage handles posting a new message to a dialog and getting a response from the AI.
func (s *ChatService) PostMessage(ctx context.Context, dialogID int64, userID int64, content string) (*domain.Message, error) {
	// 1. Verify that the user owns the dialog (security check).
	dialog, err := s.dialogRepo.FindByID(ctx, dialogID)
	if err != nil {
		return nil, fmt.Errorf("could not find dialog: %w", err)
	}
	if dialog == nil {
		return nil, errors.New("dialog not found")
	}
	if dialog.UserID != userID {
		return nil, errors.New("user does not own this dialog") // Security error
	}

	// 2. Save the user's message to the database.
	userMessage := &domain.Message{
		DialogID: dialogID,
		Role:     domain.RoleUser,
		Content:  content,
		CreatedAt: time.Now(),
	}
	if err := s.dialogRepo.AddMessage(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("could not save user message: %w", err)
	}

	// 3. Get the full, updated conversation history.
	// We add the new user message to the history we already loaded.
	dialog.Messages = append(dialog.Messages, *userMessage)
	
	// 4. Send the history and the system prompt to the LLM to get a response.
	aiContent, err := s.llmClient.GenerateResponse(ctx, dialog.Messages, s.systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("llm client failed to generate response: %w", err)
	}

	// 5. Save the AI's message to the database.
	aiMessage := &domain.Message{
		DialogID: dialogID,
		Role:     domain.RoleAI,
		Content:  aiContent,
		CreatedAt: time.Now(),
	}
	if err := s.dialogRepo.AddMessage(ctx, aiMessage); err != nil {
		return nil, fmt.Errorf("could not save ai message: %w", err)
	}

	// 6. Return the AI's message to the user.
	return aiMessage, nil
}