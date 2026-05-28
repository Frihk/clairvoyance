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
		// Safety check: if DB is nil, inform the user about Demo Mode
		if database == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "offline", "message": "Demo mode: Database unreachable"}`))
			return
		}

		// 1. Parsing would happen here, keeping it static for now as requested
		taskID := 1
		txHash := "0xRealTransactionHash123"

		// 2. Call the real database update
		err := db.UpdateTaskStatus(database, taskID, txHash, "Completed")
		if err != nil {
			http.Error(w, "Failed to update task", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "success", "message": "Task updated in database"}`))
	})

	return mux
}