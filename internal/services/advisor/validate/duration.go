// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func Duration(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^(?:[0-9]{1,2}:)?[0-9]{2}:[0-9]{2}:[0-9]{2}$`).Match([]byte(value)) {
		errors = append(errors, fmt.Errorf("%q must be in format DD:HH:MM:SS. If DD is 00, it has to be omit", k))
	}

	return warnings, errors
}
