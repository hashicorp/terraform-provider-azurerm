// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dates

import (
	"fmt"
	"time"
)

// ParseAsFormat parses the given nilable string as a time.Time using the specified
// format (for example RFC3339)
func ParseAsFormat(input *string, format string) (*time.Time, error) {
	if input == nil {
		return nil, nil
	}

	val, err := time.Parse(format, *input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", *input, err)
	}

	return &val, nil
}
