package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ibraah007/clairvoyance/backend/internal/db"
)

func setupRoutes(database *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// GET tasks endpoint
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := db.GetTasksByTeam(database, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

	// POST verify endpoint
	mux.HandleFunc("POST /verify", func(w http.ResponseWriter, r *http.Request) {
		// In a real scenario, you'd parse JSON body here
		// For the demo, just return a success confirmation
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "verified", "message": "Transaction recorded on-chain"}`))
	})

	return mux
}
