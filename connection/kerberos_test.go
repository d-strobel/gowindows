package connection

import (
	"testing"

	"github.com/masterzen/winrm"
	"github.com/stretchr/testify/assert"
)

func TestWinRMKerberosParams(t *testing.T) {
	// Create a sample WinRM configuration
	winRMConfig := &WinRMConfig{
		WinRMUsername: "testUser",
		WinRMPassword: "testPassword",
		WinRMHost:     "testHost",
		WinRMKerberos: &KerberosConfig{
			Realm:         "testRealm",
			KrbConfigFile: "/path/to/krb5.conf",
		},
		WinRMUseTLS: false,
		WinRMPort:   5985,
	}

	// Call the function to get the parameters
	params := winRMKerberosParams(winRMConfig)

	// Check that the parameters are set as expected
	assert.NotNil(t, params)
	assert.NotNil(t, params.TransportDecorator)
	assert.Equal(t, winRMConfig.WinRMUsername, params.TransportDecorator().(*winrm.ClientKerberos).Username)
	assert.Equal(t, winRMConfig.WinRMPassword, params.TransportDecorator().(*winrm.ClientKerberos).Password)
	assert.Equal(t, winRMConfig.WinRMHost, params.TransportDecorator().(*winrm.ClientKerberos).Hostname)
	assert.Equal(t, winRMConfig.WinRMKerberos.Realm, params.TransportDecorator().(*winrm.ClientKerberos).Realm)
	assert.Equal(t, winRMConfig.WinRMPort, params.TransportDecorator().(*winrm.ClientKerberos).Port)
	assert.Equal(t, "http", params.TransportDecorator().(*winrm.ClientKerberos).Proto)
	assert.Equal(t, winRMConfig.WinRMKerberos.KrbConfigFile, params.TransportDecorator().(*winrm.ClientKerberos).KrbConf)
}

func TestWinRMKerberosParamsWithTLS(t *testing.T) {
	// Create a sample WinRM configuration with WinRMUseTLS set to true
	winRMConfig := &WinRMConfig{
		WinRMUsername: "testUser",
		WinRMPassword: "testPassword",
		WinRMHost:     "testHost",
		WinRMKerberos: &KerberosConfig{
			Realm:         "testRealm",
			KrbConfigFile: "/path/to/krb5.conf",
		},
		WinRMUseTLS: true, // Set WinRMUseTLS to true
		WinRMPort:   5985,
	}

	// Call the function to get the parameters
	params := winRMKerberosParams(winRMConfig)

	// Check that the protocol is set to "https" when WinRMUseTLS is true
	assert.Equal(t, "https", params.TransportDecorator().(*winrm.ClientKerberos).Proto)
}
