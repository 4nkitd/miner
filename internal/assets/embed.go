package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed files/*
var embeddedFS embed.FS

// Extract writes embedded assets to a temporary directory and returns the path.
func Extract() (string, error) {
	tempDir, err := os.MkdirTemp("", "miner-assets-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Walk the embedded filesystem and extract all files
	err = fs.WalkDir(embeddedFS, "files", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == "files" {
			return nil
		}

		// Get relative path (remove "files/" prefix)
		relPath, err := filepath.Rel("files", path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(tempDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read embedded file
		data, err := embeddedFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write to temp directory
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		return nil
	})

	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to extract assets: %w", err)
	}

	return tempDir, nil
}

// Cleanup removes the temporary assets directory.
func Cleanup(path string) error {
	if path == "" {
		return nil
	}
	return os.RemoveAll(path)
}
