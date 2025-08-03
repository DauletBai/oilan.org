// oilan/cmd/server/main.go
package main

import (
    //"fmt"
    "log"
    "net/http"
    "github.com/DauletBai/oilan.org/internal/infrastructure/server"
)

func main() {
    // В будущем здесь будет загрузка конфигурации, подключение к БД и т.д.
    
    srv := server.NewServer() // Создаем новый сервер из нашего пакета
    
    log.Println("Starting server on :8080")
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("could not listen on :8080: %v\n", err)
    }
}