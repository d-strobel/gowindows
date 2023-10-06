package gowindows

import (
	"testing"
	"time"
)

func TestNewValidConfig(t *testing.T) {
	validConfig := &Config{
		WinRMUsername: "testuser",
		WinRMPassword: "testpassword",
		WinRMHost:     "testhost",
		WinRMPort:     5986, // Provide a valid port
		WinRMProtocol: "https",
		WinRMInsecure: true,
		WinRMTimeout:  10 * time.Second, // Provide a valid timeout
	}

	client, err := NewClient(validConfig)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if client == nil {
		t.Error("Expected a valid client, but got nil")
	}
}

func TestNewInvalidConfig(t *testing.T) {
	invalidConfig := &Config{
		// Missing required fields
	}

	client, err := NewClient(invalidConfig)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	if client != nil {
		t.Error("Expected a nil client, but got a valid client")
	}
}

func TestNewDefaultValues(t *testing.T) {
	defaultConfig := &Config{
		WinRMUsername: "testuser",
		WinRMPassword: "testpassword",
		WinRMHost:     "testhost",
	}

	client, err := NewClient(defaultConfig)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if client == nil {
		t.Error("Expected a valid client, but got nil")
	}
}
