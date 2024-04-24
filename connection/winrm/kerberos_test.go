package winrm

import (
	"github.com/d-strobel/winrm"
)

func (suite *WinRMUnitTestSuite) TestKerberosValidate() {
	suite.Run("should be a valid configuration", func() {
		tcs := []struct {
			description string
			config      *KerberosConfig
		}{
			{
				"Realm + KrbConfigFile + Protocol-https",
				&KerberosConfig{
					Realm:         "test.local",
					KrbConfigFile: "/home/test/krb5.conf",
					Protocol:      "https",
				},
			},
			{
				"Realm + KrbConfigFile + Protocol-https",
				&KerberosConfig{
					Realm:         "test.local",
					KrbConfigFile: "/home/test/krb5.conf",
					Protocol:      "http",
				},
			},
			{
				"Realm + KrbConfigFile",
				&KerberosConfig{
					Realm:         "test.local",
					KrbConfigFile: "/home/test/krb5.conf",
				},
			},
		}

		for _, tc := range tcs {
			suite.T().Logf("test case: %s", tc.description)
			err := tc.config.validate()
			suite.Assertions.NoError(err)
		}
	})

	suite.Run("should not be a valid configuration", func() {
		tcs := []struct {
			description string
			config      *KerberosConfig
		}{
			{
				"Realm + Protocol",
				&KerberosConfig{
					Realm:    "test.local",
					Protocol: "https",
				},
			},
			{
				"KrbConfigFile + Protocol",
				&KerberosConfig{
					KrbConfigFile: "/home/test/krb5.conf",
					Protocol:      "http",
				},
			},
			{
				"Protocol-http",
				&KerberosConfig{
					Protocol: "http",
				},
			},
			{
				"Protocol-https",
				&KerberosConfig{
					Protocol: "https",
				},
			},
			{
				"Realm + Protocol + Protocol-httt",
				&KerberosConfig{
					Realm:         "test.local",
					KrbConfigFile: "/home/test/krb5.conf",
					Protocol:      "httt",
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

func (suite *WinRMUnitTestSuite) TestKerberosParams() {
	suite.Run("should return a valid kerberos configuration", func() {
		KrbConf := &KerberosConfig{
			Realm:         "test.local",
			KrbConfigFile: "/home/test/krb5.conf",
			Protocol:      "http",
		}
		winRMConf := &Config{
			Username: "test",
			Password: "test",
			Host:     "winsrv",
			Port:     5555,
			UseTLS:   false,
			Insecure: true,
		}
		expectedParams := winrm.DefaultParameters
		expectedParams.TransportDecorator = func() winrm.Transporter {
			return &winrm.ClientKerberos{
				Username: "test",
				Password: "test",
				Hostname: "winsrv",
				Realm:    "test.local",
				Port:     5555,
				Proto:    "http",
				KrbConf:  "/home/test/krb5.conf",
			}
		}
		actualParams := KrbConf.kerberosParams(winRMConf)
		suite.Assertions.Equal(expectedParams, actualParams)
	})
}
