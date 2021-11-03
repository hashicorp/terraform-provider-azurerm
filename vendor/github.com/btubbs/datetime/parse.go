package datetime

import (
	"bytes"
	"time"
)

// Parse takes a string with a ISO 8601 timestamp in it, and a default location to use for
// timestamps that don't include one, and returns a time.Time.
func Parse(s string, defaultLocation *time.Location) (time.Time, error) {
	p := newParser(bytes.NewBuffer([]byte(s)))
	return p.parse(defaultLocation)
}

// ParseUTC takes a string with a ISO 8601 timestamp in it and returns a time.Time.  For inputs
// that do not specify a location, time.UTC will be used.
func ParseUTC(s string) (time.Time, error) { return Parse(s, time.UTC) }

// ParseLocal takes a string with a ISO 8601 timestamp in it and returns a time.Time.  For inputs
// that do not specify a location, time.Local will be used.
func ParseLocal(s string) (time.Time, error) { return Parse(s, time.Local) }
