package winrm

import (
	"errors"
	"time"
)

// Default values for WinRM configuration.
const (
	defaultPort     int           = 5985
	defaultPortTLS  int           = 5986
	defaultUseTLS   bool          = false
	defaultInsecure bool          = false
	defaultTimeout  time.Duration = 0
)

// Config represents the configuration details for establishing a WinRM connection.
type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	UseTLS   bool
	Insecure bool
	Timeout  time.Duration
}

// validate validates the WinRM configuration.
func (config *Config) validate() error {
	if config.Host == "" || config.Username == "" || config.Password == "" {
		return errors.New("winrm: Config parameter 'Host', 'Username', and 'Password' must be set")
	}

	return nil
}

// defaults sets the default values for the WinRM configuration.
func (config *Config) defaults() error {
	if !config.UseTLS {
		config.UseTLS = defaultUseTLS
	}

	if config.Port == 0 {
		config.Port = defaultPort

		// Set a different default port if TLS enabled
		if config.UseTLS {
			config.Port = defaultPortTLS
		}
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	if !config.Insecure {
		config.Insecure = defaultInsecure
	}

	return nil
}
