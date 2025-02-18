package configuration

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type C struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	ApiKey string `yaml:"api_key"`

	DataDirectory    string `yaml:"data_directory"`
	OcservPasswdPath string `yaml:"ocserv_passwd_path"`

	Organization         string `yaml:"organization"`
	CaCertificatePath    string `yaml:"ca_certificate_path"`
	CaPrivateKeyPath     string `yaml:"ca_private_key_path"`
	ClientPrivateKeyPath string `yaml:"client_private_key_path"`
}

func Initialize() (C, error) {
	configurationPath, err := getPathToConfigurationFileFromFlag()
	if err != nil {
		return C{}, err
	}

	file, err := os.Open(configurationPath)
	if err != nil {
		return C{}, fmt.Errorf("failed to open configuration path: %w", err)
	}
	defer file.Close()

	var configuration C
	if err := yaml.NewDecoder(file).Decode(&configuration); err != nil {
		return C{}, fmt.Errorf("failed to decode configuration: %w", err)
	}

	_, err = os.Stat(configuration.DataDirectory)
	if err != nil {
		return C{}, fmt.Errorf("data directory does not exist: %w", err)
	}

	return configuration, nil
}

func getPathToConfigurationFileFromFlag() (string, error) {
	configurationPath := flag.String("configuration", "", "path to configuration file")
	flag.StringVar(configurationPath, "c", "", "path to configuration file")
	flag.Parse()

	if *configurationPath == "" {
		return "", fmt.Errorf("configuration file path is not provided")
	}

	// check that file exists
	_, err := os.Stat(*configurationPath)
	if err != nil {
		return "", fmt.Errorf("configuration file does not exist: %w", err)
	}

	return *configurationPath, nil
}
