package parsing

import (
	"fmt"
	"time"
)

// PwshTimespanString returns a string representation of a time.Duration that can be used in a PowerShell command.
// The returned string is in the format "$(New-TimeSpan -Days <days> -Hours <hours> -Minutes <minutes> -Seconds <seconds>)".
func PwshTimespanString(d time.Duration) string {
	return fmt.Sprintf(
		"$(New-TimeSpan -Days %d -Hours %d -Minutes %d -Seconds %d)",
		int32(d.Hours())/24,
		int32(d.Hours())%24,
		int32(d.Minutes())%60,
		int32(d.Seconds())%60,
	)
}
