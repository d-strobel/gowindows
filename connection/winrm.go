package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/d-strobel/winrm"
)

// Default values for WinRM configuration.
const (
	defaultWinRMPort     int           = 5985
	defaultWinRMPortTLS  int           = 5986
	defaultWinRMUseTLS   bool          = false
	defaultWinRMInsecure bool          = false
	defaultWinRMTimeout  time.Duration = 0
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

// WinRMConnection represents a WinRM connection.
type WinRMConnection struct {
	Client *winrm.Client
}

// validate validates the WinRM configuration.
func (config *WinRMConfig) validate() error {

	if config.WinRMHost == "" || config.WinRMUsername == "" || config.WinRMPassword == "" {
		return fmt.Errorf("winrm: WinRMConfig parameter 'WinRMHost', 'WinRMUsername', and 'WinRMPassword' must be set")
	}

	return nil
}

// defaults sets the default values for the WinRM configuration.
func (config *WinRMConfig) defaults() error {

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

	if !config.WinRMInsecure {
		config.WinRMInsecure = defaultWinRMInsecure
	}

	return nil
}

// NewConnectionWithWinRM creates a new WinRM client based on the provided WinRM configuration.
func NewConnectionWithWinRM(config *WinRMConfig) (*WinRMConnection, error) {

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	// Set default values
	if err := config.defaults(); err != nil {
		return nil, err
	}

	// Initialize WinRMConnection object
	conn := &WinRMConnection{}

	// WinRM connection
	winRMEndpoint := winrm.NewEndpoint(
		config.WinRMHost,
		config.WinRMPort,
		config.WinRMUseTLS,
		config.WinRMInsecure,
		nil, // CA certificate
		nil, // Client Certificate
		nil, // Client Key
		config.WinRMTimeout,
	)

	// Kerberos transport
	if config.WinRMKerberos != nil {
		params := config.winRMKerberosParams()

		client, err := winrm.NewClientWithParameters(winRMEndpoint, config.WinRMUsername, config.WinRMPassword, params)
		if err != nil {
			return nil, err
		}

		conn.Client = client

	} else {
		client, err := winrm.NewClient(winRMEndpoint, config.WinRMUsername, config.WinRMPassword)

		if err != nil {
			return nil, err
		}

		conn.Client = client
	}

	return conn, nil
}

// Close closes the WinRM connection.
// Satisfies the Connection interface.
func (c *WinRMConnection) Close() error {
	return nil
}

// Run runs a command using the configured WinRM connection and context.
// It returns the result of the command execution, including stdout and stderr.
func (c *WinRMConnection) Run(ctx context.Context, cmd string) (CMDResult, error) {

	var r CMDResult

	// Prepare base64 encoded powershell command to pass into the run functions
	pwshCmd := winrm.Powershell(cmd)

	stdout, stderr, _, err := c.Client.RunWithContextWithString(ctx, pwshCmd, "")
	if err != nil {
		return r, err
	}
	if stderr != "" {
		r.StdErr = stderr
		return r, nil
	}

	r.StdOut = stdout
	return r, nil
}
