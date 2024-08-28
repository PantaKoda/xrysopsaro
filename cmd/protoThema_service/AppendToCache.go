package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func appendSliceToFile(filename string, data []string) error {
	if filename == "" {
		filename = DEFAULT_FILENAME
	}

	// Construct full path relative to the current directory
	filepath := filepath.Join(GetCurrentDirectory(), filename)

	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(file)

	// Write each string in the slice to a new line
	for _, line := range data {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	// Flush the buffer to ensure all data is written to the file
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	return nil
}
