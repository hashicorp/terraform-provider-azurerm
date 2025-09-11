// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"database/sql/driver"
	"fmt"
)

// Scan parses some value, which can be either string or []byte.
// It implements sql.Scanner, https://golang.org/pkg/database/sql/#Scanner
func (period *Period) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	err = nil
	switch v := value.(type) {
	case []byte:
		*period, err = Parse(string(v), false)
	case string:
		*period, err = Parse(v, false)
	default:
		err = fmt.Errorf("%T %+v is not a meaningful period", value, value)
	}

	return err
}

// Value converts the period to a string. It implements driver.Valuer,
// https://golang.org/pkg/database/sql/driver/#Valuer
func (period Period) Value() (driver.Value, error) {
	return period.String(), nil
}
