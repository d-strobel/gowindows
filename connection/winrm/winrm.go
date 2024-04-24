package winrm

import (
	"context"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/winrm"
)

// Connection represents a WinRM connection.
type Connection struct {
	Client *winrm.Client
}

// NewConnectionWithKerberos creates a new WinRM client with Kerberos authentication.
// It takes a WinRM configuration and a Kerberos configuration as input parameters.
func NewConnectionWithKerberos(config *Config, krbConfig *KerberosConfig) (*Connection, error) {
	// Validate configuration for WinRM and Kerberos.
	if err := config.validate(); err != nil {
		return nil, err
	}
	if err := krbConfig.validate(); err != nil {
		return nil, err
	}

	// Set default values for WinRM and Kerberos.
	if err := config.defaults(); err != nil {
		return nil, err
	}
	if err := krbConfig.defaults(); err != nil {
		return nil, err
	}

	// WinRM connection
	winRMEndpoint := winrm.NewEndpoint(
		config.Host,
		config.Port,
		config.UseTLS,
		config.Insecure,
		nil, // CA certificate
		nil, // Client Certificate
		nil, // Client Key
		config.Timeout,
	)

	params := krbConfig.kerberosParams(config)

	// Create a new WinRM client with Kerberos authentication.
	client, err := winrm.NewClientWithParameters(winRMEndpoint, config.Username, config.Password, params)
	if err != nil {
		return nil, err
	}

	return &Connection{Client: client}, nil
}

// NewConnection creates a new WinRM client based on the provided WinRM configuration.
func NewConnection(config *Config) (*Connection, error) {
	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	// Set default values
	if err := config.defaults(); err != nil {
		return nil, err
	}

	// WinRM connection
	winRMEndpoint := winrm.NewEndpoint(
		config.Host,
		config.Port,
		config.UseTLS,
		config.Insecure,
		nil, // CA certificate
		nil, // Client Certificate
		nil, // Client Key
		config.Timeout,
	)

	// Create a new WinRM client.
	client, err := winrm.NewClient(winRMEndpoint, config.Username, config.Password)
	if err != nil {
		return nil, err
	}

	return &Connection{Client: client}, nil
}

// Close closes the WinRM connection.
// Satisfies the Connection interface.
func (c *Connection) Close() error {
	return nil
}

// RunWithPowershell runs a command using the configured WinRM connection and context via Powershell.
// It returns the result of the command execution, including stdout and stderr.
func (c *Connection) RunWithPowershell(ctx context.Context, cmd string) (connection.CMDResult, error) {
	// Prepare powershell command.
	pwshCmd := winrm.Powershell(cmd)
	return c.Run(ctx, pwshCmd)
}

// Run runs a command using the configured WinRM connection and context.
// It returns a connection.CMDResult object, including stdout and stderr.
func (c *Connection) Run(ctx context.Context, cmd string) (connection.CMDResult, error) {
	var r connection.CMDResult

	stdout, stderr, _, err := c.Client.RunWithContextWithString(ctx, cmd, "")
	if err != nil {
		return r, err
	}

	r.StdErr = stderr
	r.StdOut = stdout

	return r, nil
}
