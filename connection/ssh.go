package connection

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSHConfig represents the configuration details for establishing an SSH connection.
type SSHConfig struct {
	SSHHost                  string
	SSHPort                  int
	SSHUsername              string
	SSHPassword              string
	SSHPrivateKey            string
	SSHPrivateKeyPath        string
	SSHKnownHostsPath        string
	SSHInsecureIgnoreHostKey bool
}

// Default values for SSH configuration.
const (
	defaultSSHPort        int    = 22
	defaultKnownHostsPath string = ".ssh/known_hosts"
)

// newSSHClient creates a new SSH client based on the provided configuration.
func newSSHClient(config *SSHConfig) (*ssh.Client, error) {

	// Assert
	if (config.SSHHost == "" || config.SSHUsername == "") || (config.SSHPassword == "" && config.SSHPrivateKey == "" && config.SSHPrivateKeyPath == "") {
		return nil, fmt.Errorf("ssh: SSHConfig parameter 'SSHHost', 'SSHUsername' and one of 'SSHPassword', 'SSHPrivateKey', 'SSHPrivateKeyPath' must be set")
	}

	// Parse SSH host string
	sshHost := fmt.Sprintf("%s:%d", config.SSHHost, config.SSHPort)

	// Check known host key callback
	knownHostCallback, err := knownHostCallback(config)
	if err != nil {
		return nil, fmt.Errorf("ssh: known host callback failed with error: %s", err)
	}

	// Authentication method
	authMethod, err := authenticationMethod(config)
	if err != nil {
		return nil, fmt.Errorf("ssh: %s", err)
	}

	// Configuration
	sshConfig := &ssh.ClientConfig{
		User:            config.SSHUsername,
		Auth:            authMethod,
		HostKeyCallback: knownHostCallback,
	}

	// Connect to the remote server and perform the SSH handshake
	client, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh: %s", err)
	}

	return client, nil
}

// runSSH runs a command on the SSH connection and returns the stdout and stderr.
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
	if err := s.Start(cmd); err != nil {
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

	// Some Powershell functions does not provide any output.
	// In that case we return empty strings.
	if len(stdoutBytes) == 0 && len(stderrBytes) == 0 {
		return "", "", nil
	}

	return string(stdoutBytes), "", nil
}

// knownHostCallback generates a host key callback based on the SSH configuration.
func knownHostCallback(config *SSHConfig) (ssh.HostKeyCallback, error) {

	// Ignore host key
	if config.SSHInsecureIgnoreHostKey {
		return ssh.InsecureIgnoreHostKey(), nil
	}

	// Get the current user from the system
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	// Set default values
	knownHostsPath := fmt.Sprintf("%s/%s", user.HomeDir, defaultKnownHostsPath)
	if config.SSHKnownHostsPath != "" {
		knownHostsPath = config.SSHKnownHostsPath
	}

	// Create the callback from the known hosts file
	callback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, err
	}

	return callback, nil
}

// authenticationMethod generates authentication methods based on the SSH configuration.
func authenticationMethod(config *SSHConfig) ([]ssh.AuthMethod, error) {
	var authMethod []ssh.AuthMethod = []ssh.AuthMethod{}

	// Private key authentication
	if config.SSHPrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(config.SSHPrivateKey))
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, ssh.PublicKeys(signer))
	} else if config.SSHPrivateKeyPath != "" {
		privateKey, err := os.ReadFile(config.SSHPrivateKeyPath)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, ssh.PublicKeys(signer))
	}

	// Password authentication
	if config.SSHPassword != "" {
		authMethod = append(authMethod, ssh.Password(config.SSHPassword))
	}

	return authMethod, nil
}
