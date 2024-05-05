// Package winrm provides functionality for establishing and managing WinRM (Windows Remote Management) connections.
// It facilitates executing commands on remote Windows machines securely over HTTP(S) using the WinRM protocol.
//
// Key Features:
//   - Establishes WinRM connections based on provided configuration.
//   - Handles authentication and secure communication with remote Windows hosts.
//   - Supports execution of commands including cmd and powershell commands.
package winrm

import (
	"context"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parsing"
	"github.com/masterzen/winrm"
)

// Connection represents a WinRM connection.
type Connection struct {
	Client *winrm.Client
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
func (c *Connection) RunWithPowershell(ctx context.Context, cmd string) (connection.CmdResult, error) {
	// Prepare powershell command.
	pwshCmd, err := parsing.EncodePwshCmd(cmd)
	if err != nil {
		return connection.CmdResult{}, err
	}

	return c.Run(ctx, pwshCmd)
}

// Run runs a command using the configured WinRM connection and context.
// It returns a connection.CMDResult object, including stdout and stderr.
func (c *Connection) Run(ctx context.Context, cmd string) (connection.CmdResult, error) {
	var r connection.CmdResult

	stdout, stderr, _, err := c.Client.RunWithContextWithString(ctx, cmd, "")
	if err != nil {
		return r, err
	}

	r.StdErr = stderr
	r.StdOut = stdout

	return r, nil
}
