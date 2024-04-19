package connection

import (
	"testing"

	"github.com/d-strobel/winrm"
	"github.com/stretchr/testify/assert"
)

func TestWinRMKerberosParams(t *testing.T) {
	t.Run("should return default kerberos configuration", func(t *testing.T) {
		KrbConf := &KerberosConfig{
			Realm:         "test.local",
			KrbConfigFile: "/home/test/krb5.conf",
		}
		winRMConf := &WinRMConfig{
			WinRMUsername: "test",
			WinRMPassword: "test",
			WinRMHost:     "winsrv",
			WinRMPort:     5555,
			WinRMUseTLS:   false,
			WinRMInsecure: true,
			WinRMKerberos: KrbConf,
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
		actualParams := winRMConf.winRMKerberosParams()
		assert.Equal(t, expectedParams, actualParams)
	})

	t.Run("should return kerberos configuration with https", func(t *testing.T) {
		KrbConf := &KerberosConfig{
			Realm:         "test.local",
			KrbConfigFile: "/home/test/krb5.conf",
		}
		winRMConf := &WinRMConfig{
			WinRMUsername: "test",
			WinRMPassword: "test",
			WinRMHost:     "winsrv",
			WinRMPort:     5555,
			WinRMUseTLS:   true,
			WinRMInsecure: true,
			WinRMKerberos: KrbConf,
		}
		expectedParams := winrm.DefaultParameters
		expectedParams.TransportDecorator = func() winrm.Transporter {
			return &winrm.ClientKerberos{
				Username: "test",
				Password: "test",
				Hostname: "winsrv",
				Realm:    "test.local",
				Port:     5555,
				Proto:    "https",
				KrbConf:  "/home/test/krb5.conf",
			}
		}
		actualParams := winRMConf.winRMKerberosParams()
		assert.Equal(t, expectedParams, actualParams)
	})
}
