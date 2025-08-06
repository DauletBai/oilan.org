// oilan/cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"oilan/internal/app/services"
	//"oilan/internal/auth"
	"oilan/internal/domain/repository"
	"oilan/internal/infrastructure/handlers"
	"oilan/internal/infrastructure/llm"
	//"oilan/internal/infrastructure/middleware"
	"oilan/internal/infrastructure/repository/postgres"
	"oilan/internal/infrastructure/server"
	"oilan/internal/view"
	"os"

	//"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// This function runs on startup to ensure the configured admin user exists and has the correct role.
func bootstrapAdmin(userRepo repository.UserRepository) {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		log.Println("ADMIN_EMAIL not set, skipping admin bootstrap.")
		return
	}

	ctx := context.Background()
	user, err := userRepo.FindByEmail(ctx, adminEmail)
	if err != nil {
		log.Printf("Error finding admin user by email: %v", err)
		return
	}

	if user != nil {
		if user.Role != "admin" {
			log.Printf("Promoting user %s to admin...", adminEmail)
			user.Role = "admin"
			if err := userRepo.Update(ctx, user); err != nil {
				log.Printf("Failed to promote user to admin: %v", err)
			} else {
				log.Printf("User %s successfully promoted to admin.", adminEmail)
			}
		}
	} else {
		log.Printf("Admin user %s not found in database. Please log in once to create the user.", adminEmail)
	}
}

func main() {
	// --- Goth Configuration ---
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := "http://localhost:8080/auth/google/callback"
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	gothic.Store = store
	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, callbackURL, "email", "profile"),
	)

	// --- Database Connection ---
	db, err := postgres.NewConnection()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection successful")

	// --- Repositories ---
	userRepo := postgres.NewUserRepository(db)
	dialogRepo := postgres.NewDialogRepository(db)

	bootstrapAdmin(userRepo)

	// --- Call the bootstrap function right after creating the repositories ---
	bootstrapAdmin(userRepo)

	// --- LLM Client ---
	llmClient, err := llm.NewGeminiClient()
	if err != nil {
		log.Fatalf("failed to create gemini client: %v", err)
	}

	// --- Services ---
	chatService, err := services.NewChatService(dialogRepo, llmClient)
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
	if err != nil { log.Fatalf("could not parse welcome template: %v", err) }

	chatTpl, err := view.NewTemplate(
		"web/templates/base.html",
		"web/templates/parts/head.html",
		"web/templates/parts/header.html",
		"web/templates/parts/footer.html",
		"web/templates/pages/chat.html",
	)
	if err != nil { log.Fatalf("could not parse chat template: %v", err) }

	dashboardTpl, err := view.NewTemplate(
		"web/templates/admin/base.html",
		"web/templates/admin/parts/admin_head.html",
		"web/templates/admin/parts/admin_header.html",
		"web/templates/admin/pages/dashboard.html",
	)
	if err != nil { log.Fatalf("could not parse dashboard template: %v", err) }
	
	usersTpl, err := view.NewTemplate(
		"web/templates/admin/base.html",
		"web/templates/admin/parts/admin_head.html",
		"web/templates/admin/parts/admin_header.html",
		"web/templates/admin/pages/users.html",
	)
	if err != nil { log.Fatalf("could not parse users template: %v", err) }

	// --- Handlers ---
	apiHandlers := handlers.NewAPIHandlers(chatService, userRepo)
	pageHandlers := &handlers.PageHandlers{
		WelcomeTemplate: welcomeTpl,
		ChatTemplate:    chatTpl,
	}
	
	// --- THIS IS THE FIX ---
	adminHandlers := &handlers.AdminHandlers{
		DashboardTemplate: dashboardTpl,
		UsersTemplate:     usersTpl,
		UserRepo:          userRepo, // We must pass the userRepo here
	}

	// --- Server ---
	srv := server.NewServer(apiHandlers, pageHandlers, adminHandlers, userRepo)

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on :8080: %v\n", err)
	}
}