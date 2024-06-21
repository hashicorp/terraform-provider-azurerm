// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ContainerRegistryCacheRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"alpha numeric characters optionally separated by '-' only are allowed in %q: %q", k, value))
	}

	if 5 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 5 characters: %q", k, value))
	}

	if len(value) >= 50 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 50 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
