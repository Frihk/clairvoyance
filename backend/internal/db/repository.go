package db

import (
	"database/sql"
	"github.com/ibraah007/clairvoyance/backend/internal/models"
)

// GetTasksByTeam returns mock data so the API works regardless of network
func GetTasksByTeam(db *sql.DB, teamID int) ([]models.Task, error) {
	// MOCK DATA: This allows you to develop the frontend and flow
	mockTasks := []models.Task{
		{ID: 1, TeamID: teamID, Title: "On-chain verification setup", TxHash: "0xabc123", Status: "Pending"},
		{ID: 2, TeamID: teamID, Title: "Database audit log sync", TxHash: "0xdef456", Status: "Completed"},
	}
	return mockTasks, nil
}