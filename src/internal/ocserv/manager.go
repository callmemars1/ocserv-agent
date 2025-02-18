package ocserv

import (
	"fmt"
	"os/exec"
	"strings"
)

type Manager struct {
	passwdPath string
}

func NewManager(passwdPath string) *Manager {
	return &Manager{
		passwdPath: passwdPath,
	}
}

func (m *Manager) AddUser(username, password string) error {
	cmd := exec.Command("ocpasswd", "-c", m.passwdPath, username)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	go func() {
		defer stdin.Close()
		_, _ = stdin.Write([]byte(password + "\n"))
		_, _ = stdin.Write([]byte(password + "\n"))
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

func (m *Manager) DisconnectUser(username string) error {
	cmd := exec.Command("occtl", "disconnect", "user", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "could not disconnect user") {
			return nil // ignore error when user is not connected
		}
		return fmt.Errorf("failed to disconnect user: %w", err)
	}
	return nil
}
