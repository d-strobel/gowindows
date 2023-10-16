package connection

import (
	"testing"
)

func TestNewSSHClient(t *testing.T) {
	t.Run("InvalidConfig", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		expectedErrorMsg := "SSHHost, SSHUsername, and SSHPassword must be set"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
		}
	})
}
