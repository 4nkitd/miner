package elevation

import (
	"fmt"
	"os"
	"runtime"
)

// IsElevated checks if the current process has admin/root privileges
func IsElevated() bool {
	switch runtime.GOOS {
	case "windows":
		return isWindowsElevated()
	default:
		return os.Geteuid() == 0
	}
}

// RequestElevation attempts to restart the application with elevated privileges
func RequestElevation() error {
	if IsElevated() {
		return nil
	}

	switch runtime.GOOS {
	case "windows":
		return requestWindowsElevation()
	case "darwin":
		return requestDarwinElevation()
	default:
		return requestLinuxElevation()
	}
}

// CheckAndElevate checks privileges and elevates if necessary
func CheckAndElevate() error {
	if !IsElevated() {
		fmt.Println("Miner requires administrator privileges to configure hosts file and PATH.")
		fmt.Println("Requesting elevation...")
		return RequestElevation()
	}
	return nil
}
