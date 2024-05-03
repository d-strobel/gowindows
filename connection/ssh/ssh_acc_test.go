package ssh_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/d-strobel/gowindows/connection/ssh"
	"github.com/stretchr/testify/suite"
)

// Init acceptance test suite for SSH
type SSHAccTestSuite struct {
	suite.Suite

	// Fixtures
	host                  string
	username              string
	password              string
	port                  int
	privateKeyPathED25519 string
	privateKeyED25519     string
	privateKeyPathRSA     string
	privateKeyRSA         string
	connection            *ssh.Connection
	adHost                string
	adUsernamePre2k       string
	adPassword            string
	adPort                int
}

// SetupSuite setups all neccessary fixtures for running the ssh tests.
func (suite *SSHAccTestSuite) SetupSuite() {
	var err error

	// Load environment variables
	suite.host = os.Getenv("GOWINDOWS_TEST_HOST")
	suite.Require().NotEmpty(suite.host, "Environment variable not set: GOWINDOWS_TEST_HOST")

	suite.username = os.Getenv("GOWINDOWS_TEST_USERNAME")
	suite.Require().NotEmpty(suite.username, "Environment variable not set: GOWINDOWS_TEST_USERNAME")

	suite.password = os.Getenv("GOWINDOWS_TEST_PASSWORD")
	suite.Require().NotEmpty(suite.password, "Environment variable not set: GOWINDOWS_TEST_PASSWORD")

	suite.port, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_SSH_PORT"))
	suite.Require().NoError(err)

	suite.privateKeyPathED25519 = os.Getenv("GOWINDOWS_TEST_SSH_PRIVATE_KEY_ED25519_PATH")
	suite.Require().NotEmpty(suite.privateKeyPathED25519, "Environment variable not set: GOWINDOWS_TEST_SSH_PRIVATE_KEY_ED25519_PATH")

	suite.privateKeyPathRSA = os.Getenv("GOWINDOWS_TEST_SSH_PRIVATE_KEY_RSA_PATH")
	suite.Require().NotEmpty(suite.privateKeyPathRSA, "Environment variable not set: GOWINDOWS_TEST_SSH_PRIVATE_KEY_RSA_PATH")

	privateKeyRSA, err := os.ReadFile(suite.privateKeyPathRSA)
	suite.Require().NoError(err)
	suite.privateKeyRSA = string(privateKeyRSA)

	privateKeyED25519, err := os.ReadFile(suite.privateKeyPathED25519)
	suite.Require().NoError(err)
	suite.privateKeyED25519 = string(privateKeyED25519)

	suite.adHost = os.Getenv("GOWINDOWS_TEST_AD_HOST")
	suite.Require().NotEmpty(suite.host, "Environment variable not set: GOWINDOWS_TEST_AD_HOST")

	suite.adUsernamePre2k = os.Getenv("GOWINDOWS_TEST_AD_USERNAME_PRE2K")
	suite.Require().NotEmpty(suite.username, "Environment variable not set: GOWINDOWS_TEST_AD_USERNAME_PRE2K")

	suite.adPassword = os.Getenv("GOWINDOWS_TEST_AD_PASSWORD")
	suite.Require().NotEmpty(suite.password, "Environment variable not set: GOWINDOWS_TEST_AD_PASSWORD")

	suite.adPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_AD_SSH_PORT"))
	suite.Require().NoError(err)

	// Setup SSH connection
	sshConfig := &ssh.Config{
		Host:     suite.host,
		Port:     suite.port,
		Username: suite.username,
		Password: suite.password,
	}
	suite.connection, err = ssh.NewConnection(sshConfig)
	suite.Require().NoError(err)
}

// TearDownSuite closes all ssh connections after running the tests.
func (suite *SSHAccTestSuite) TearDownSuite() {
	// Close connections
	suite.connection.Close()
}

// TestSSHAccTestSuite runs the acceptance test suite for the ssh package.
// It will be skipped if the short flag is set.
func TestSSHAccTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, &SSHAccTestSuite{})
}

