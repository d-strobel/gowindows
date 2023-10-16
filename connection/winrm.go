package connection

import (
	"errors"
	"time"

	"github.com/masterzen/winrm"
)

type WinRMConfig struct {
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMUseTLS   bool
	WinRMInsecure bool
	WinRMTimeout  time.Duration
	WinRMKerberos *KerberosConfig
}

func newWinRMClient(config *WinRMConfig) (*winrm.Client, error) {

	// Assert
	if config.WinRMHost == "" || config.WinRMUsername == "" || config.WinRMPassword == "" {
		return nil, errors.New("WinRMHost, WinRMUsername, and WinRMPassword must be set")
	}

	// Set default values
	if config.WinRMPort == 0 {
		config.WinRMPort = defaultWinRMPort
	}
	if !config.WinRMUseTLS {
		config.WinRMUseTLS = defaultWinRMUseTLS
	}
	if config.WinRMTimeout == 0 {
		config.WinRMTimeout = defaultWinRMTimeout
	}

	winRMInsecure := defaultWinRMInsecure
	if config.WinRMInsecure {
		winRMInsecure = config.WinRMInsecure
	}

	// WinRM connection
	winRMEndpoint := winrm.NewEndpoint(
		config.WinRMHost,
		config.WinRMPort,
		config.WinRMUseTLS,
		winRMInsecure,
		nil, // CA certificate
		nil, // Client Certificate
		nil, // Client Key
		config.WinRMTimeout,
	)

	// Kerberos transport
	if config.WinRMKerberos != nil {
		params := winRMKerberosParams(config)
		return winrm.NewClientWithParameters(winRMEndpoint, config.WinRMUsername, config.WinRMPassword, params)
	}

	return winrm.NewClient(winRMEndpoint, config.WinRMUsername, config.WinRMPassword)
}
