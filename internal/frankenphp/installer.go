package frankenphp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// EnsureInstalled checks whether the frankenphp binary is available in PATH.
// If missing (on macOS/Linux), it attempts to download and install it using the
// official install script. On Windows, users must install via WSL manually.
func EnsureInstalled() error {
	if _, err := exec.LookPath("frankenphp"); err == nil {
		return nil
	}
	if runtime.GOOS == "windows" {
		return errors.New("FrankenPHP auto-install unsupported on Windows; use WSL: curl https://frankenphp.dev/install.sh | sh")
	}
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	fmt.Println("FrankenPHP not found. Attempting automatic install...")
	tmpDir, scriptPath, err := downloadScript()
	if err != nil {
		return err
	}
	if err = runScript(tmpDir, scriptPath); err != nil {
		return err
	}
	if err = placeBinary(tmpDir); err != nil {
		return err
	}
	fmt.Println("FrankenPHP installed successfully.")
	return nil
}

func downloadScript() (string, string, error) {
	resp, err := http.Get("https://frankenphp.dev/install.sh")
	if err != nil {
		return "", "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	tmpDir, err := os.MkdirTemp("", "frankenphp-install-*")
	if err != nil {
		return "", "", fmt.Errorf("temp dir error: %w", err)
	}
	scriptPath := filepath.Join(tmpDir, "install.sh")
	f, err := os.Create(scriptPath)
	if err != nil {
		return "", "", fmt.Errorf("script create error: %w", err)
	}
	if _, err = io.Copy(f, resp.Body); err != nil {
		f.Close()
		return "", "", fmt.Errorf("write script error: %w", err)
	}
	f.Close()
	if err = os.Chmod(scriptPath, 0o755); err != nil {
		return "", "", fmt.Errorf("chmod script error: %w", err)
	}
	return tmpDir, scriptPath, nil
}

func runScript(tmpDir, scriptPath string) error {
	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install script failed: %w", err)
	}
	return nil
}

func placeBinary(tmpDir string) error {
	localBinary := filepath.Join(tmpDir, "frankenphp")
	if _, err := os.Stat(localBinary); err != nil {
		return fmt.Errorf("binary missing after install: %w", err)
	}
	targetDir := "/usr/local/bin"
	if _, err := os.Stat(targetDir); err != nil {
		return fmt.Errorf("target dir missing: %w", err)
	}
	targetPath := filepath.Join(targetDir, "frankenphp")
	if _, err := os.Stat(targetPath); err == nil {
		_ = os.Remove(targetPath)
	}
	if err := os.Rename(localBinary, targetPath); err != nil {
		in, ierr := os.Open(localBinary)
		if ierr != nil {
			return fmt.Errorf("open binary error: %w", ierr)
		}
		defer in.Close()
		out, oerr := os.Create(targetPath)
		if oerr != nil {
			return fmt.Errorf("create target error: %w", oerr)
		}
		if _, cerr := io.Copy(out, in); cerr != nil {
			out.Close()
			return fmt.Errorf("copy binary error: %w", cerr)
		}
		out.Close()
	}
	if err := os.Chmod(targetPath, 0o755); err != nil {
		return fmt.Errorf("chmod target error: %w", err)
	}
	if _, err := exec.LookPath("frankenphp"); err != nil {
		return fmt.Errorf("frankenphp not in PATH after install: %w", err)
	}
	return nil
}
