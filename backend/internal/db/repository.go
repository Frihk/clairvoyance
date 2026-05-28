package db

import (
	"database/sql"
	"github.com/ibraah007/clairvoyance/backend/internal/models"
)

func GetTasksByTeam(db *sql.DB, teamID int) ([]models.Task, error) {
	// Safety check: if DB is nil, return empty list instead of panicking
	if db == nil {
		return []models.Task{}, nil
	}

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

func IssueCredential(db *sql.DB, userID int, hash string) error {
	// Safety check
	if db == nil {
		return nil
	}
	_, err := db.Exec("INSERT INTO credentials (user_id, hash) VALUES ($1, $2)", userID, hash)
	return err
}

func UpdateTaskStatus(db *sql.DB, taskID int, txHash string, status string) error {
	// Safety check
	if db == nil {
		return nil
	}
	_, err := db.Exec("UPDATE tasks SET tx_hash = $1, status = $2 WHERE id = $3", txHash, status, taskID)
	return err
}