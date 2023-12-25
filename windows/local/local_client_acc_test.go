package local_test

import (
	"testing"

	"github.com/d-strobel/gowindows/connection"
	"github.com/d-strobel/gowindows/parser"
	"github.com/d-strobel/gowindows/windows/local"
	"github.com/stretchr/testify/suite"
)

// Connection parameters for acceptance test
const (
	testHost  = "127.0.0.1"
	username  = "vagrant"
	password  = "vagrant"
	winRMPort = 15986
	sshPort   = 1222
)

// Acceptance test suite for all local functions
type LocalAccTestSuite struct {
	suite.Suite
	clients []local.LocalClient
}

// Setup acceptance test suite for all local functions
func (suite *LocalAccTestSuite) SetupSuite() {
	parser := parser.NewParser()

	// Connection configs
	testConfigs := []connection.Config{
		// WinRM
		{
			WinRM: &connection.WinRMConfig{
				WinRMHost:     testHost,
				WinRMUsername: username,
				WinRMPassword: password,
				WinRMUseTLS:   true,
				WinRMInsecure: true,
				WinRMPort:     winRMPort,
			},
		},
		// SSH
		{
			SSH: &connection.SSHConfig{
				SSHHost:                  testHost,
				SSHPort:                  sshPort,
				SSHUsername:              username,
				SSHPassword:              password,
				SSHInsecureIgnoreHostKey: true,
			},
		},
	}

	// Setup clients for tests
	for _, conf := range testConfigs {
		conn, err := connection.NewConnection(&conf)
		suite.Require().NoError(err)

		suite.clients = append(suite.clients, *local.NewLocalClient(conn, parser))
	}
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
