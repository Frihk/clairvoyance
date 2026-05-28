package db

import (
	"database/sql"
	"github.com/ibraah007/clairvoyance/backend/internal/models"
)

// GetTasksByTeam fetches tasks for your dashboard
func GetTasksByTeam(db *sql.DB, teamID int) ([]models.Task, error) {
	rows, err := db.Query("SELECT id, team_id, title, tx_hash, status FROM tasks WHERE team_id = $1", teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.TeamID, &t.Title, &t.TxHash, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
