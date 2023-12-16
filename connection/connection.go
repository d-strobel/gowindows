package connection

import (
	"context"
	"fmt"

	"github.com/d-strobel/winrm"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	WinRM *winrm.Client
	SSH   *ssh.Client
}

type ConnectionInterface interface {
	Run(ctx context.Context, cmd string) (CMDResult, error)
	Close() error
}

type Config struct {
	WinRM *WinRMConfig
	SSH   *SSHConfig
}

type CMDResult struct {
	StdOut string
	StdErr string
}

// NewConnection returns a Connection object.
// If WinRMConfig is specified the Connection object contains a WinRM connection.
// If SSHConfig is specified the Connection object contains a SSH connection.
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

// Close closes any open connection.
func (c *Connection) Close() error {
	if c.SSH != nil {
		err := c.SSH.Close()
		if err != nil {
			return fmt.Errorf("connection: %s", err)
		}
	}

	return nil
}

// Run runs a command with a connection and context.
// It returns stdout and stderr within a CMDResult object.
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
