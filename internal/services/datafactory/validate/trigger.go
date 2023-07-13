// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func TriggerTimespan(i interface{}, k string) (warnings []string, errors []error) {
	value, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if !regexp.MustCompile(`^\-?((\d+)\.)?(\d\d):(60|([0-5][0-9])):(60|([0-5][0-9]))`).MatchString(value) {
		errors = append(errors, fmt.Errorf("invalid timespan, must be of format hh:mm:ss %q: %q", k, value))
	}

	return warnings, errors
}
