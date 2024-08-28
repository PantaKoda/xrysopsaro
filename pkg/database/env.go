// pkg/database/env.go
package database

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnvFile loads environment variables from the .env file
func LoadEnvFile() error {
	rootDir, err := GetRootDirectory()
	if err != nil {
		return fmt.Errorf("failed to get root directory: %w", err)
	}

	envPath := filepath.Join(rootDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

// GetRootDirectory returns the root directory of the project
func GetRootDirectory() (string, error) {
	return filepath.Abs(filepath.Join("..", ".."))
}
