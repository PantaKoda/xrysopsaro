// pkg/database/init.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import for PostgreSQL driver
)

// InitializeDatabase initializes the database connection and returns it
func InitializeDatabase() *sql.DB {
	// Load the environment variables from the .env file
	if err := LoadEnvFile(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the database URL from the environment
	dbURL := os.Getenv("DATABASE_PROJECT_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_PROJECT_URL not set in the environment")
	}

	// Open the database connection
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Parse and print the database host (for debugging)
	PrintDatabaseHost(dbURL)

	log.Println("Connected to database")
	return db
}

// PrintDatabaseHost parses the DB URL and prints the host for debugging
func PrintDatabaseHost(dbURL string) {
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}
	fmt.Printf("Connecting to host: %s\n", u.Host)
}

// PingDatabase checks if the database connection is alive
func PingDatabase(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("cannot ping the database: %w", err)
	}
	return nil
}
