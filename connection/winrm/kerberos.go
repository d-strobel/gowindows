package winrm

import (
	"errors"

	"github.com/d-strobel/winrm"
)

// Default values for Kerberos configuration.
const (
	defaultKerberosProtocol string = "https"
)

// KerberosConfig represents the configuration details for Kerberos authentication.
type KerberosConfig struct {
	Realm         string
	KrbConfigFile string
	Protocol      string
}

// validate validates the Kerberos configuration parameters.
func (krbConfig *KerberosConfig) validate() error {
	if krbConfig.Realm == "" || krbConfig.KrbConfigFile == "" {
		return errors.New("winrm: KerberosConfig parameter 'Realm' and 'KrbConfigFile' must be set")
	}

	if krbConfig.Protocol != "" {
		if krbConfig.Protocol != "http" && krbConfig.Protocol != "https" {
			return errors.New("winrm: KerberosConfig parameter 'Protocol' must be one of 'http' or 'https'")
		}
	}

	return nil
}

// defaults sets the default values for the Kerberos configuration.
func (krbConfig *KerberosConfig) defaults() error {
	if krbConfig.Protocol == "" {
		krbConfig.Protocol = defaultKerberosProtocol
	}

	return nil
}

// kerberosParams returns the necessary parameters
// to pass into the Kerberos WinRM connection.
func (krbConfig *KerberosConfig) kerberosParams(config *Config) *winrm.Parameters {
	// Init default parameters
	params := winrm.NewParameters("PT60S", "en-US", 153600)

	// Configure kerberos transporter
	params.TransportDecorator = func() winrm.Transporter {
		return &winrm.ClientKerberos{
			Username: config.Username,
			Password: config.Password,
			Hostname: config.Host,
			Port:     config.Port,
			Proto:    krbConfig.Protocol,
			Realm:    krbConfig.Realm,
			KrbConf:  krbConfig.KrbConfigFile,
		}
	}

	return params
}
