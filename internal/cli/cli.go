package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Manager handles CLI command registration
type Manager struct {
	binaryPath string
}

// NewManager creates a new CLI manager
func NewManager(binaryPath string) *Manager {
	return &Manager{
		binaryPath: binaryPath,
	}
}

// Register registers CLI commands (php, fphp, miner) in system PATH
func (m *Manager) Register() error {
	baseDir := filepath.Dir(m.binaryPath)
	binDir := ""
	if runtime.GOOS == "windows" {
		// On Windows keep wrappers next to the binary by default
		if filepath.Base(baseDir) == "bin" {
			binDir = baseDir
		} else {
			binDir = filepath.Join(baseDir, "bin")
		}
	} else {
		// On Unix/macOS, install to a stable system-wide location
		binDir = "/usr/local/miner/bin"
	}

	// Create bin directory if it doesn't exist
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Create wrapper scripts for each command
	commands := map[string]string{
		"php":   m.getPHPScript(),
		"fphp":  m.getFPHPScript(),
		"miner": m.getMinerScript(),
	}

	for cmdName, scriptContent := range commands {
		if err := m.createWrapper(cmdName, scriptContent, binDir); err != nil {
			return fmt.Errorf("failed to create %s command: %w", cmdName, err)
		}
	}

	// Add bin directory to PATH
	if err := m.addToPath(binDir); err != nil {
		return fmt.Errorf("failed to add to PATH: %w", err)
	}

	return nil
}

// Unregister removes CLI commands from system PATH
func (m *Manager) Unregister() error {
	binDir := ""
	if runtime.GOOS == "windows" {
		binDir = filepath.Join(filepath.Dir(m.binaryPath), "bin")
	} else {
		binDir = "/usr/local/miner/bin"
	}

	// Remove bin directory
	if err := os.RemoveAll(binDir); err != nil {
		return fmt.Errorf("failed to remove bin directory: %w", err)
	}

	// Remove from PATH
	if err := m.removeFromPath(binDir); err != nil {
		return fmt.Errorf("failed to remove from PATH: %w", err)
	}

	return nil
}

func (m *Manager) createWrapper(cmdName, scriptContent, binDir string) error {
	var scriptPath string
	var scriptExt string

	if runtime.GOOS == "windows" {
		scriptExt = ".bat"
	} else {
		scriptExt = ""
	}

	scriptPath = filepath.Join(binDir, cmdName+scriptExt)

	// Write script file
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return err
	}

	return nil
}

func (m *Manager) getPHPScript() string {
	if runtime.GOOS == "windows" {
		return "@echo off\nfrankenphp php-cli %*"
	}
	// Minimal alias-like wrapper to FrankenPHP's php-cli
	return "#!/bin/sh\nexec frankenphp php-cli \"$@\""
}

func (m *Manager) getFPHPScript() string {
	if runtime.GOOS == "windows" {
		return "@echo off\nfrankenphp %*"
	}
	return "#!/bin/sh\nexec frankenphp \"$@\""
}

func (m *Manager) getMinerScript() string {
	if runtime.GOOS == "windows" {
		return "@echo off\nstart http://miner.local"
	}
	return "#!/bin/sh\nopen http://miner.local:88 2>/dev/null || xdg-open http://miner.local:88 2>/dev/null || sensible-browser http://miner.local:88"
}
