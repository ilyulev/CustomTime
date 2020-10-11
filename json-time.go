package jsontime

import (
	"fmt"
	"time"
)

// JSONTime converts a time from a different string formats to time.Time
type JSONTime struct{ time.Time }

// MarshalJSON outputs JSON.
func (d JSONTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.Format(time.RFC3339) + "\""), nil
}

// UnmarshalJSON handles incoming JSON.
func (d *JSONTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	//attempt 1 - RFC3339 format
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		*d = JSONTime{t}
		return
	}

	//attempt 2 - datetime with milliseconds format
	t, err = time.Parse("2006-01-02T15:04:05.000", s)
	if err == nil {
		*d = JSONTime{t}
		return
	}
	return fmt.Errorf("No suitable format found for a string %s", s)
}
