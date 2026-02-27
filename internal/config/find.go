package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const FileName = ".surf-links.yml"

// Find walks up from startDir looking for .surf-links.yml.
// Returns the absolute path to the file, or an error if not found.
func Find(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(dir, FileName)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no %s found — run surf init to create one", FileName)
}
