// Package datetime provides a ParseTime function for turning commonly-used ISO 8601 date/time
// formats into Golang time.Time variables.
//
// Unlike Go's built-in RFC-3339 time format, this package automatically supports ISO 8601 date and
// time stamps with varying levels of granularity.  Examples:
package datetime

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"
)

// DefaultUTC is just like a time.Time but serializes as an RFC3339Nano format when stringified or
// marshaled to JSON.  When parsed/unmarshaled, it uses time.UTC as the location for timestamps that
// don't specify one.
type DefaultUTC time.Time

// String returns the DefaultUTC's RFC3339Nano representation.
func (d DefaultUTC) String() string {
	t := time.Time(d)
	return t.Format(time.RFC3339Nano)
}

// UnmarshalJSON implements the JSON Unmarshaler interface, allowing datetime.DefaultUTC struct fields
// to be read from JSON string fields.
func (d *DefaultUTC) UnmarshalJSON(data []byte) error {
	t, err := JSONParse(data, time.UTC)
	*d = DefaultUTC(t)
	return err
}

// Scan implements the sql Scanner interface, allowing datetime.DefaultUTC fields to be read from
// database columns.
func (d *DefaultUTC) Scan(value interface{}) error {
	t, err := sqlScan(value, time.UTC)
	*d = DefaultUTC(t)
	return err
}

// Value implements the sql Valuer interface, allowing datetime.DefaultUTC fields to be saved to
// database columns.
func (d DefaultUTC) Value() (driver.Value, error) {
	return d.String(), nil
}

// DefaultLocal is just like a time.Time but serializes as an RFC3339Nano format when stringified or
// marshaled to JSON.  When parsed/unmarshaled, it uses time.Local as the location for timestamps
// that don't specify one.
type DefaultLocal time.Time

// String returns the DefaultLocal's RFC3339Nano representation.
func (d DefaultLocal) String() string {
	t := time.Time(d)
	return t.Format(time.RFC3339Nano)
}

// UnmarshalJSON implements the JSON Unmarshaler interface, allowing datetime.DefaultLocal struct fields
// to be read from JSON string fields.
func (d *DefaultLocal) UnmarshalJSON(data []byte) error {
	t, err := JSONParse(data, time.Local)
	*d = DefaultLocal(t)
	return err
}

// Scan implements the sql Scanner interface, allowing datetime.DefaultLocal fields to be read from
// database columns.
func (d *DefaultLocal) Scan(value interface{}) error {
	t, err := sqlScan(value, time.Local)
	*d = DefaultLocal(t)
	return err
}

// Value implements the sql Valuer interface, allowing datetime.DefaultLocal fields to be saved to
// database columns.
func (d DefaultLocal) Value() (driver.Value, error) {
	return d.String(), nil
}

// Below here are helper funcs used by the DefaultUTC and Local types.
func parseBytes(b []byte, loc *time.Location) (time.Time, error) {
	p := newParser(bytes.NewBuffer(b))
	return p.parse(loc)
}

func sqlScan(value interface{}, loc *time.Location) (time.Time, error) {
	switch v := value.(type) {
	case []byte:
		t, err := parseBytes(v, loc)
		if err != nil {
			return zeroTime, err
		}
		return t, nil
	case string:
		t, err := parseBytes([]byte(v), loc)
		if err != nil {
			return zeroTime, err
		}
		return t, nil
	default:
		return zeroTime, fmt.Errorf("can only scan string and []byte, not %v", reflect.TypeOf(value))
	}
}

const doubleQuote byte = 34

// JSONParse will take a JSON bytes value with quotes around it, and parse it into a time.Time.
func JSONParse(data []byte, loc *time.Location) (time.Time, error) {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return zeroTime, nil
	}
	if data[0] != doubleQuote || data[len(data)-1] != doubleQuote {
		return zeroTime, fmt.Errorf("%s does not begin and end with double quotes", data)
	}
	trimmed := data[1 : len(data)-1]

	return parseBytes(trimmed, loc)
}
