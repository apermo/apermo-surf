package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const FileName = ".surf-links.yml"
const FileNameDist = ".surf-links.yml.dist"

// Find walks up from startDir looking for .surf-links.yml.
// At each level, .surf-links.yml is checked first; if absent,
// .surf-links.yml.dist is used as fallback. Closest ancestor wins.
func Find(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		primary := filepath.Join(dir, FileName)
		if _, err := os.Stat(primary); err == nil {
			return primary, nil
		}

		dist := filepath.Join(dir, FileNameDist)
		if _, err := os.Stat(dist); err == nil {
			return dist, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no %s found — run surf init to create one", FileName)
}
