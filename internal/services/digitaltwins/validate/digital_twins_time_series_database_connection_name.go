// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DigitalTwinsTimeSeriesDatabaseConnectionName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if len(name) < 3 {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", 3, v))
		return
	}

	if len(name) > 50 {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", 50, v))
		return
	}

	if regexp.MustCompile(`^[0-9]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q should not contain only numbers, got %v", k, v))
		return
	}

	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter or number, end with a letter or number and contain only letters, numbers, and hyphens, got %v", k, v))
		return
	}
	return
}
