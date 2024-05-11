package ssh

import (
	"fmt"
	"os"
	"os/user"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// Default values for SSH configuration.
const (
	defaultPort           int    = 22
	defaultKnownHostsPath string = ".ssh/known_hosts"
)

// Config represents the configuration details for establishing an SSH connection.
type Config struct {
	Host           string
	Port           int
	Username       string
	Password       string
	PrivateKey     string
	PrivateKeyPath string
	KnownHostsPath string
	Insecure       bool
}

// validate validates the SSH configuration parameters.
func (config *Config) validate() error {
	if (config.Host == "" || config.Username == "") || (config.Password == "" && config.PrivateKey == "" && config.PrivateKeyPath == "") {
		return fmt.Errorf("ssh: Config parameter 'Host', 'Username' and one of 'Password', 'PrivateKey', 'PrivateKeyPath' must be set")
	}

	return nil
}

// defaults sets the default values for the SSH configuration.
func (config *Config) defaults() error {
	if config.Port == 0 {
		config.Port = defaultPort
	}

	if config.KnownHostsPath == "" {
		// Get the current user from the system
		user, err := user.Current()
		if err != nil {
			return err
		}

		config.KnownHostsPath = fmt.Sprintf("%s/%s", user.HomeDir, defaultKnownHostsPath)
	}

	return nil
}

// knownHostCallback generates a host key callback based on the SSH configuration.
func (config *Config) knownHostCallback() (ssh.HostKeyCallback, error) {

	// Ignore host key
	if config.Insecure {
		return ssh.InsecureIgnoreHostKey(), nil
	}

	// Create the callback from the known hosts file
	callback, err := knownhosts.New(config.KnownHostsPath)
	if err != nil {
		return nil, err
	}

	return callback, nil
}

// authenticationMethod generates authentication methods based on the SSH configuration.
func (config *Config) authenticationMethod() ([]ssh.AuthMethod, error) {
	var authMethod []ssh.AuthMethod = []ssh.AuthMethod{}

	// Private key authentication
	if config.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, ssh.PublicKeys(signer))
	} else if config.PrivateKeyPath != "" {
		privateKey, err := os.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, ssh.PublicKeys(signer))
	}

	// Password authentication
	if config.Password != "" {
		authMethod = append(authMethod, ssh.Password(config.Password))
	}

	return authMethod, nil
}
