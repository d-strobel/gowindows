package connection

import (
	"testing"
)

func TestNewWinRMClient(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		// No need to set a real host since this function does not open a connection
		config := &WinRMConfig{
			WinRMHost:     "example.com",
			WinRMPort:     5986,
			WinRMUseTLS:   true,
			WinRMTimeout:  30,
			WinRMUsername: "username",
			WinRMPassword: "password",
		}

		client, err := newWinRMClient(config)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if client == nil {
			t.Error("Expected a non-nil client, but got nil")
		}
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		// Test case with missing required fields
		config := &WinRMConfig{}

		client, err := newWinRMClient(config)

		if err == nil {
			t.Error("Expected an error, but got nil")
		}

		if client != nil {
			t.Error("Expected a nil client, but got non-nil")
		}

		expectedErrorMsg := "WinRMHost, WinRMUsername, and WinRMPassword must be set"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
		}
	})
}
