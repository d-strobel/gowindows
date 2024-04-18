package connection_test

import (
	"github.com/d-strobel/gowindows/connection"
)

func (suite *ConnectionAccTestSuite) TestNewConnectionWithSSH() {
	suite.T().Parallel()

	suite.Run("should establish a connection via password", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:     suite.host,
			SSHPort:     suite.sshPort,
			SSHUsername: suite.username,
			SSHPassword: suite.password,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey path with ed25519", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:           suite.host,
			SSHPort:           suite.sshPort,
			SSHUsername:       suite.username,
			SSHPrivateKeyPath: suite.sshKeyPathED25519,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey path with rsa", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:           suite.host,
			SSHPort:           suite.sshPort,
			SSHUsername:       suite.username,
			SSHPrivateKeyPath: suite.sshKeyPathRSA,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey with ed25519", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:       suite.host,
			SSHPort:       suite.sshPort,
			SSHUsername:   suite.username,
			SSHPrivateKey: suite.sshKeyED25519,
		}

		suite.T().Logf("key: '%s'", suite.sshKeyED25519)
		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey with rsa", func() {
		sshConfig := connection.SSHConfig{
			SSHHost:       suite.host,
			SSHPort:       suite.sshPort,
			SSHUsername:   suite.username,
			SSHPrivateKey: suite.sshKeyRSA,
		}

		conn, err := connection.NewConnectionWithSSH(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})
}
