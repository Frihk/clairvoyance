package main

import (
	"log"
	"net/http"
	"github.com/ibraah007/clairvoyance/backend/internal/db"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.Close()

	mux := setupRoutes(database)

	log.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", mux)
}
