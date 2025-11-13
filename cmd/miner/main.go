package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/4nkitd/miner/internal/assets"
	"github.com/4nkitd/miner/internal/cli"
	"github.com/4nkitd/miner/internal/config"
	"github.com/4nkitd/miner/internal/elevation"
	"github.com/4nkitd/miner/internal/frankenphp"
	"github.com/4nkitd/miner/internal/hosts"
	"github.com/4nkitd/miner/internal/server"
	"github.com/4nkitd/miner/internal/systray"
)

func main() {
	// Check for subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if err := runInstall(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "daemon":
			if err := runDaemon(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "uninstall":
			if err := runUninstall(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "help", "--help", "-h":
			printHelp()
			return
		case "version", "--version", "-v":
			fmt.Printf("Miner v%s\n", config.AppVersion)
			return
		}
	}

	// Default: run the app
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Miner - Standalone Database Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  miner              Start the Miner system tray application")
	fmt.Println("  miner daemon       Run headless server (no tray) in foreground")
	fmt.Println("  miner install      Install and configure Miner (requires admin/root)")
	fmt.Println("  miner uninstall    Remove Miner configuration")
	fmt.Println("  miner help         Show this help message")
	fmt.Println("  miner version      Show version information")
	fmt.Println()
	fmt.Println("After installation, access Adminer at: http://miner.local")
}

func run() error {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if miner is installed (hosts entry exists)
	hostsManager := hosts.NewManager(cfg.HostsPath)
	hasHostEntry, _ := hostsManager.HasEntry(cfg.Domain)

	if !hasHostEntry {
		fmt.Println("Miner is not installed yet. Please run:")
		fmt.Println("  sudo miner install")
		return fmt.Errorf("installation required")
	}

	// Initialize server (no privileges needed)
	srv := server.NewServer(cfg.Port, cfg.Domain, cfg.AssetsDir)

	// Start server
	fmt.Printf("Starting server on %s\n", cfg.URL())
	if err := srv.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Initialize managers for systray
	cliManager := cli.NewManager(cfg.BinaryPath)
	svc, err := config.NewService(cfg)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Create and run system tray app
	app := systray.NewApp(srv, hostsManager, cliManager, svc, cfg)

	fmt.Println("Miner is running. Check system tray for menu.")
	fmt.Printf("Access Adminer at: %s\n", cfg.URL())

	app.Run()

	// Cleanup temp assets if any
	if cfg.TempAssets != "" {
		assets.Cleanup(cfg.TempAssets)
	}

	return nil
}

func runInstall() error {
	fmt.Println("Installing Miner...")
	fmt.Println()

	// Check for elevation early (needed for placing binary in /usr/local/bin)
	if !elevation.IsElevated() {
		fmt.Println("Installation requires administrator/root privileges.")
		return elevation.RequestElevation()
	}

	// Ensure FrankenPHP is installed (auto-download on macOS/Linux)
	if err := frankenphp.EnsureInstalled(); err != nil {
		// If auto-install failed, fall back to manual guidance
		fmt.Printf("Warning: %v\n", err)
		if _, lookupErr := exec.LookPath("frankenphp"); lookupErr != nil {
			fmt.Println("FrankenPHP is required for PHP execution.")
			fmt.Println("Install manually with:")
			fmt.Println("  curl https://frankenphp.dev/install.sh | sh")
			fmt.Println("  sudo mv frankenphp /usr/local/bin/")
			fmt.Print("Continue Miner installation without FrankenPHP? [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				return fmt.Errorf("installation cancelled")
			}
			fmt.Println()
		}
	}

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize managers
	hostsManager := hosts.NewManager(cfg.HostsPath)
	cliManager := cli.NewManager(cfg.BinaryPath)
	svc, err := config.NewService(cfg)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Setup hosts entry
	fmt.Printf("✓ Adding hosts entry: %s -> %s\n", cfg.Domain, cfg.Host)
	if err := hostsManager.AddEntry(cfg.Domain, cfg.Host); err != nil {
		return fmt.Errorf("failed to add hosts entry: %w", err)
	}

	// Register CLI commands
	fmt.Println("✓ Registering CLI commands: php, fphp, miner")
	if err := cliManager.Register(); err != nil {
		return fmt.Errorf("failed to register CLI commands: %w", err)
	}

	// Install auto-start service
	fmt.Println("✓ Installing auto-start service")
	if err := svc.Install(); err != nil {
		fmt.Printf("  Warning: Failed to install auto-start: %v\n", err)
	} else {
		if err := svc.Start(); err != nil {
			fmt.Printf("  Warning: Failed to start service: %v\n", err)
		} else {
			fmt.Println("  Service started in background (persistent daemon).")
		}
	}

	fmt.Println()
	fmt.Println("✓ Installation complete!")
	fmt.Println()
	fmt.Println("To start Miner, run:")
	fmt.Println("  miner")
	fmt.Println()
	fmt.Printf("Then access Adminer at: %s\n", cfg.URL())

	return nil
}

func runUninstall() error {
	fmt.Println("Uninstalling Miner...")
	fmt.Println()

	// Check for elevation
	if !elevation.IsElevated() {
		fmt.Println("Uninstallation requires administrator/root privileges.")
		return elevation.RequestElevation()
	}

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize managers
	hostsManager := hosts.NewManager(cfg.HostsPath)
	cliManager := cli.NewManager(cfg.BinaryPath)
	svc, err := config.NewService(cfg)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Remove hosts entry
	fmt.Printf("✓ Removing hosts entry: %s\n", cfg.Domain)
	if err := hostsManager.RemoveEntry(cfg.Domain); err != nil {
		fmt.Printf("  Warning: Failed to remove hosts entry: %v\n", err)
	}

	// Unregister CLI commands
	fmt.Println("✓ Unregistering CLI commands")
	if err := cliManager.Unregister(); err != nil {
		fmt.Printf("  Warning: Failed to unregister CLI commands: %v\n", err)
	}

	// Uninstall auto-start service
	fmt.Println("✓ Uninstalling auto-start service")
	if err := svc.Uninstall(); err != nil {
		fmt.Printf("  Warning: Failed to uninstall service: %v\n", err)
	}

	fmt.Println()
	fmt.Println("✓ Uninstallation complete!")

	return nil
}

// runDaemon starts the server headlessly (no tray UI) and blocks until interrupted.
func runDaemon() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	hostsManager := hosts.NewManager(cfg.HostsPath)
	if has, _ := hostsManager.HasEntry(cfg.Domain); !has {
		return fmt.Errorf("hosts entry missing; run 'sudo miner install' first")
	}

	srv := server.NewServer(cfg.Port, cfg.Domain, cfg.AssetsDir)
	fmt.Printf("Starting headless server on %s\n", cfg.URL())
	if err := srv.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	fmt.Println("Headless server running. Press Ctrl+C to stop.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Stopping server...")
	_ = srv.Stop()
	if cfg.TempAssets != "" {
		assets.Cleanup(cfg.TempAssets)
	}
	fmt.Println("Server stopped")
	return nil
}
