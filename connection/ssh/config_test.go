package ssh

import (
	"fmt"
)

func (suite *SSHUnitTestSuite) TestValidate() {
	suite.Run("should be a valide SSH config", func() {
		tcs := []struct {
			description string
			config      *Config
		}{
			{
				"Host + Username + Password",
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
				},
			},
			{
				"Host + Username + PrivateKey",
				&Config{
					Host:       "test",
					Username:   "test",
					PrivateKey: "test",
				},
			},
			{
				"Host + Username + PrivateKeyPath",
				&Config{
					Host:           "test",
					Username:       "test",
					PrivateKeyPath: "/test/test",
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			err := tc.config.validate()
			suite.Assertions.NoError(err)
		}
	})

	suite.Run("should return an error", func() {
		tcs := []struct {
			description string
			config      *Config
		}{
			{
				"Host",
				&Config{
					Host: "test",
				},
			},
			{
				"Host + Username",
				&Config{
					Host:     "test",
					Username: "test",
				},
			},
			{
				"Host + PrivateKey",
				&Config{
					Host:       "test",
					PrivateKey: "test",
				},
			},
			{
				"Username + Password",
				&Config{
					Username: "test",
					Password: "test",
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

func (suite *SSHUnitTestSuite) TestDefaults() {

	suite.Run("should set the default values", func() {
		input := &Config{
			Host:     "test",
			Username: "test",
			Password: "test",
		}
		expected := &Config{
			Host:           "test",
			Username:       "test",
			Password:       "test",
			Port:           22,
			KnownHostsPath: fmt.Sprintf("%s/%s", suite.currentUserHomeDir, defaultKnownHostsPath),
		}
		err := input.defaults()
		suite.Assertions.NoError(err)
		suite.Assertions.EqualValues(input, expected)
	})

	suite.Run("should not overwrite user input", func() {
		input := &Config{
			Host:           "test",
			Username:       "test",
			Password:       "test",
			Port:           2222,
			KnownHostsPath: "/home/test/.ssh/known_hosts",
		}
		expected := &Config{
			Host:           "test",
			Username:       "test",
			Password:       "test",
			Port:           2222,
			KnownHostsPath: "/home/test/.ssh/known_hosts",
		}
		err := input.defaults()
		suite.Assertions.NoError(err)
		suite.Assertions.EqualValues(input, expected)
	})
}
