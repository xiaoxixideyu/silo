package isotime

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Date struct {
	time.Time
	Valid bool
}

func NewDateFromTime(t time.Time) Date {
	d := Date{}
	d.Time = t
	d.Valid = true
	return d
}

func NewDateFromISOTime(t ISOTime) Date {
	d := Date{}
	d.Time = t.ToTime()
	d.Valid = true
	return d
}

func NewDateFromString(str string) (Date, error) {
	d := Date{}
	t, err := time.Parse(dateFormat, str)
	if err != nil {
		d.Valid = false
		d.Time = time.Time{}
		return d, err
	}
	d.Time = t
	d.Valid = true
	return d, nil
}

func (d Date) FromISOTime(t ISOTime) Date {
	d.Time = t.ToTime()
	d.Valid = true
	return d
}

// String returns the date in the format YYYY-MM-DD
func (d *Date) String() string {
	if !d.Valid {
		return ""
	}
	return d.Format(dateFormat)
}

// JsonString returns the date in the format YYYY-MM-DD
func (d Date) JsonString() string {
	if !d.Valid {
		return ""
	}
	return d.Format(jsonFormat)
}

func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.String(), nil
}

func (d *Date) Scan(value any) error {
	if value == nil {
		d.Valid = false
		return nil
	}
	var err error
	switch p := value.(type) {
	case time.Time:
		d.Time = p
		d.Valid = true
	case string:
		d.Time, err = time.Parse(dateFormat, p)
		if err != nil {
			return err
		}
		d.Valid = true
	case []byte:
		d.Time, err = time.Parse(dateFormat, string(p))
		if err != nil {
			return err
		}
		d.Valid = true
	default:
		return errors.New("invalid type for Date")
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of d.String().
func (d Date) MarshalText() ([]byte, error) {
	if !d.Valid {
		return []byte(""), nil
	}
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseDate.
func (d *Date) UnmarshalText(data []byte) error {
	var err error
	*d, err = NewDateFromString(string(data))
	return err
}

// MarshalJSON .
func (d *Date) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return nil, nil
	}
	return []byte(d.JsonString()), nil
}

// UnmarshalJSON .
func (d *Date) UnmarshalJSON(b []byte) error {
	if b == nil {
		d.Valid = false
		return nil
	}
	if string(b) == "null" {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return json.Unmarshal(b, &d.Time)
}

// Equal .
func (d Date) Equal(v Date) bool {
	return d.String() == v.String()
}
