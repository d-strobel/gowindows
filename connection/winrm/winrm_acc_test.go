package winrm_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/d-strobel/gowindows/connection/winrm"
	"github.com/stretchr/testify/suite"
)

// Init acceptance test suite for SSH
type WinRMAccTestSuite struct {
	suite.Suite

	// Fixtures
	host        string
	username    string
	password    string
	httpPort    int
	httpsPort   int
	adHost      string
	adUsername  string
	adPassword  string
	adHttpPort  int
	adHttpsPort int
}

// SetupSuite setups all neccessary fixtures for running the winrm tests.
func (suite *WinRMAccTestSuite) SetupSuite() {
	var err error

	// Load environment variables
	suite.host = os.Getenv("GOWINDOWS_TEST_HOST")
	suite.Require().NotEmpty(suite.host, "Environment variable not set: GOWINDOWS_TEST_HOST")

	suite.username = os.Getenv("GOWINDOWS_TEST_USERNAME")
	suite.Require().NotEmpty(suite.username, "Environment variable not set: GOWINDOWS_TEST_USERNAME")

	suite.password = os.Getenv("GOWINDOWS_TEST_PASSWORD")
	suite.Require().NotEmpty(suite.password, "Environment variable not set: GOWINDOWS_TEST_PASSWORD")

	suite.httpPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_WINRM_HTTP_PORT"))
	suite.Require().NoError(err)

	suite.httpsPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_WINRM_HTTPS_PORT"))
	suite.Require().NoError(err)

	suite.adHost = os.Getenv("GOWINDOWS_TEST_AD_HOST")
	suite.Require().NotEmpty(suite.host, "Environment variable not set: GOWINDOWS_TEST_AD_HOST")

	suite.adUsername = os.Getenv("GOWINDOWS_TEST_AD_USERNAME")
	suite.Require().NotEmpty(suite.username, "Environment variable not set: GOWINDOWS_TEST_AD_USERNAME")

	suite.adPassword = os.Getenv("GOWINDOWS_TEST_AD_PASSWORD")
	suite.Require().NotEmpty(suite.password, "Environment variable not set: GOWINDOWS_TEST_AD_PASSWORD")

	suite.adHttpPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_AD_WINRM_HTTP_PORT"))
	suite.Require().NoError(err)

	suite.adHttpsPort, err = strconv.Atoi(os.Getenv("GOWINDOWS_TEST_AD_WINRM_HTTPS_PORT"))
	suite.Require().NoError(err)
}

// TestWinRMAccTestSuite runs the acceptance test suite for the winrm package.
// It will be skipped if the short flag is set.
func TestWinRMAccTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, &WinRMAccTestSuite{})
}

// TestNewConnection tests the NewConnection method.
// It needs to execute the Run() method to test the connection.
// This is neccessary because the winrm library is designed to
// open a connection only when a command is run.
func (suite *WinRMAccTestSuite) TestNewConnection() {
	suite.Run("should create a valid connection", func() {
		tcs := []struct {
			description string
			config      *winrm.Config
		}{
			{
				"Host + Username + Password + Port + http + insecure",
				&winrm.Config{
					Host:     suite.host,
					Username: suite.username,
					Password: suite.password,
					Port:     suite.httpPort,
					UseTLS:   false,
					Insecure: true,
				},
			},
			{
				"Host + Username + Password + Port + https + insecure",
				&winrm.Config{
					Host:     suite.host,
					Username: suite.username,
					Password: suite.password,
					Port:     suite.httpsPort,
					UseTLS:   true,
					Insecure: true,
				},
			},
			{
				"Domain Host + Username + Password + Port + http + insecure",
				&winrm.Config{
					Host:     suite.adHost,
					Username: suite.adUsername,
					Password: suite.adPassword,
					Port:     suite.adHttpPort,
					UseTLS:   false,
					Insecure: true,
				},
			},
			{
				"Domain Host + Username + Password + Port + https + insecure",
				&winrm.Config{
					Host:     suite.adHost,
					Username: suite.adUsername,
					Password: suite.adPassword,
					Port:     suite.adHttpsPort,
					UseTLS:   true,
					Insecure: true,
				},
			},
		}
		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			conn, err := winrm.NewConnection(tc.config)
			suite.Assertions.NoError(err)
			defer conn.Close()
			// We run a command to test the connection.
			_, err = conn.Run(context.Background(), "ipconfig")
			// We only need to know that the connection has no errors.
			suite.Assertions.NoError(err)
		}
	})

	suite.Run("should fail to create a valid connection", func() {
		tcs := []struct {
			description string
			config      *winrm.Config
		}{
			{
				"UseTLS + http",
				&winrm.Config{
					Host:     suite.host,
					Username: suite.username,
					Password: suite.password,
					Port:     suite.httpPort,
					UseTLS:   true,
					Insecure: true,
				},
			},
			{
				"Without TLS + https",
				&winrm.Config{
					Host:     suite.host,
					Username: suite.username,
					Password: suite.password,
					Port:     suite.httpsPort,
					UseTLS:   false,
					Insecure: true,
				},
			},
			{
				"cannot validate certificate",
				&winrm.Config{
					Host:     suite.host,
					Username: suite.username,
					Password: suite.password,
					Port:     suite.httpsPort,
					UseTLS:   true,
					Insecure: false,
				},
			},
		}
		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			conn, err := winrm.NewConnection(tc.config)
			suite.Assertions.NoError(err)
			defer conn.Close()
			// We run a command to test the connection.
			_, err = conn.Run(context.Background(), "ipconfig")
			// We only need to know that the connection has errors.
			suite.Assertions.Error(err)
		}
	})
}

// TestRun tests the Run method.
func (suite *WinRMAccTestSuite) TestRun() {
	suite.Run("should run the command successfully and return stdout", func() {
		winRMConfig := &winrm.Config{
			Host:     suite.host,
			Username: suite.username,
			Password: suite.password,
			Port:     suite.httpPort,
			UseTLS:   false,
			Insecure: true,
		}
		conn, err := winrm.NewConnection(winRMConfig)
		suite.Assertions.NoError(err)
		defer conn.Close()
		result, err := conn.Run(context.Background(), "ipconfig")
		suite.Assertions.NoError(err)
		suite.Assertions.NotEmpty(result.StdOut)
		suite.Assertions.Empty(result.StdErr)
	})

	suite.Run("should get an stderr result", func() {
		winRMConfig := &winrm.Config{
			Host:     suite.host,
			Username: suite.username,
			Password: suite.password,
			Port:     suite.httpPort,
			UseTLS:   false,
			Insecure: true,
		}
		conn, err := winrm.NewConnection(winRMConfig)
		suite.Assertions.NoError(err)
		defer conn.Close()
		result, err := conn.Run(context.Background(), "iipconfig")
		suite.Assertions.NoError(err)
		suite.Assertions.Empty(result.StdOut)
		suite.Assertions.NotEmpty(result.StdErr)
	})
}
