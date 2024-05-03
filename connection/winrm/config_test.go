package winrm

func (suite *WinRMUnitTestSuite) TestValidate() {
	suite.T().Parallel()

	suite.Run("should be a valide WinRM config", func() {
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
				"Host + Username + Password + Port + UseTLS + Insecure + Timeout",
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					Port:     5543,
					UseTLS:   true,
					Insecure: false,
					Timeout:  0,
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
			config      *Config
		}{
			{
				"Host + Username",
				&Config{
					Host:     "test",
					Username: "test",
				},
			},
			{
				"Host + Password",
				&Config{
					Host:     "test",
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

func (suite *WinRMUnitTestSuite) TestDefaults() {
	suite.T().Parallel()

	suite.Run("should set the default values", func() {
		tcs := []struct {
			description string
			input       *Config
			expected    *Config
		}{
			{
				"minimal config",
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
				},
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					UseTLS:   false,
					Insecure: false,
					Port:     5985,
					Timeout:  0,
				},
			},
			{
				"minimal config + TLS",
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					UseTLS:   true,
				},
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					UseTLS:   true,
					Insecure: false,
					Port:     5986,
					Timeout:  0,
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			err := tc.input.defaults()
			suite.Assertions.NoError(err)
			suite.Assertions.Equal(tc.expected, tc.input)
		}
	})

	suite.Run("should not overwrite user input", func() {
		tcs := []struct {
			description string
			input       *Config
			expected    *Config
		}{
			{
				"TLS + Port",
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					UseTLS:   true,
					Port:     5555,
				},
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					UseTLS:   true,
					Insecure: false,
					Port:     5555,
					Timeout:  0,
				},
			},
			{
				"all",
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					UseTLS:   true,
					Port:     5555,
					Insecure: true,
					Timeout:  5,
				},
				&Config{
					Host:     "test",
					Username: "test",
					Password: "test",
					UseTLS:   true,
					Port:     5555,
					Insecure: true,
					Timeout:  5,
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			err := tc.input.defaults()
			suite.Assertions.NoError(err)
			suite.Assertions.Equal(tc.expected, tc.input)
		}
	})
}
