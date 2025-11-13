package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/4nkitd/miner/internal/assets"
)

const (
	// Application constants
	AppName    = "miner"
	AppVersion = "1.0.0"

	// Server configuration
	ServerPort   = "88"
	ServerDomain = "miner.local"
	ServerHost   = "127.0.0.1"

	// CLI command names
	CLICommandPHP   = "php"
	CLICommandFPHP  = "fphp"
	CLICommandMiner = "miner"
)

// Config holds application configuration
type Config struct {
	Port       string
	Domain     string
	Host       string
	AppDir     string
	AssetsDir  string
	BinaryPath string
	HostsPath  string
	AutoStart  bool
	TempAssets string // Path to extracted embedded assets (empty if using filesystem)
}

// New creates a new configuration with defaults
func New() (*Config, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Dir(execPath)

	// Try multiple locations for assets directory
	assetsDir := ""
	tempAssets := ""
	possiblePaths := []string{
		filepath.Join(appDir, "assets"),       // Same directory as binary
		filepath.Join(appDir, "..", "assets"), // Parent directory (for dev)
		"assets",                              // Current working directory
	}

	for _, path := range possiblePaths {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			if _, err := os.Stat(filepath.Join(path, "adminer.php")); err == nil {
				assetsDir, _ = filepath.Abs(path)
				break
			}
		}
	}

	// If no filesystem assets found, extract embedded assets
	if assetsDir == "" {
		extracted, err := assets.Extract()
		if err == nil {
			assetsDir = extracted
			tempAssets = extracted
		} else {
			// Fallback to expected path (will fail later if missing)
			assetsDir = filepath.Join(appDir, "assets")
		}
	}

	cfg := &Config{
		Port:       ServerPort,
		Domain:     ServerDomain,
		Host:       ServerHost,
		AppDir:     appDir,
		AssetsDir:  assetsDir,
		BinaryPath: execPath,
		HostsPath:  getHostsPath(),
		AutoStart:  true,
		TempAssets: tempAssets,
	}

	return cfg, nil
}

// URL returns the full server URL
func (c *Config) URL() string {
	if c.Port == "80" {
		return "http://" + c.Domain
	}
	return "http://" + c.Domain + ":" + c.Port
}

// getHostsPath returns the platform-specific hosts file path
func getHostsPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("WINDIR"), "System32", "drivers", "etc", "hosts")
	default:
		return "/etc/hosts"
	}
}
