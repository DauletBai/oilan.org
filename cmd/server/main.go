// oilan/cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"oilan/internal/app/services"
	// "oilan/internal/auth"
	"oilan/internal/infrastructure/handlers"
	"oilan/internal/infrastructure/llm"
	"oilan/internal/infrastructure/repository/postgres"
	"oilan/internal/infrastructure/server"
	"oilan/internal/view"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	// Goth Configuration
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := "http://localhost:8080/auth/google/callback"
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	gothic.Store = store
	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, callbackURL, "email", "profile"),
	)

	// Database Connection
	db, err := postgres.NewConnection()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection successful")

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	dialogRepo := postgres.NewDialogRepository(db)

	// LLM Client
	// We are now switching from the mock client to the real OpenAI client.
	// mockLLM := llm.NewMockLLMClient() 
	// llmClient := llm.NewOpenAIClient()
	llmClient, err := llm.NewGeminiClient() 
	if err != nil {
		log.Fatalf("failed to create gemini client: %v", err)
	}

	// Services
	chatService, err := services.NewChatService(dialogRepo, llmClient) // Pass the real client
	if err != nil {
		log.Fatalf("failed to create chat service: %v", err)
	}

	// --- Template Parsing ---
	welcomeTpl, err := view.NewTemplate(
		"web/templates/base.html",
		"web/templates/parts/head.html",
		"web/templates/parts/header.html",
		"web/templates/parts/footer.html",
		"web/templates/pages/welcome.html",
	)
	if err != nil {
		log.Fatalf("could not parse welcome template: %v", err)
	}

	// --- ADDED MISSING BLOCK HERE ---
	chatTpl, err := view.NewTemplate(
		"web/templates/base.html",
		"web/templates/parts/head.html",
		"web/templates/parts/header.html",
		"web/templates/parts/footer.html",
		"web/templates/pages/chat.html",
	)
	if err != nil {
		log.Fatalf("could not parse chat template: %v", err)
	}

	// --- DashboardTpl Parsing ---
	dashboardTpl, err := view.NewTemplate(
		"web/templates/admin/base.html",
		"web/templates/admin/parts/admin_head.html",
		"web/templates/admin/parts/admin_header.html",
		"web/templates/admin/pages/dashboard.html",
	)
	if err != nil {
		log.Fatalf("could not parse dashboard template: %v", err)
	}

	// --- Handlers ---
	apiHandlers := handlers.NewAPIHandlers(chatService, userRepo)
	pageHandlers := &handlers.PageHandlers{
		WelcomeTemplate: welcomeTpl,
		ChatTemplate:    chatTpl, 
	}
	adminHandlers := &handlers.AdminHandlers{DashboardTemplate: dashboardTpl}

	// --- Server ---
	srv := server.NewServer(apiHandlers, pageHandlers, adminHandlers, userRepo)

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on :8080: %v\n", err)
	}
}