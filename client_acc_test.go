package gowindows_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/connection/ssh"
	"github.com/d-strobel/gowindows/connection/winrm"
	"github.com/stretchr/testify/suite"
)

type GowindowsAccTestSuite struct {
	suite.Suite

	// Fixtures
	host     string
	username string
	password string
	sshPort  int
	httpPort int
}

// SetupSuite setups all neccessary fixtures for running the gowindows tests.
func (suite *GowindowsAccTestSuite) SetupSuite() {
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

	suite.httpPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_WINRM_HTTP_PORT"))
	suite.Require().NoError(err)
}

// TestGowindowsAccTestSuite runs the acceptance test suite for the gowindwos package.
// It will be skipped if the short flag is set.
func TestGowindowsAccTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, &GowindowsAccTestSuite{})
}

func (suite *GowindowsAccTestSuite) TestNewClient() {
	// Test subpackages to ensure the gowindows client works as expected.
	suite.Run("should return local users with an ssh connection", func() {
		sshConfig := &ssh.Config{
			Host:     suite.host,
			Port:     suite.sshPort,
			Username: suite.username,
			Password: suite.password,
		}

		conn, err := ssh.NewConnection(sshConfig)
		suite.Require().NoError(err)
		defer conn.Close()

		client := gowindows.NewClient(conn)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err = client.Local.UserList(ctx)
		suite.NoError(err)
	})
	suite.Run("should return local users with a winrm connection", func() {
		winrmConfig := &winrm.Config{
			Host:     suite.host,
			Port:     suite.httpPort,
			Username: suite.username,
			Password: suite.password,
		}

		conn, err := winrm.NewConnection(winrmConfig)
		suite.Require().NoError(err)
		defer conn.Close()

		client := gowindows.NewClient(conn)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err = client.Local.UserList(ctx)
		suite.NoError(err)
	})
}
