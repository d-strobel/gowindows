package parsing

import (
	"encoding/json"
	"time"
)

// CimTimeDuration is a custom time type that embeds the time.Duration type.
// It is designed to handle the unmarshalling of CimInstance time-span json blocks.
type CimTimeDuration struct {
	time.Duration
}

// cimTimeDurationObject is a struct that represents the unmarshalled
// json of a CimInstance time duration object.
// It is used to do the initial unmarshalling of the json block.
type cimTimeDurationObject struct {
	Days         int32 `json:"Days"`
	Hours        int32 `json:"Hours"`
	Minutes      int32 `json:"Minutes"`
	Seconds      int32 `json:"Seconds"`
	MilliSeconds int32 `json:"Milliseconds"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for the CimTimeDuration type.
// It parses a JSON-encoded CimInstance time duration JSON block and converts it into a CimTimeDuration object.
func (t *CimTimeDuration) UnmarshalJSON(b []byte) error {
	var d cimTimeDurationObject

	// Unmarshal the json block into the cimTimeDurationObject struct.
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}

	// Convert the fields into a time.Duration object.
	duration := time.Duration(d.Days)*24*time.Hour +
		time.Duration(d.Hours)*time.Hour +
		time.Duration(d.Minutes)*time.Minute +
		time.Duration(d.Seconds)*time.Second +
		time.Duration(d.MilliSeconds)*time.Millisecond

	// Set the time.Duration object to the CimTimeDuration object.
	t.Duration = duration

	return nil
}
