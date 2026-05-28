package main

import (
	"database/sql"
	"net/http"
	"encoding/json"
	"github.com/ibraah007/clairvoyance/backend/internal/db"
)

func setupRoutes(database *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := db.GetTasksByTeam(database, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

	return mux
}
