// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ManagedRedisAccessPolicyName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q expected type of to be string", k))
		return
	}

	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("%q length should be greater than %d, got %q", k, 1, v))
		return
	}

	if len(v) > 64 {
		errors = append(errors, fmt.Errorf("%q length should be less than or equal to %d characters, got %d", k, 64, len(v)))
		return
	}

	// Pattern: start/end with alphanumeric, can contain letters, numbers, hyphens, and spaces
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\s-]*[a-zA-Z0-9]$|^[a-zA-Z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must start and end with alphanumeric characters and can contain letters, numbers, hyphens, and spaces, got %v", k, v))
		return
	}

	return
}
