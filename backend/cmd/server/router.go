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
		// 1. Parsing would happen here, keeping it static for now as requested
		taskID := 1 
		txHash := "0xRealTransactionHash123"

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