// pkg/database/db.go
package database

import (
	"database/sql"
	"log"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitializeDB initializes the database connection using the provided *sql.DB
func InitializeDB(conn *sql.DB) {
	once.Do(func() {
		db = conn
	})
}

// GetDB returns the singleton database connection
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	return db
}
