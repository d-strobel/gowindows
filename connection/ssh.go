package connection

import (
	"context"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

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

	// Wait for the command to complete
	err = s.Wait()
	if err != nil {
		return "", "", err
	}

	return string(stdoutBytes), string(stderrBytes), nil
}
