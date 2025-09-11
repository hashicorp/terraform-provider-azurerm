// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

// MarshalBinary implements the encoding.BinaryMarshaler interface.
// This also provides support for gob encoding.
func (period Period) MarshalBinary() ([]byte, error) {
	// binary method would take more space in many cases, so we simply use text
	return period.MarshalText()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// This also provides support for gob decoding.
func (period *Period) UnmarshalBinary(data []byte) error {
	return period.UnmarshalText(data)
}

// MarshalText implements the encoding.TextMarshaler interface for Periods.
// This also provides support for JSON encoding.
func (period Period) MarshalText() ([]byte, error) {
	return []byte(period.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Periods.
// This also provides support for JSON decoding.
func (period *Period) UnmarshalText(data []byte) (err error) {
	u, err := Parse(string(data), false)
	if err == nil {
		*period = u
	}
	return err
}
