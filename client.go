package client

import (
	"errors"
	"strings"
	"time"

	"github.com/masterzen/winrm"
)

const (
	// WinRM default values
	defaultWinRMPort     int           = 5986
	defaultWinRMProtocol string        = "https"
	defaultWinRMInsecure bool          = true
	defaultWinRMTimeout  time.Duration = 0
)

type Config struct {
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMProtocol string
	WinRMInsecure bool
	WinRMTimeout  time.Duration
}

func New(config *Config) (*winrm.Client, error) {

	// WinRM assert config
	if config == nil {
		return nil, errors.New("Config cannot be nil")
	}

	if config.WinRMHost == "" || config.WinRMUsername == "" || config.WinRMPassword == "" {
		return nil, errors.New("WinRMHost, WinRMUsername, and WinRMPassword must be set")
	}

	// Set default values
	if config.WinRMPort == 0 {
		config.WinRMPort = defaultWinRMPort
	}
	if config.WinRMProtocol == "" {
		config.WinRMProtocol = defaultWinRMProtocol
	}
	if config.WinRMTimeout == 0 {
		config.WinRMTimeout = defaultWinRMTimeout
	}

	// WinRM port
	winRMPort := defaultWinRMPort
	if config.WinRMPort != 0 {
		winRMPort = config.WinRMPort
	}

	// WinRM TLS
	winRMUseTLS := false
	if strings.ToLower(config.WinRMProtocol) == "https" {
		winRMUseTLS = true
	}

	// WinRM insecure
	winRMInsecure := defaultWinRMInsecure
	if config.WinRMInsecure {
		winRMInsecure = config.WinRMInsecure
	}

	// WinRM timeout
	winRMTimeout := defaultWinRMTimeout
	if config.WinRMTimeout != 0 {
		winRMTimeout = config.WinRMTimeout
	}

	// WinRM connection
	winRMEndpoint := winrm.NewEndpoint(
		config.WinRMHost,
		winRMPort,
		winRMUseTLS,
		winRMInsecure,
		nil, // CA certificate
		nil, // Client Certificate
		nil, // Client Key
		winRMTimeout,
	)

	client, err := winrm.NewClient(winRMEndpoint, config.WinRMUsername, config.WinRMPassword)
	if err != nil {
		return nil, err
	}

	return client, nil
}
