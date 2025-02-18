package certs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/callmemars1/setka/src/bot/src/internal/users"
)

const (
	clientCertName   = "cert.pem"
	metadataFileName = "metadata.json"
)

type Configuration struct {
	Organization string

	DataDirectory string

	CaCertificatePath    string
	CaPrivateKeyPath     string
	ClientPrivateKeyPath string
}

type Manager struct {
	organization   string
	certsDirectory string
	caCertPath     string
	caPKPath       string
	clientPKPath   string
}

func NewManager(configuration Configuration) *Manager {
	certsPath := filepath.Join(configuration.DataDirectory, "certificates")

	if err := os.MkdirAll(certsPath, 0755); err != nil {
		panic(fmt.Errorf("failed to create certificates directory: %w", err))
	}

	return &Manager{
		organization:   configuration.Organization,
		certsDirectory: certsPath,
		caCertPath:     configuration.CaCertificatePath,
		caPKPath:       configuration.CaPrivateKeyPath,
		clientPKPath:   configuration.ClientPrivateKeyPath,
	}
}

func (m *Manager) IssueCertificateForUser(user *users.User) ([]byte, error) {
	userDirPath := filepath.Join(m.certsDirectory, user.Username)

	if err := os.MkdirAll(userDirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for user certificates: %w", err)
	}

	if err := m.createCertificateTemplateIfNotExists(user.Username); err != nil {
		return nil, err
	}

	if err := m.createClientCertificateIfNotExists(user.Username); err != nil {
		return nil, err
	}

	if err := m.packClientCertificateToP12(user.Username, user.Password); err != nil {
		return nil, err
	}

	return m.ReadClientCertificateP12(user.Username)
}

func (m *Manager) createClientCertificateIfNotExists(username string) error {
	templatePath := filepath.Join(m.certsDirectory, username, templateFileName)
	if found, err := checkFileExists(templatePath); err != nil {
		return fmt.Errorf("failed to check if client certificate template exists: %w", err)
	} else if !found {
		return fmt.Errorf("client certificate template does not exist")
	}

	clientCertificateDirPath := filepath.Join(m.certsDirectory, username)
	if err := os.MkdirAll(clientCertificateDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory for client certificate: %w", err)
	}

	clientCertificatePath := filepath.Join(clientCertificateDirPath, clientCertName)
	if found, err := checkFileExists(clientCertificatePath); err != nil {
		return fmt.Errorf("failed to check if client certificate exists: %w", err)
	} else if found {
		return nil
	}

	if _, err := exec.Command(
		"certtool",
		"--generate-certificate",
		"--load-privkey", m.clientPKPath,
		"--load-ca-certificate", m.caCertPath,
		"--load-ca-privkey", m.caPKPath,
		"--template", templatePath,
		"--outfile", clientCertificatePath,
	).Output(); err != nil {
		return fmt.Errorf("failed to generate client certificate: %w", err)
	}

	return nil
}

func (m *Manager) packClientCertificateToP12(username, password string) error {
	clientCertificatePath := filepath.Join(m.certsDirectory, username, clientCertName)
	if found, err := checkFileExists(clientCertificatePath); err != nil {
		return fmt.Errorf("failed to check if client certificate exists: %w", err)
	} else if !found {
		return fmt.Errorf("client certificate does not exist")
	}

	p12DirPath := filepath.Join(m.certsDirectory, username)
	if err := os.MkdirAll(p12DirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory for p12 certificate: %w", err)
	}

	p12Path := filepath.Join(p12DirPath, "cert.p12")
	if found, err := checkFileExists(p12Path); err != nil {
		return fmt.Errorf("failed to check if p12 certificate exists: %w", err)
	} else if found {
		return nil
	}

	_, err := exec.Command(
		"certtool",
		"--to-p12",
		"--pkcs-cipher", "3des-pkcs12",
		"--hash", "SHA1",
		"--p12-name", "setka",
		"--load-certificate", clientCertificatePath,
		"--load-privkey", m.clientPKPath,
		"--outfile", p12Path,
		"--outder",
		"--password", password,
	).Output()

	return err
}

func (m *Manager) ReadClientCertificateP12(username string) ([]byte, error) {
	p12Path := filepath.Join(m.certsDirectory, username, "cert.p12")
	if found, err := checkFileExists(p12Path); err != nil {
		return nil, fmt.Errorf("failed to check if p12 certificate exists: %w", err)
	} else if !found {
		return nil, fmt.Errorf("p12 certificate does not exist")
	}

	return os.ReadFile(p12Path)
}
