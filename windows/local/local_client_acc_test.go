package local_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/d-strobel/gowindows/connection/ssh"
	"github.com/d-strobel/gowindows/connection/winrm"
	"github.com/d-strobel/gowindows/windows/local"
	"github.com/stretchr/testify/suite"
)

// Acceptance test suite for all local functions.
type LocalAccTestSuite struct {
	suite.Suite

	// Fixtures
	host      string
	username  string
	password  string
	winRMPort int
	sshPort   int
	clients   []local.Client
}

// SetupSuite setups the acceptance test suite for all local functions.
// We ensure that all commands return the same output with WinRM and SSH.
func (suite *LocalAccTestSuite) SetupSuite() {
	var err error

	// Load environment variables
	suite.host = os.Getenv("GOWINDOWS_TEST_HOST")
	suite.Require().NotEmpty(suite.host, "Environment variable not set: GOWINDOWS_TEST_HOST")

	suite.username = os.Getenv("GOWINDOWS_TEST_USERNAME")
	suite.Require().NotEmpty(suite.username, "Environment variable not set: GOWINDOWS_TEST_USERNAME")

	suite.password = os.Getenv("GOWINDOWS_TEST_PASSWORD")
	suite.Require().NotEmpty(suite.password, "Environment variable not set: GOWINDOWS_TEST_PASSWORD")

	suite.winRMPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_WINRM_HTTP_PORT"))
	suite.Require().NoError(err)

	suite.sshPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_SSH_PORT"))
	suite.Require().NoError(err)

	// Setup WinRM connection
	winRMConfig := &winrm.Config{
		Host:     suite.host,
		Username: suite.username,
		Password: suite.password,
		UseTLS:   false,
		Insecure: true,
		Port:     suite.winRMPort,
	}
	winRMConn, err := winrm.NewConnection(winRMConfig)
	suite.Require().NoError(err)
	suite.clients = append(suite.clients, *local.NewClient(winRMConn))

	// Setup SSH connection
	sshConfig := &ssh.Config{
		Host:                  suite.host,
		Username:              suite.username,
		Password:              suite.password,
		Port:                  suite.sshPort,
		InsecureIgnoreHostKey: true,
	}
	sshConn, err := ssh.NewConnection(sshConfig)
	suite.Require().NoError(err)
	suite.clients = append(suite.clients, *local.NewClient(sshConn))
}

func (suite *LocalAccTestSuite) TearDownSuite() {
	// Close connections
	for _, c := range suite.clients {
		c.Connection.Close()
	}
}

func TestLocalAccTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, &LocalAccTestSuite{})
}
