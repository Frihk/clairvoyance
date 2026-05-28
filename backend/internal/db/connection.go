package db

import (
    "database/sql"
    _ "github.com/lib/pq"
)

// Connect returns a database handle for the provided connection string
func Connect(connStr string) (*sql.DB, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    return db, db.Ping()
}
