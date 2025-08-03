// oilan/internal/infrastructure/handlers/routes.go
package handlers

import "net/http"

// RegisterRoutes регистрирует все маршруты приложения.
func RegisterRoutes(mux *http.ServeMux) {
    // Маршрут для главной страницы
    mux.HandleFunc("/", HomeHandler)

    // В будущем здесь будут маршруты для аутентификации, чата и т.д.
    // mux.HandleFunc("/auth/google/login", GoogleLoginHandler)
    // mux.HandleFunc("/api/chat", ChatAPIHandler)
}