package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type NullableTime struct {
	Time  time.Time
	Valid bool
}

func NewNullableTime(t time.Time) NullableTime {
	return NullableTime{
		Time:  t,
		Valid: true,
	}
}

func (nt *NullableTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time = time.Time{}
		nt.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time = v
		nt.Valid = true
		return nil
	case string:
		return nt.parse(v)
	case []byte:
		return nt.parse(string(v))
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into NullableTime", value)
	}
}

func (nt NullableTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}

	return nt.Time.Format("2006-01-02 15:04:05"), nil
}

func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}

	return []byte(`"` + nt.Time.Format(time.RFC3339) + `"`), nil
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Time = time.Time{}
		nt.Valid = false
		return nil
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	return nt.parse(raw)
}

func (nt *NullableTime) parse(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		nt.Time = time.Time{}
		nt.Valid = false
		return nil
	}

	layouts := []string{
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05.999",
		"2006-01-02 15:04:05",
		time.RFC3339Nano,
		time.RFC3339,
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, raw, time.Local)
		if err == nil {
			nt.Time = parsed
			nt.Valid = true
			return nil
		}
	}

	return fmt.Errorf("unable to parse time value %q", raw)
}
