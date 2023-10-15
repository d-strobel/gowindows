package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectionErrorMessages(t *testing.T) {
	t.Run("Error - Neither WinRM nor SSH", func(t *testing.T) {
		conf := &Config{}

		_, err := New(conf)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "one of WinRMConfig and SSHConfig must be set")
	})

	t.Run("Error - Both WinRM and SSH", func(t *testing.T) {
		winRMConfig := &WinRMConfig{
			WinRMUsername: "test",
			WinRMPassword: "test",
			WinRMHost:     "test",
		}

		sshConfig := &SSHConfig{
			SSHUsername: "test",
			SSHPassword: "test",
			SSHHost:     "test",
		}

		conf := &Config{
			WinRM: winRMConfig,
			SSH:   sshConfig,
		}

		_, err := New(conf)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only one of WinRMConfig and SSHConfig must be set")
	})
}
