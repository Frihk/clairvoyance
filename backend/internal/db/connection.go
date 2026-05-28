package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Connect loads the DB_URL from .env and connects to the database
func Connect() (*sql.DB, error) {
	// 1. Try to load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Could not load .env file, checking system variables instead")
	}

	// 2. Fetch the URL from the environment
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	// 3. Connect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// 4. Verify the connection
	return db, db.Ping()
}
