// Package ssh offers utilities for creating and managing SSH (Secure Shell) connections, enabling secure communication with Windows operating systems.
// It allows executing commands on remote hosts securely and efficiently.
//
// Key Features:
//   - Establishes SSH connections with remote hosts based on provided configuration.
//   - Handles authentication mechanisms such as password-based and privatekey-based authentication.
//   - Supports execution of commands including cmd and powershell commands.
package ssh

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/winrm"
	"golang.org/x/crypto/ssh"
)

// Connection represents an SSH connection.
// It holds a client object for interacting with the remote system.
type Connection struct {
	Client *ssh.Client
}

// NewConnection creates a new SSH client based on the provided configuration.
func NewConnection(config *Config) (*Connection, error) {

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	// Set default values
	if err := config.defaults(); err != nil {
		return nil, err
	}

	// Parse SSH host string
	sshHost := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// Check known host key callback
	knownHostCallback, err := config.knownHostCallback()
	if err != nil {
		return nil, fmt.Errorf("ssh: known host callback failed with error: %s", err)
	}

	// Authentication method
	authMethod, err := config.authenticationMethod()
	if err != nil {
		return nil, fmt.Errorf("ssh: authentication method failed with error: %s", err)
	}

	// Configuration
	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            authMethod,
		HostKeyCallback: knownHostCallback,
	}

	// Connect to the remote server and perform the SSH handshake
	client, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh: %s", err)
	}

	return &Connection{Client: client}, nil
}

// Close closes the SSH connection.
func (c *Connection) Close() error {
	return c.Client.Close()
}

// RunWithPowershell runs a command using the configured SSH connection and context via Powershell.
func (c *Connection) RunWithPowershell(ctx context.Context, cmd string) (connection.CmdResult, error) {
	// Prepare base64 encoded powershell command
	pwshCmd := winrm.Powershell(cmd)
	return c.Run(ctx, pwshCmd)
}

// Run runs a command using the configured SSH connection and context.
// It returns the result of the command execution, including stdout and stderr.
func (c *Connection) Run(ctx context.Context, cmd string) (connection.CmdResult, error) {
	var r connection.CmdResult

	// Open a new SSH session.
	s, err := c.Client.NewSession()
	if err != nil {
		return r, err
	}
	defer s.Close()

	// Prepare channels for stdout, stderr and errors.
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	errChan := make(chan error)

	// Use a WaitGroup to wait for the goroutines to finish.
	var wg sync.WaitGroup

	// Goroutine to handle stdout.
	wg.Add(1)
	go func() {
		defer wg.Done()

		stdout, err := s.StdoutPipe()
		if err != nil {
			errChan <- err
			return
		}

		stdoutBytes, err := io.ReadAll(stdout)
		if err != nil {
			errChan <- err
			return
		}
		stdoutChan <- string(stdoutBytes)
	}()

	// Goroutine to handle stderr.
	wg.Add(1)
	go func() {
		defer wg.Done()

		stderr, err := s.StderrPipe()
		if err != nil {
			errChan <- err
			return
		}

		stderrBytes, err := io.ReadAll(stderr)
		if err != nil {
			errChan <- err
			return
		}
		stderrChan <- string(stderrBytes)
	}()

	// Start the command execution.
	if err := s.Start(cmd); err != nil {
		return r, err
	}

	// Wait for the goroutines to finish.
	go func() {
		wg.Wait()
		close(stdoutChan)
		close(stderrChan)
		close(errChan)
	}()

	// Wait for the command to complete with context support.
	select {
	case <-ctx.Done():
		_ = s.Signal(ssh.SIGINT)
		return r, ctx.Err()
	case err := <-errChan:
		return r, err
	case r.StdOut = <-stdoutChan:
		r.StdErr = <-stderrChan
		return r, nil
	}
}
