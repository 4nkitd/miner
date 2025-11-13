//go:build windows
// +build windows

package cli

// Windows stubs: PATH modification handled differently by installer; wrappers placed in bin directory already on PATH via user configuration.
func (m *Manager) addToPath(binDir string) error      { return nil }
func (m *Manager) removeFromPath(binDir string) error { return nil }
