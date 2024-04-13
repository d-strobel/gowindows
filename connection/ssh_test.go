package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Unit test suite for all SSH functions.
type ConnectionSSHUnitTestSuite struct {
	suite.Suite
}

func TestConnectionSSHUnitTestSuite(t *testing.T) {
	suite.Run(t, &ConnectionSSHUnitTestSuite{})
}

func (suite *ConnectionSSHUnitTestSuite) TestValidate() {

	suite.Run("should be a valide SSH config", func() {
		tcs := []struct {
			description string
			config      *SSHConfig
		}{
			{
				"Host + Username + Password",
				&SSHConfig{
					SSHHost:     "test",
					SSHUsername: "test",
					SSHPassword: "test",
				},
			},
			{
				"Host + Username + PrivateKey",
				&SSHConfig{
					SSHHost:       "test",
					SSHUsername:   "test",
					SSHPrivateKey: "test",
				},
			},
			{
				"Host + Username + PrivateKeyPath",
				&SSHConfig{
					SSHHost:           "test",
					SSHUsername:       "test",
					SSHPrivateKeyPath: "/test/test",
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			err := tc.config.Validate()
			suite.Assertions.NoError(err)
		}
	})

	suite.Run("should return an error", func() {
		tcs := []struct {
			description string
			config      *SSHConfig
		}{
			{
				"Host",
				&SSHConfig{
					SSHHost: "test",
				},
			},
			{
				"Host + Username",
				&SSHConfig{
					SSHHost:     "test",
					SSHUsername: "test",
				},
			},
			{
				"Host + PrivateKey",
				&SSHConfig{
					SSHHost:       "test",
					SSHPrivateKey: "test",
				},
			},
			{
				"Username + Password",
				&SSHConfig{
					SSHUsername: "test",
					SSHPassword: "test",
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