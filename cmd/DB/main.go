package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib" // Import for PostgreSQL driver
	"github.com/joho/godotenv"
)

func main() {
	// Get the root directory of your project
	rootDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		log.Fatalf("Failed to get root directory: %v", err)
	}

	log.Println("RootDir :", rootDir)
	// Load the .env file
	envPath := filepath.Join(rootDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the database URL from the environment
	dbURL := os.Getenv("DATABASE_PROJECT_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_PROJECT_URL not set in the environment")
	}

	// Parse the database URL
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}
	fmt.Printf("Connecting to host: %s\n", u.Host)

	// Open the database connection
	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("Error closing the database connection: %v", cerr)
		}
	}()
	log.Println("Connected to database")

	// Ping the database to ensure it's reachable
	if err := conn.Ping(); err != nil {
		log.Fatalf("Cannot ping the database: %v", err)
	}
	log.Println("Pinged database successfully")
}
