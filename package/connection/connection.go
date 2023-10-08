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

func NewConnection(conf *Config) (*Connection, error) {

	// Assert
	if conf.WinRM == nil && conf.SSH == nil {
		return nil, errors.New("one of WinRMConfig and SSHConfig must be set")
	}
	if conf.WinRM != nil && conf.SSH != nil {
		return nil, errors.New("only one of WinRMConfig and SSHConfig must be set")
	}

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

func newWinRMClient(config *WinRMConfig) (*winrm.Client, error) {

	// Assert
	if config.WinRMHost == "" || config.WinRMUsername == "" || config.WinRMPassword == "" {
		return nil, errors.New("WinRMHost, WinRMUsername, and WinRMPassword must be set")
	}

	// Set default values
	if config.WinRMPort == 0 {
		config.WinRMPort = defaultWinRMPort
	}
	if !config.WinRMUseTLS {
		config.WinRMUseTLS = defaultWinRMUseTLS
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

	winRMClient, err := winrm.NewClient(winRMEndpoint, config.WinRMUsername, config.WinRMPassword)
	if err != nil {
		return nil, err
	}

	return winRMClient, nil
}

func newSSHClient(config *SSHConfig) (*ssh.Client, error) {

	// Assert
	if config.SSHHost == "" || config.SSHUsername == "" || config.SSHPassword == "" {
		return nil, errors.New("SSHHost, SSHUsername, and SSHPassword must be set")
	}

	// Parse ssh host string
	sshHost := fmt.Sprintf("%s:%s", config.SSHHost, fmt.Sprint(config.SSHPort))

	// Configuration
	sshConfig := &ssh.ClientConfig{
		User: config.SSHUsername,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.Password(config.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Close closes any open connection.
// Only ssh connection will be terminated here.
// To avoid surprises in the future, this should always be called in a defer statement.
func (conn *Connection) Close() error {
	if conn.SSH != nil {
		err := conn.SSH.Close()
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
