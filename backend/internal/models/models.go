package models

type Task struct {
	ID     int    `json:"id"`
	TeamID int    `json:"team_id"`
	Title  string `json:"title"`
	TxHash string `json:"tx_hash"`
	Status string `json:"status"`
}
