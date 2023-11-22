package connection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("EmptyConfig", func(t *testing.T) {
		conf := &Config{}

		_, err := New(conf)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection client: Connection object 'WinRMConfig' or 'SSHConfig' must be set")
	})

	t.Run("Both-WinRM-And-SSH", func(t *testing.T) {
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
		assert.Contains(t, err.Error(), "connection client: Connection object must only contain 'WinRMConfig' or 'SSHConfig'")
	})
}

func TestRun(t *testing.T) {
	t.Run("EmptyConfig", func(t *testing.T) {
		c := &Connection{}

		result, err := c.Run(context.Background(), "test")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "connection client: Connection object 'WinRMConfig' or 'SSHConfig' must be set")
		assert.IsType(t, CMDResult{}, result)
	})
}
