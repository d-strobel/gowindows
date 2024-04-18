package connection_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Init acceptance test suite for SSH
type ConnectionAccTestSuite struct {
	suite.Suite

	// Fixtures
	host              string
	username          string
	password          string
	sshPort           int
	sshKeyPathED25519 string
	sshKeyED25519     string
	sshKeyPathRSA     string
	sshKeyRSA         string
	winRMPort         int
}

// SetupSuite setups all neccessary fixtures for running the connection tests.
func (suite *ConnectionAccTestSuite) SetupSuite() {
	var err error

	// Load environment variables
	suite.host = os.Getenv("GOWINDOWS_TEST_HOST")
	suite.Require().NotEmpty(suite.host, "Environment variable not set: GOWINDOWS_TEST_HOST")

	suite.username = os.Getenv("GOWINDOWS_TEST_USERNAME")
	suite.Require().NotEmpty(suite.username, "Environment variable not set: GOWINDOWS_TEST_USERNAME")

	suite.password = os.Getenv("GOWINDOWS_TEST_PASSWORD")
	suite.Require().NotEmpty(suite.password, "Environment variable not set: GOWINDOWS_TEST_PASSWORD")

	suite.sshPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_SSH_PORT"))
	suite.Require().NoError(err)

	suite.sshKeyPathED25519 = os.Getenv("GOWINDOWS_TEST_SSH_PRIVATE_KEY_ED25519_PATH")
	suite.Require().NotEmpty(suite.sshKeyPathED25519, "Environment variable not set: GOWINDOWS_TEST_SSH_PRIVATE_KEY_ED25519_PATH")

	suite.sshKeyPathRSA = os.Getenv("GOWINDOWS_TEST_SSH_PRIVATE_KEY_RSA_PATH")
	suite.Require().NotEmpty(suite.sshKeyPathRSA, "Environment variable not set: GOWINDOWS_TEST_SSH_PRIVATE_KEY_RSA_PATH")

	privateKeyRSA, err := os.ReadFile(suite.sshKeyPathRSA)
	suite.Require().NoError(err)
	suite.sshKeyRSA = string(privateKeyRSA)

	privateKeyED25519, err := os.ReadFile(suite.sshKeyPathED25519)
	suite.Require().NoError(err)
	suite.sshKeyED25519 = string(privateKeyED25519)

	suite.winRMPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_WINRM_HTTP_PORT"))
	suite.Require().NoError(err)
}

// TestConnectionAccTestSuite runs the acceptance test suite for the connection package.
// It will be skipped if the short flag is set.
func TestConnectionAccTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, &ConnectionAccTestSuite{})
}
