package parsing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPwshTimespanString(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		description    string
		inputDuration  string
		expectedString string
	}{
		{
			"duration of 0",
			"0h",
			"$(New-TimeSpan -Days 0 -Hours 0 -Minutes 0 -Seconds 0)",
		},
		{
			"duration with 2 days 30 minutes and 15 seconds",
			"48h30m15s",
			"$(New-TimeSpan -Days 2 -Hours 0 -Minutes 30 -Seconds 15)",
		},
		{
			"duration with 1 day 2 hours 30 minutes and 15 seconds",
			"26h30m15s",
			"$(New-TimeSpan -Days 1 -Hours 2 -Minutes 30 -Seconds 15)",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			d, err := time.ParseDuration(tc.inputDuration)
			assert.NoError(t, err)
			actualString := PwshTimespanString(d)
			assert.Equal(t, tc.expectedString, actualString)
		})
	}
}
