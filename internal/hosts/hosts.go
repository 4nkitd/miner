package hosts

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/txn2/txeh"
)

// Manager handles hosts file modifications
type Manager struct {
	hostsPath string
}

// NewManager creates a new hosts file manager
func NewManager(hostsPath string) *Manager {
	return &Manager{
		hostsPath: hostsPath,
	}
}

// AddEntry adds a domain to IP mapping in the hosts file
func (m *Manager) AddEntry(domain, ip string) error {
	cfg := &txeh.HostsConfig{ReadFilePath: m.hostsPath, WriteFilePath: m.hostsPath}
	hosts, err := txeh.NewHosts(cfg)
	if err != nil {
		return fmt.Errorf("failed to load hosts file: %w", err)
	}
	hosts.AddHost(ip, domain)
	if err := hosts.Save(); err != nil {
		return fmt.Errorf("failed to save hosts file: %w", err)
	}
	return nil
}

// RemoveEntry removes a domain from the hosts file
func (m *Manager) RemoveEntry(domain string) error {
	cfg := &txeh.HostsConfig{ReadFilePath: m.hostsPath, WriteFilePath: m.hostsPath}
	hosts, err := txeh.NewHosts(cfg)
	if err != nil {
		return fmt.Errorf("failed to load hosts file: %w", err)
	}
	hosts.RemoveHost(domain)
	if err := hosts.Save(); err != nil {
		return fmt.Errorf("failed to save hosts file: %w", err)
	}
	return nil
}

// HasEntry checks if a domain exists in the hosts file
func (m *Manager) HasEntry(domain string) (bool, error) {
	f, err := os.Open(m.hostsPath)
	if err != nil {
		return false, fmt.Errorf("failed to open hosts file: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split by whitespace
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		// domains can follow the IP; check any field after the first
		for _, fld := range fields[1:] {
			if fld == domain {
				return true, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading hosts file: %w", err)
	}
	return false, nil
}
