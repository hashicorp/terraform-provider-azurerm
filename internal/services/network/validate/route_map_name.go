// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func RouteMapName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z_.-]+[a-zA-Z_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("The name must begin with a letter, end with a letter or underscore, and may contain only letters, underscores, periods or hyphens. %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 characters: %q", k, value))
	}

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
