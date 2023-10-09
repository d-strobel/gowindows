package connection

import (
	"errors"
	"fmt"

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
