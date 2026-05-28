package main

import (
    "net/http"
    "github.com/ibraah007/clairvoyance/backend/internal/db"
    "database/sql"
)

func setupRoutes(database *sql.DB) *http.ServeMux {
    mux := http.NewServeMux()

    // Example endpoint to list tasks
    mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
        tasks, err := db.GetTasksByTeam(database, 1) // Hardcoded ID for now
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        // Add JSON encoding logic here
    })

    return mux
}
