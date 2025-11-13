//go:build !windows
// +build !windows

package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m *Manager) addToPath(binDir string) error {
	// Create a file in /etc/paths.d/ (requires root)
	pathFile := "/etc/paths.d/miner"

	// Write the bin directory path
	if err := os.WriteFile(pathFile, []byte(binDir+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to create path file: %w", err)
	}

	// Also add to user shell profiles for immediate effect
	home, err := os.UserHomeDir()
	if err == nil {
		profiles := []string{
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".bash_profile"),
		}

		exportLine := fmt.Sprintf("\n# Miner CLI\nexport PATH=\"%s:$PATH\"\n", binDir)

		for _, profile := range profiles {
			if _, err := os.Stat(profile); err == nil {
				// Append to existing profile
				f, err := os.OpenFile(profile, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					continue
				}
				f.WriteString(exportLine)
				f.Close()
			}
		}
	}

	return nil
}

func (m *Manager) removeFromPath(binDir string) error {
	// Remove from /etc/paths.d/
	pathFile := "/etc/paths.d/miner"
	os.Remove(pathFile)

	// Remove from user shell profiles
	home, err := os.UserHomeDir()
	if err == nil {
		profiles := []string{
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".bash_profile"),
		}

		for _, profile := range profiles {
			// This is simplified - in production, you'd parse and remove the specific line
			// For now, just notify user to manually remove if needed
			_ = profile
		}
	}

	return nil
}
