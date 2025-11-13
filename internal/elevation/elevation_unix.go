//go:build !windows
// +build !windows

package elevation

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func isWindowsElevated() bool {
	return false
}

func requestWindowsElevation() error {
	return fmt.Errorf("not on Windows")
}

func requestDarwinElevation() error {
	if runtime.GOOS != "darwin" {
		return requestLinuxElevation()
	}

	// On macOS, try to use osascript for graphical elevation
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// Build command with arguments
	cmdLine := executable
	for _, arg := range os.Args[1:] {
		cmdLine += " " + arg
	}

	script := fmt.Sprintf(`do shell script "%s" with administrator privileges`, cmdLine)
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to elevate privileges: %w", err)
	}

	// If elevation succeeded, exit current process
	os.Exit(0)
	return nil
}

func requestLinuxElevation() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// Try pkexec first (graphical), then sudo
	var cmd *exec.Cmd

	if _, err := exec.LookPath("pkexec"); err == nil {
		cmd = exec.Command("pkexec", append([]string{executable}, os.Args[1:]...)...)
	} else if _, err := exec.LookPath("sudo"); err == nil {
		cmd = exec.Command("sudo", append([]string{executable}, os.Args[1:]...)...)
	} else {
		return fmt.Errorf("no elevation method found (pkexec or sudo required)")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: 0, Gid: 0},
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to elevate privileges: %w", err)
	}

	// If elevation succeeded, exit current process
	os.Exit(0)
	return nil
}
