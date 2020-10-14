package jsontime

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// JSONTime converts a time from a different string formats to time.Time
type JSONTime struct{ time.Time }

// MarshalJSON outputs JSON.
func (d JSONTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.Local().Format(time.RFC3339Nano) + "\""), nil
}

// UnmarshalJSON handles incoming JSON.
func (d *JSONTime) UnmarshalJSON(b []byte) error {
	return d.tryParse(strings.Trim(string(b), "\""))
}

func (d *JSONTime) tryParse(s string) (err error) {
	//attempt 1 - RFC3339 format
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		*d = JSONTime{t.Local()}
		return
	}
	//attempt 1.5 - RFC3339Nano format
	t, err = time.Parse(time.RFC3339Nano, s)
	if err == nil {
		*d = JSONTime{t.Local()}
		return
	}

	//attempt 2 - datetime with milliseconds format
	t, err = time.ParseInLocation("2006-01-02T15:04:05.999999999", s, time.Local)
	if err == nil {
		*d = JSONTime{t.Local()}
		return
	}

	//attempt 3 - sql server
	t, err = time.ParseInLocation("2006-01-02T15:04:05", s, time.Local)
	if err == nil {
		*d = JSONTime{t}
		return
	}
	//attempt 4 - sql server with Z
	t, err = time.Parse("2006-01-02T15:04:05Z", s)
	if err == nil {
		*d = JSONTime{t}
		return
	}
	err = fmt.Errorf("No suitable format found for a string %s", s)
	return
}

//MarshalDynamoDBAttributeValue marshals object to a dynamodb attribute
func (d *JSONTime) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	t := d.Time.Local().Format(time.RFC3339Nano)
	av.S = &t
	return nil
}

//UnmarshalDynamoDBAttributeValue marshals object from a dynamodb attribute
func (d *JSONTime) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.S == nil {
		return nil
	}

	return d.tryParse(*av.S)
}
