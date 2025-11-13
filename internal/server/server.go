package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Server struct {
	port          string
	domain        string
	assetsDir     string
	running       bool
	frankenphpCmd *exec.Cmd
}

func NewServer(port, domain, assetsDir string) *Server {
	return &Server{
		port:      port,
		domain:    domain,
		assetsDir: assetsDir,
		running:   false,
	}
}

func (s *Server) Start() error {
	if s.running {
		return fmt.Errorf("server is already running")
	}

	frankenphpPath, err := exec.LookPath("frankenphp")
	if err != nil {
		return fmt.Errorf("frankenphp not found. Install it with: curl https://frankenphp.dev/install.sh | sh")
	}

	// Ensure assets directory exists and contains adminer.php
	adminerPath := filepath.Join(s.assetsDir, "adminer.php")
	if _, err := os.Stat(adminerPath); err != nil {
		return fmt.Errorf("adminer.php not found in assets directory: %w", err)
	}

	// Use php-server mode; set working directory to assetsDir so adminer.php is document root
	// Bind explicitly to the requested domain:port using --listen if available, else rely on hosts + default
	// FrankenPHP's php-server listens on 3000 by default; we need port 80
	// Bind only to the port; rely on hosts file + Host header for domain routing
	listen := ":" + s.port

	// Command: frankenphp php-server -r <assetsDir> --listen :port
	args := []string{"php-server", "-r", s.assetsDir, "--listen", listen}

	s.frankenphpCmd = exec.Command(frankenphpPath, args...)
	s.frankenphpCmd.Dir = s.assetsDir
	s.frankenphpCmd.Stdout = os.Stdout
	s.frankenphpCmd.Stderr = os.Stderr

	if err := s.frankenphpCmd.Start(); err != nil {
		return fmt.Errorf("failed to start FrankenPHP php-server: %w", err)
	}

	s.running = true
	fmt.Printf("FrankenPHP php-server started on %s\n", s.URL())

	go func() {
		if err := s.frankenphpCmd.Wait(); err != nil {
			fmt.Printf("FrankenPHP exited: %v\n", err)
			s.running = false
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	if !s.running {
		return fmt.Errorf("server is not running")
	}

	if s.frankenphpCmd != nil && s.frankenphpCmd.Process != nil {
		// Send SIGTERM for graceful shutdown
		if err := s.frankenphpCmd.Process.Signal(syscall.SIGTERM); err != nil {
			// If SIGTERM fails, force kill
			s.frankenphpCmd.Process.Kill()
		}
	}

	s.running = false
	return nil
}

func (s *Server) IsRunning() bool {
	return s.running
}

func (s *Server) URL() string {
	if s.port == "80" {
		return fmt.Sprintf("http://%s", s.domain)
	}
	return fmt.Sprintf("http://%s:%s", s.domain, s.port)
}
