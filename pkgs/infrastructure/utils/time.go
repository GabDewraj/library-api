package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}
type CustomDate struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ut *CustomTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	ut.Time = time.Unix(timestamp, 0)
	return nil
}

// Scan implements the sql.Scanner interface
func (ut *CustomTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		ut.Time = time.Time{} // Set to the zero value for time.Time
	case time.Time:
		ut.Time = v
	case []byte:
		// Assuming the value from the database is a string representation
		strValue := string(v)
		t, err := time.Parse("2006-01-02 15:04:05", strValue)
		if err != nil {
			return fmt.Errorf("parsing time %q: %w", strValue, err)
		}
		ut.Time = t
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ut *CustomDate) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	ut.Time = time.Unix(timestamp, 0)
	return nil
}

// Scan implements the sql.Scanner interface
func (ut *CustomDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ut.Time = v
	case []byte:
		// Assuming the value from the database is a string representation
		strValue := string(v)
		t, err := time.Parse("2006-01-02", strValue)
		if err != nil {
			return fmt.Errorf("parsing time %q: %w", strValue, err)
		}
		ut.Time = t
	}
	return nil
}
