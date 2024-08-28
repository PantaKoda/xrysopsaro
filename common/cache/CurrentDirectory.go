package cache

import (
	"log"
	"os"
	"path/filepath"
)

// Get the directory where the go executable is located
// Not good if running go run
// it takes in to account where the binary is
func GetCurrentDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)

	return exPath
}