func (suite *SSHAccTestSuite) TestNewConnection() {
	suite.Run("should establish a connection via password", func() {
		sshConfig := ssh.Config{
			Host:     suite.host,
			Port:     suite.port,
			Username: suite.username,
			Password: suite.password,
		}

		conn, err := ssh.NewConnection(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey path with ed25519", func() {
		sshConfig := ssh.Config{
			Host:           suite.host,
			Port:           suite.port,
			Username:       suite.username,
			PrivateKeyPath: suite.privateKeyPathED25519,
		}

		conn, err := ssh.NewConnection(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey path with rsa", func() {
		sshConfig := ssh.Config{
			Host:           suite.host,
			Port:           suite.port,
			Username:       suite.username,
			PrivateKeyPath: suite.privateKeyPathRSA,
		}

		conn, err := ssh.NewConnection(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey with ed25519", func() {
		sshConfig := ssh.Config{
			Host:       suite.host,
			Port:       suite.port,
			Username:   suite.username,
			PrivateKey: suite.privateKeyED25519,
		}
		conn, err := ssh.NewConnection(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	suite.Run("should establish a connection via privatekey with rsa", func() {
		sshConfig := ssh.Config{
			Host:       suite.host,
			Port:       suite.port,
			Username:   suite.username,
			PrivateKey: suite.privateKeyRSA,
		}

		conn, err := ssh.NewConnection(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})

	// Domaincontroller tests
	suite.Run("domain: should establish a connection via password", func() {
		sshConfig := ssh.Config{
			Host:                  suite.adHost,
			Port:                  suite.adPort,
			Username:              suite.adUsernamePre2k,
			Password:              suite.adPassword,
			InsecureIgnoreHostKey: true,
		}

		conn, err := ssh.NewConnection(&sshConfig)
		suite.Assertions.NoError(err)
		conn.Close()
	})
}

func (suite *SSHAccTestSuite) TestRunWithPowershell() {
	suite.Run("should return a valid output", func() {
		ctx := context.Background()
		result, err := suite.connection.RunWithPowershell(ctx, "Get-LocalUser")
		suite.Assertions.NoError(err)
		suite.Assertions.NotEmpty(result.StdOut)
		suite.Assertions.Empty(result.StdErr)
	})

	suite.Run("should return a non valid output", func() {
		ctx := context.Background()
		result, err := suite.connection.RunWithPowershell(ctx, "Get-ocalser")
		suite.Assertions.NoError(err)
		suite.Assertions.Empty(result.StdOut)
		suite.Assertions.NotEmpty(result.StdErr)
	})

	suite.Run("should return error with context canceled", func() {
		ctx, cancel := context.WithCancel(context.Background())
		// Cancel the context immediately
		cancel()
		result, err := suite.connection.RunWithPowershell(ctx, "Start-Sleep -Seconds 20")
		suite.Assertions.ErrorContains(err, "context canceled")
		suite.Assertions.Empty(result.StdOut)
		suite.Assertions.Empty(result.StdErr)
	})
}

func (suite *SSHAccTestSuite) TestRun() {
	suite.Run("should return a valid output", func() {
		ctx := context.Background()
		result, err := suite.connection.Run(ctx, "ipconfig")
		suite.Assertions.NoError(err)
		suite.Assertions.NotEmpty(result.StdOut)
		suite.Assertions.Empty(result.StdErr)
	})

	suite.Run("should return a non valid output", func() {
		ctx := context.Background()
		result, err := suite.connection.RunWithPowershell(ctx, "asdawe")
		suite.Assertions.NoError(err)
		suite.Assertions.Empty(result.StdOut)
		suite.Assertions.NotEmpty(result.StdErr)
	})

	suite.Run("should return error with context canceled", func() {
		ctx, cancel := context.WithCancel(context.Background())
		// Cancel the context immediately
		cancel()
		result, err := suite.connection.RunWithPowershell(ctx, "Start-Sleep -Seconds 20")
		suite.Assertions.ErrorContains(err, "context canceled")
		suite.Assertions.Empty(result.StdOut)
		suite.Assertions.Empty(result.StdErr)
	})
}
