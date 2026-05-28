package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	godotenv.Load(".env")
	godotenv.Load("../../.env")

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("WARNING: Could not reach database: %v\n", err)
		fmt.Println("Proceeding with mock data for demo...")
		return db, nil
	}

	return db, nil
}
