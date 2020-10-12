package jsontime

import (
	"encoding/json"
	"testing"
	"time"
)

func Test(t *testing.T) {
	loc, _ := time.LoadLocation("Pacific/Auckland")
	tests := []struct {
		input        string
		expectedRFC  time.Time
		expectedCust time.Time
	}{
		{
			input:        `{"t1":"2020-09-22T08:26:52.8585767+12:00","t2":"2020-09-15T14:45:33.3034643"}`,
			expectedRFC:  time.Date(2020, time.September, 22, 8, 26, 52, 858576700, loc),
			expectedCust: time.Date(2020, time.September, 15, 14, 45, 33, 303464300, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var e struct {
				T1 JSONTime `json:"t1"`
				T2 JSONTime `json:"t2"`
			}
			err := json.Unmarshal([]byte(test.input), &e)
			if err != nil {
				t.Errorf("unexpected error during unmarshalling: %v", err)
			}
			if !test.expectedRFC.Equal(e.T1.Time) {
				t.Errorf("expected RFC3339 of %v, got %v", test.expectedRFC, e.T1)
			}
			if !test.expectedCust.Equal(time.Time(e.T2.Time)) {
				t.Errorf("expected \"2006-01-02T15:04:05.0000000\" of %v, got %v", test.expectedCust, e.T2)
			}
		})
	}
}
