package main

import (
	"log"
	"net/http"
	"github.com/ibraah007/clairvoyance/backend/internal/db"
)

func main() {
    database, err := db.Connect()
    if err != nil {
        // Only log, don't use log.Fatalf (which stops the server)
        log.Printf("Database connection status: %v", err)
    }

    mux := setupRoutes(database)
    log.Println("Server successfully started on port :8080")
    http.ListenAndServe(":8080", mux)
}