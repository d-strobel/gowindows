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
			err := tc.config.validate()
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
			err := tc.config.validate()
			suite.Assertions.Error(err)
		}
	})
}

func (suite *ConnectionWinRMUnitTestSuite) TestDefaults() {

	suite.Run("should set the default values", func() {
		tcs := []struct {
			description string
			input       *WinRMConfig
			expected    *WinRMConfig
		}{
			{
				"minimal config",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
				},
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMUseTLS:   false,
					WinRMInsecure: false,
					WinRMPort:     5985,
					WinRMTimeout:  0,
				},
			},
			{
				"minimal config + TLS",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMUseTLS:   true,
				},
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMUseTLS:   true,
					WinRMInsecure: false,
					WinRMPort:     5986,
					WinRMTimeout:  0,
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			tc.input.defaults()
			suite.Assertions.Equal(tc.expected, tc.input)
		}
	})

	suite.Run("should not overwrite user input", func() {
		tcs := []struct {
			description string
			input       *WinRMConfig
			expected    *WinRMConfig
		}{
			{
				"TLS + Port",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMUseTLS:   true,
					WinRMPort:     5555,
				},
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMUseTLS:   true,
					WinRMInsecure: false,
					WinRMPort:     5555,
					WinRMTimeout:  0,
				},
			},
			{
				"all",
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMUseTLS:   true,
					WinRMPort:     5555,
					WinRMInsecure: true,
					WinRMTimeout:  5,
				},
				&WinRMConfig{
					WinRMHost:     "test",
					WinRMUsername: "test",
					WinRMPassword: "test",
					WinRMUseTLS:   true,
					WinRMPort:     5555,
					WinRMInsecure: true,
					WinRMTimeout:  5,
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			tc.input.defaults()
			suite.Assertions.Equal(tc.expected, tc.input)
		}
	})
}
