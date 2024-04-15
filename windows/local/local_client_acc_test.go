package local_test

import (
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/windows/local"
	"github.com/stretchr/testify/suite"
)

// Connection parameters for acceptance test.
const (
	testHost  = "127.0.0.1"
	username  = "vagrant"
	password  = "vagrant"
	winRMPort = 15986
	sshPort   = 1222
)

// Acceptance test suite for all local functions.
type LocalAccTestSuite struct {
	suite.Suite
	clients []local.LocalClient
}

// Setup acceptance test suite for all local functions.
// We ensure that all commands return the same output with WinRM and SSH.
func (suite *LocalAccTestSuite) SetupSuite() {
	// Setup WinRM connection
	winRMConfig := &connection.WinRMConfig{
		WinRMHost:     testHost,
		WinRMUsername: username,
		WinRMPassword: password,
		WinRMUseTLS:   true,
		WinRMInsecure: true,
		WinRMPort:     winRMPort,
	}
	winRMConn, err := connection.NewConnectionWithWinRM(winRMConfig)
	suite.Require().NoError(err)
	suite.clients = append(suite.clients, *local.NewClient(winRMConn))

	// Setup SSH connection
	sshConfig := &connection.SSHConfig{
		SSHHost:                  testHost,
		SSHPort:                  sshPort,
		SSHUsername:              username,
		SSHPassword:              password,
		SSHInsecureIgnoreHostKey: true,
	}
	sshConn, err := connection.NewConnectionWithSSH(sshConfig)
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
