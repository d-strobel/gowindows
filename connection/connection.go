package connection

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/masterzen/winrm"
	"golang.org/x/crypto/ssh"
)

const (
	// WinRM default values
	defaultWinRMPort     int           = 5986
	defaultWinRMUseTLS   bool          = false
	defaultWinRMInsecure bool          = true
	defaultWinRMTimeout  time.Duration = 0

	// SSH default values
	defaultSSHPort int = 22
)

type Connection struct {
	WinRM *winrm.Client
	SSH   *ssh.Client
}

type Config struct {
	WinRM *WinRMConfig
	SSH   *SSHConfig
}

type WinRMConfig struct {
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMUseTLS   bool
	WinRMInsecure bool
	WinRMTimeout  time.Duration
}

type SSHConfig struct {
	SSHHost     string
	SSHPort     int
	SSHUsername string
	SSHPassword string
}

// New returns a Connection object.
// If WinRMConfig is specified then the Connection contains a WinRM conenction.
// If SSHConfig is specified then the Connection contains a SSH conenction.
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

func (c *Connection) Run(ctx context.Context, cmd string) (string, error) {
	var result string

	// Get encoded powershell command to execute over WinRM or SSH
	pwshCmd := winrm.Powershell(cmd)

	// WinRM execution
	if c.WinRM != nil {
		stdout, _, _, err := c.WinRM.RunWithContextWithString(ctx, pwshCmd, "")
		if err != nil {
			return "", err
		}

		result = stdout
	}

	// SSH execution
	if c.SSH != nil {
		// Open a new SSH session
		s, err := c.SSH.NewSession()
		if err != nil {
			return "", err
		}
		defer s.Close()

		// Create pipes to capture STDOUT and STDERR
		stdout, err := s.StdoutPipe()
		if err != nil {
			return "", err
		}
		stderr, err := s.StderrPipe()
		if err != nil {
			return "", err
		}

		// Run the command
		err = s.Start(pwshCmd)
		if err != nil {
			return "", nil
		}

		// Read output from pipes
		stdoutBytes, err := io.ReadAll(stdout)
		if err != nil {
			return "", nil
		}
		stderrBytes, err := io.ReadAll(stderr)
		if err != nil {
			return "", nil
		}

		// Wait for the command to complete
		err = s.Wait()
		if err != nil {
			return "", nil
		}

		if stderrBytes == nil {
			return "", errors.New(fmt.Sprintf("Command failed with following error:\n%s", string(stderrBytes)))
		}

		result = string(stdoutBytes)
	}

	return result, nil
}
