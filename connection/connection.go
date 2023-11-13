package connection

import (
	"context"
	"errors"

	"github.com/masterzen/winrm"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	WinRM *winrm.Client
	SSH   *ssh.Client
}

type Config struct {
	WinRM *WinRMConfig
	SSH   *SSHConfig
}

type CMDResult struct {
	StdOut string
	StdErr string
}

// New returns a Connection object.
// If WinRMConfig is specified the Connection contains a WinRM conenction.
// If SSHConfig is specified the Connection contains a SSH conenction.
func New(conf *Config) (*Connection, error) {

	// Assert WinRM and SSH configuration
	if conf.WinRM == nil && conf.SSH == nil {
		return nil, errors.New("one of WinRMConfig and SSHConfig must be set")
	}
	if conf.WinRM != nil && conf.SSH != nil {
		return nil, errors.New("only one of WinRMConfig and SSHConfig must be set")
	}

	// Allocate a new Connection
	c := new(Connection)

	// WinRM configuration
	if conf.WinRM != nil {
		winRMClient, err := newWinRMClient(conf.WinRM)
		if err != nil {
			return nil, err
		}

		c = &Connection{
			WinRM: winRMClient,
		}
	}

	// SSH configuration
	if conf.SSH != nil {
		sshClient, err := newSSHClient(conf.SSH)
		if err != nil {
			return nil, err
		}

		c = &Connection{
			SSH: sshClient,
		}
	}

	return c, nil
}

// Close closes any open connection.
func (c *Connection) Close() error {
	if c.SSH != nil {
		err := c.SSH.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Run runs a command with a connection and context
// It returns stdout and stderr within a CMDResult object
func (c *Connection) Run(ctx context.Context, cmd string) (*CMDResult, error) {

	// Prepare base64 encoded powershell command to pass into the run functions
	pwshCmd := winrm.Powershell(cmd)

	// Allocate CMDResult
	r := new(CMDResult)

	// WinRM execution
	if c.WinRM != nil {
		stdout, stderr, _, err := c.WinRM.RunWithContextWithString(ctx, pwshCmd, "")
		if err != nil {
			return nil, err
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
