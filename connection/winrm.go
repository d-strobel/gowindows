package connection

import (
	"fmt"
	"time"

	"github.com/d-strobel/winrm"
)

// WinRMConfig represents the configuration details for establishing a WinRM connection.
type WinRMConfig struct {
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMUseTLS   bool
	WinRMInsecure bool
	WinRMTimeout  time.Duration
	WinRMKerberos *KerberosConfig
}

// Default values for WinRM configuration.
const (
	defaultWinRMPort     int           = 5985
	defaultWinRMPortTLS  int           = 5986
	defaultWinRMUseTLS   bool          = false
	defaultWinRMInsecure bool          = true
	defaultWinRMTimeout  time.Duration = 0
)

// newWinRMClient creates a new WinRM client based on the provided configuration.
func newWinRMClient(config *WinRMConfig) (*winrm.Client, error) {

	// Assert
	if config.WinRMHost == "" || config.WinRMUsername == "" || config.WinRMPassword == "" {
		return nil, fmt.Errorf("winrm: WinRMConfig parameter 'WinRMHost', 'WinRMUsername', and 'WinRMPassword' must be set")
	}

	// Set default values
	if !config.WinRMUseTLS {
		config.WinRMUseTLS = defaultWinRMUseTLS
	}
	if config.WinRMPort == 0 {
		config.WinRMPort = defaultWinRMPort

		// Set a different default port if TLS enabled
		if config.WinRMUseTLS {
			config.WinRMPort = defaultWinRMPortTLS
		}
	}
	if config.WinRMTimeout == 0 {
		config.WinRMTimeout = defaultWinRMTimeout
	}

	winRMInsecure := defaultWinRMInsecure
	if config.WinRMInsecure {
		winRMInsecure = config.WinRMInsecure
	}

	// WinRM connection
	winRMEndpoint := winrm.NewEndpoint(
		config.WinRMHost,
		config.WinRMPort,
		config.WinRMUseTLS,
		winRMInsecure,
		nil, // CA certificate
		nil, // Client Certificate
		nil, // Client Key
		config.WinRMTimeout,
	)

	// Kerberos transport
	if config.WinRMKerberos != nil {
		params := winRMKerberosParams(config)
		return winrm.NewClientWithParameters(winRMEndpoint, config.WinRMUsername, config.WinRMPassword, params)
	}

	return winrm.NewClient(winRMEndpoint, config.WinRMUsername, config.WinRMPassword)
}
