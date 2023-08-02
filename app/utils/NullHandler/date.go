package utils

import (
	"database/sql/driver"
	"time"
)

const dateLayout = "2006-01-02"

type NullDate struct {
	time.Time
	Valid bool
}

func (n NullDate) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n *NullDate) Scan(src interface{}) error {
	n.Time, n.Valid = src.(time.Time)
	return nil
}

func (n NullDate) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(dateLayout)+2)
	b = append(b, '"')
	b = n.AppendFormat(b, dateLayout)
	b = append(b, '"')
	return b, nil
}

func (n *NullDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	t, err := time.Parse(s, dateLayout)
	if err != nil {
		return err
	}
	n.Time, n.Valid = t, true
	return nil
}
