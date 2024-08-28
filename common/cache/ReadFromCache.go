package cache

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadFromCache(filename string) ([]string, error) {
	if filename == "" {
		filename = "LOCAL_CACHE_DEFAULT.txt"
	}

	// Construct full path relative to the current directory
	filepath := filepath.Join(GetCurrentDirectory(), filename)

	// Try to open the file
	file, err := os.Open(filepath)

	if err != nil {
		// If the file doesn't exist, create an empty one
		if os.IsNotExist(err) {
			_, err := os.Create(filepath)
			if err != nil {
				return nil, fmt.Errorf("failed to create file: %v", err)
			}
			return []string{}, nil // Return empty slice for newly created file
		}
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}
