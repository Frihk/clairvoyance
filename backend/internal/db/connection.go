package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Connect loads the environment and establishes a secure DB connection
func Connect() (*sql.DB, error) {
	// 1. Load .env file
	godotenv.Load(".env")
	godotenv.Load("../../.env")

	// 2. Fetch the DB_URL
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	// 3. Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// 4. Ping the database
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping database: %v", err)
	}

	return db, nil
}