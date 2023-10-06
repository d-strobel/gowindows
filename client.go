package gowindows

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

type Client struct {
	WinRM *winrm.Client
}

type Config struct {
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMProtocol string
	WinRMInsecure bool
	WinRMTimeout  time.Duration
}

func NewClient(config *Config) (*Client, error) {

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

	winRMUseTLS := false
	if strings.ToLower(config.WinRMProtocol) == "https" {
		winRMUseTLS = true
	}

	winRMInsecure := defaultWinRMInsecure
	if config.WinRMInsecure {
		winRMInsecure = config.WinRMInsecure
	}

	// WinRM connection
	winRMEndpoint := winrm.NewEndpoint(
		config.WinRMHost,
		config.WinRMPort,
		winRMUseTLS,
		winRMInsecure,
		nil, // CA certificate
		nil, // Client Certificate
		nil, // Client Key
		config.WinRMTimeout,
	)

	winRMClient, err := winrm.NewClient(winRMEndpoint, config.WinRMUsername, config.WinRMPassword)
	if err != nil {
		return nil, err
	}

	c := &Client{
		WinRM: winRMClient,
	}

	return c, nil
}
