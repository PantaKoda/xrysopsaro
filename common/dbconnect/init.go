// pkg/database_sqlc/init.go
package dbconnect

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // Import for PostgreSQL driver
)

// InitializeDatabase initializes the database_sqlc connection and returns it
func InitializeDatabase() *pgx.Conn {
	// Load the environment variables from the .env file
	if err := LoadEnvFile(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the database_sqlc URL from the environment
	dbURL := os.Getenv("DATABASE_PROJECT_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_PROJECT_URL not set in the environment")
	}

	// Open the database_sqlc connection
	//db, err := sql.Open("pgx", dbURL)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)

	if err != nil {
		log.Fatalf("Failed to open database_sqlc connection: %v", err)
	}

	log.Println("Connected to database_sqlc")
	return conn
}

// PrintDatabaseHost parses the migrations_goose URL and prints the host for debugging
func PrintDatabaseHost(dbURL string) {
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database_sqlc URL: %v", err)
	}
	fmt.Printf("Connecting to host: %s\n", u.Host)
}

// PingDatabase checks if the database_sqlc connection is alive
func PingDatabase(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("cannot ping the database_sqlc: %w", err)
	}
	return nil
}
