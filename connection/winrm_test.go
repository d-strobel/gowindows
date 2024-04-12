package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all WinRM functions.
type ConnectionWinRMUnitTestSuite struct {
	suite.Suite
}

func TestConnectionWinRMUnitTestSuite(t *testing.T) {
	suite.Run(t, &ConnectionWinRMUnitTestSuite{})
}

func (suite *ConnectionWinRMUnitTestSuite) TestValidate() {

	suite.Run("should be a valide WinRM config", func() {
		tcs := []struct {
			description string
			config      *WinRMConfig
		}{
			{
				"Host + Username + Password",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
				},
			},
			{
				"Host + Username + Password + Port + UseTLS + Insecure + Timeout",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMPort:     5543,
					WinRMUseTLS:   true,
					WinRMInsecure: false,
					WinRMTimeout:  0,
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			err := tc.config.Validate()
			suite.Assertions.NoError(err)
		}
	})

	suite.Run("should return an Error", func() {
		tcs := []struct {
			description string
			config      *WinRMConfig
		}{
			{
				"Host + Username",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
				},
			},
			{
				"Host + Password",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMPassword: "test",
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			err := tc.config.Validate()
			suite.Assertions.Error(err)
		}
	})
}
