package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LocalTime is a custom time type that can scan both time.Time and []uint8 (byte slice)
// from MySQL. This handles the case where parseTime=True is not working or the column
// type is VARCHAR instead of DATETIME.
type LocalTime struct {
	time.Time
	Valid bool
}

// Scan implements the sql.Scanner interface
func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		t.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		t.Valid = true
		return nil
	case []uint8:
		str := string(v)
		if str == "" || str == "0000-00-00 00:00:00" || str == "0000-00-00 00:00:00.000" {
			t.Valid = false
			return nil
		}
		// Try multiple formats
		formats := []string{
			"2006-01-02 15:04:05.000",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, f := range formats {
			parsed, err := time.Parse(f, str)
			if err == nil {
				t.Time = parsed
				t.Valid = true
				return nil
			}
		}
		// If all formats fail, just mark as invalid
		t.Valid = false
		return nil
	case string:
		if v == "" || v == "0000-00-00 00:00:00" {
			t.Valid = false
			return nil
		}
		formats := []string{
			"2006-01-02 15:04:05.000",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, f := range formats {
			parsed, err := time.Parse(f, v)
			if err == nil {
				t.Time = parsed
				t.Valid = true
				return nil
			}
		}
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("LocalTime.Scan: unsupported type %T", value)
	}
}

// Value implements the driver.Valuer interface
func (t LocalTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}
