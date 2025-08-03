// oilan/internal/infrastructure/handlers/pages.go
package handlers

import (
	"fmt"
	"net/http"
)

// HomeHandler обрабатывает запросы к главной странице.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Oilan Project!")
}