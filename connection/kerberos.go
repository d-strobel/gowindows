package connection

import (
	"github.com/d-strobel/winrm"
)

// KerberosConfig represents the configuration details for Kerberos authentication.
type KerberosConfig struct {
	Realm         string
	KrbConfigFile string
}

// winRMKerberosParams returns the necessary parameters
// to pass into the Kerberos WinRM connection.
func (config *WinRMConfig) winRMKerberosParams() *winrm.Parameters {

	// Init default parameters
	params := winrm.DefaultParameters

	// When TLS is enabled use https, otherwise use http
	kerberosProtocol := "http"
	if config.WinRMUseTLS {
		kerberosProtocol = "https"
	}

	// Configure kerberos transporter
	params.TransportDecorator = func() winrm.Transporter {
		return &winrm.ClientKerberos{
			Username: config.WinRMUsername,
			Password: config.WinRMPassword,
			Hostname: config.WinRMHost,
			Realm:    config.WinRMKerberos.Realm,
			Port:     config.WinRMPort,
			Proto:    kerberosProtocol,
			KrbConf:  config.WinRMKerberos.KrbConfigFile,
		}
	}

	return params
}
