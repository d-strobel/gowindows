package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

		assert.Contains(t, err.Error(), "ssh client: SSHConfig parameter 'SSHHost', 'SSHUsername' and 'SSHPassword' must be set")
	})
}

func TestKnownHostCallback(t *testing.T) {
	t.Run("IgnoreHostKey", func(t *testing.T) {
		config := &SSHConfig{
			SSHInsecureIgnoreHostKey: true,
		}
		callback, err := knownHostCallback(config)
		assert.NoError(t, err)
		assert.NotNil(t, callback)
	})
}
