package main

import (
	"log"
	"net/http"

	"github.com/ibraah007/clairvoyance/backend/internal/db"
)

func main() {
	// Initialize database connection
	database, err := db.Connect()
	if err != nil {
		log.Printf("Warning: Database unreachable: %v", err)
		// Database remains nil; repository.go and router.go handle this gracefully
	}

	// Setup routes
	mux := setupRoutes(database)

	log.Println("Server started on port :8080. If DB failed, app is in Read-Only Demo Mode.")
	
	// Start server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}