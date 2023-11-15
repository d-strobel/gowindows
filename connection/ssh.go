package connection

import (
	"context"
	"fmt"
	"io"

	"github.com/d-strobel/gowindows/winerror"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	SSHHost     string
	SSHPort     int
	SSHUsername string
	SSHPassword string
}

const (
	// SSH default values
	defaultSSHPort int = 22
)

func newSSHClient(config *SSHConfig) (*ssh.Client, error) {

	// Assert
	if config.SSHHost == "" || config.SSHUsername == "" || config.SSHPassword == "" {
		return nil, winerror.Errorf(winerror.ConfigError, "ssh client: SSHConfig parameter 'SSHHost', 'SSHUsername' and 'SSHPassword' must be set")
	}

	// Parse SSH host string
	sshHost := fmt.Sprintf("%s:%d", config.SSHHost, config.SSHPort)

	// Configuration
	sshConfig := &ssh.ClientConfig{
		User: config.SSHUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote server and perform the SSH handshake
	client, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, winerror.Errorf(winerror.ConnectionError, "ssh client: %s", err)
	}

	return client, nil
}

func (c *Connection) runSSH(ctx context.Context, cmd string) (string, string, error) {

	// Open a new SSH session
	s, err := c.SSH.NewSession()
	if err != nil {
		return "", "", err
	}
	defer s.Close()

	// Create pipes to capture stdout and stderr
	stdout, err := s.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderr, err := s.StderrPipe()
	if err != nil {
		return "", "", err
	}

	// Run the command
	err = s.Start(cmd)
	if err != nil {
		return "", "", err
	}

	// Read output from pipes
	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return "", "", err
	}
	stderrBytes, err := io.ReadAll(stderr)
	if err != nil {
		return "", "", err
	}

	// Wait for the command to complete with context support
	select {
	case <-ctx.Done():
		_ = s.Signal(ssh.SIGINT)
		return "", "", ctx.Err()
	default:
		err = s.Wait()
	}

	// Return the error if stderr has no value
	if err != nil && len(stderrBytes) == 0 {
		return "", "", err
	}

	// Return stderr over the error when stderr has a value
	if len(stderrBytes) > 0 {
		return "", string(stderrBytes), nil
	}

	// Return error when stdout and stderr have no values
	if len(stdoutBytes) == 0 && len(stderrBytes) == 0 {
		return "", "", winerror.Errorf(winerror.WindowsError, "ssh session: stdout and stderr are empty")
	}

	return string(stdoutBytes), "", nil
}
