// oilan/cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"oilan/internal/app/services"
	"oilan/internal/infrastructure/handlers"
	//"oilan/internal/auth"
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
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := "http://localhost:8080/auth/google/callback"

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, callbackURL, "email", "profile"),
	)

	// 1. Database Connection
	db, err := postgres.NewConnection()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection successful")

	// 2. Repositories
	userRepo := postgres.NewUserRepository(db)
	dialogRepo := postgres.NewDialogRepository(db)

	// 3. LLM Client
	mockLLM := llm.NewMockLLMClient()

	// 4. Services
	chatService, err := services.NewChatService(dialogRepo, mockLLM)
	if err != nil {
		log.Fatalf("failed to create chat service: %v", err)
	}

	// 5. Template Parsing ---
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

	// 6. Handlers
	apiHandlers := handlers.NewAPIHandlers(chatService, userRepo)
	pageHandlers := &handlers.PageHandlers{WelcomeTemplate: welcomeTpl}

	// 7. Server
	srv := server.NewServer(apiHandlers, pageHandlers)

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on :8080: %v\n", err)
	}
}
