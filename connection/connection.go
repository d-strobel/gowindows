// Package connection provides a Go library for handling connections to Windows-based systems using WinRM and SSH protocols.
package connection

import (
	"context"
	"fmt"

	"github.com/d-strobel/winrm"
	"golang.org/x/crypto/ssh"
)

// Connection represents a connection object that can be used to interact with a Windows system.
type Connection struct {
	WinRM *winrm.Client
	SSH   *ssh.Client
}

// ConnectionInterface defines the interface for a connection, specifying methods like Run and Close.
type ConnectionInterface interface {
	Run(ctx context.Context, cmd string) (CMDResult, error)
	Close() error
}

// Config contains configuration details for creating a Connection object.
type Config struct {
	WinRM *WinRMConfig
	SSH   *SSHConfig
}

// CMDResult represents the result of executing a command, including stdout and stderr.
type CMDResult struct {
	StdOut string
	StdErr string
}

// NewConnection returns a Connection object based on the provided configuration.
// The returned Connection object may contain either a WinRM or SSH connection.
func NewConnection(conf *Config) (*Connection, error) {

	// Assert WinRM and SSH configuration
	if conf.WinRM == nil && conf.SSH == nil {
		return nil, fmt.Errorf("connection: Connection object 'WinRMConfig' or 'SSHConfig' must be set")
	}
	if conf.WinRM != nil && conf.SSH != nil {
		return nil, fmt.Errorf("connection: Connection object must only contain 'WinRMConfig' or 'SSHConfig'")
	}

	// Init a new Connection
	c := &Connection{}

	// WinRM configuration
	if conf.WinRM != nil {
		winRMClient, err := newWinRMClient(conf.WinRM)
		if err != nil {
			return nil, err
		}

		c.WinRM = winRMClient
	}

	// SSH configuration
	if conf.SSH != nil {
		sshClient, err := newSSHClient(conf.SSH)
		if err != nil {
			return nil, err
		}

		c.SSH = sshClient
	}

	return c, nil
}

// Close closes any open connection, whether it's WinRM or SSH.
func (c *Connection) Close() error {
	if c.SSH != nil {
		if err := c.SSH.Close(); err != nil {
			return fmt.Errorf("connection: %s", err)
		}
	}

	return nil
}

// Run runs a command using the configured connection and context.
// It returns the result of the command execution, including stdout and stderr.
func (c *Connection) Run(ctx context.Context, cmd string) (CMDResult, error) {

	var r CMDResult

	// Assert configuration
	if c.WinRM == nil && c.SSH == nil {
		return r, fmt.Errorf("connection: Connection object 'WinRMConfig' or 'SSHConfig' must be set")
	}

	// Prepare base64 encoded powershell command to pass into the run functions
	pwshCmd := winrm.Powershell(cmd)

	// WinRM execution
	if c.WinRM != nil {
		stdout, stderr, _, err := c.WinRM.RunWithContextWithString(ctx, pwshCmd, "")
		if err != nil {
			return r, err
		}
		if stderr != "" {
			r.StdErr = stderr
			return r, nil
		}

		r.StdOut = stdout
	}

	// SSH execution
	if c.SSH != nil {
		stdout, stderr, err := c.runSSH(ctx, pwshCmd)
		if err != nil {
			r.StdErr = stderr
			return r, err
		}
		if stderr != "" {
			r.StdErr = stderr
			return r, nil
		}

		r.StdOut = stdout
	}

	return r, nil
}
