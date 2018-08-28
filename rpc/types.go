package rpc

import (
	"encoding/json"
	"time"
)

const dateFormat = "2006-01-02"

// Time defines custom marshalling and unmarshalling format compatible with RFC3339
type Time time.Time

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) error {
	str := ""
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	tPtr, err := time.Parse(time.RFC3339, str)
	*t = Time(tPtr)

	return err
}

// Date defines custom marshalling and unmarshalling format for time which only includes date without time
type Date time.Time

// MarshalJSON implements the json.Marshaler interface.
func (t Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(dateFormat))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *Date) UnmarshalJSON(data []byte) error {
	str := ""
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	tPtr, err := time.Parse(dateFormat, str)
	*t = Date(tPtr)

	return err
}
