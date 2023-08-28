// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func TimeInterval(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}
	if matched := regexp.MustCompile(`^([0-9][0-9]):([0-5][0-9]):([0-5][0-9])$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must be in the form HH:MM:SS between 00:00:00 and 99:59:59", k))
	}

	return warnings, errors
}
