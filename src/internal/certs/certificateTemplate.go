package certs

import (
	"fmt"
	"path/filepath"
	"time"
)

const (
	expirationDays   = 365 * 3
	templateFileName = "template.cnf"
)

func (m *Manager) createCertificateTemplateIfNotExists(username string) error {
	templatePath := filepath.Join(m.certsDirectory, username, templateFileName)

	if found, err := checkFileExists(templatePath); err != nil {
		return fmt.Errorf("failed to check if template exists: %w", err)
	} else if found {
		return nil
	}

	certificateTemplate := certificateTemplate{
		Organization:     m.organization,
		CertificateOwner: username,
		UserID:           username,
		SerialNumber:     time.Now().UTC().Unix(),
		ExpirationDays:   expirationDays,
	}

	if err := writeToFile(templatePath, certificateTemplate.Render()); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	return nil
}

type certificateTemplate struct {
	Organization     string
	CertificateOwner string
	UserID           string
	SerialNumber     int64
	ExpirationDays   int
}

func (ct *certificateTemplate) Render() string {
	return fmt.Sprintf(`
	# X.509 Certificate options
	# The organization of the subject.
	organization = "%s"

	# The common name of the certificate owner.
	cn = "%s"

	# A user id of the certificate owner.
	uid = "%s"
	serial = %d

	# In how many days, counting from today, this certificate will expire. Use -1 if there is no expiration date.
	expiration_days = %d

	tls_www_client

	signing_key

	encryption_key
	`, ct.Organization, ct.CertificateOwner, ct.UserID, ct.SerialNumber, ct.ExpirationDays)
}
