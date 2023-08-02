package utils

import (
	"database/sql/driver"
	"time"
)

const layout = "2006-01-02 15:04:05"

type NullTime struct {
	time.Time
	Valid bool
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n *NullTime) Scan(src interface{}) error {
	n.Time, n.Valid = src.(time.Time)
	return nil
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(layout)+2)
	b = append(b, '"')
	b = n.AppendFormat(b, layout)
	b = append(b, '"')
	return b, nil
}

func (n *NullTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	t, err := time.Parse(s, layout)
	if err != nil {
		return err
	}
	n.Time, n.Valid = t, true
	return nil
}
