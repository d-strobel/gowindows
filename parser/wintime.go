package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// WinTime is a custom time type that embeds the time.Time type. It is designed to handle
// the unmarshalling of dotnet JSON datetime strings in the format "\"/Date(timestamp)/\""
// when used as a field in a struct that is being unmarshalled from JSON.
type WinTime struct {
	time.Time
}

// dotNetTimeStringPattern represents the regex pattern for validating a dotnet JSON datetime string.
const dotNetTimeStringPattern = `^"\\/Date\(\d+\)\\/"$`

// timestampPattern represents the regex pattern for extracting the timestamp from a dotnet JSON datetime string.
const timestampPattern = `\d+`

// UnmarshalJSON implements the json.Unmarshaler interface for the WinTime type.
// It parses a JSON-encoded dotnet JSON datetime string and converts it into a WinTime object.
// The input byte slice is expected to represent a dotnet JSON datetime string in the format "\"/Date(timestamp)/\"".
func (t *WinTime) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package
	if string(b) == "null" || string(b) == `""` {
		return nil
	}

	// Check for valid dotnet timestring
	re, err := regexp.Compile(dotNetTimeStringPattern)
	if err != nil {
		return fmt.Errorf("WinTime.UnmarshalJSON: %s", err)
	}

	if !re.Match(b) {
		return fmt.Errorf("WinTime.UnmarshalJSON: input string is not a dotnet JSON datetime: %s", string(b))
	}

	// Extract timestamp
	re, err = regexp.Compile(timestampPattern)
	if err != nil {
		return fmt.Errorf("WinTime.UnmarshalJSON: %s", err)
	}
	timestamp := re.Find(b)

	// Convert to seconds
	i, err := strconv.Atoi(string(timestamp))
	if err != nil {
		return fmt.Errorf("WinTime.UnmarshalJSON: %s", err)
	}
	seconds := int64(i / 1000)

	// Unmarshal unix time into WinTime object
	unixTime := time.Unix(seconds, 0).UTC()
	*t = WinTime{unixTime}

	return nil
}
