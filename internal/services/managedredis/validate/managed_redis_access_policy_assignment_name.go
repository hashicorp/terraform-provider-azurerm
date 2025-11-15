// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ManagedRedisAccessPolicyAssignmentName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q expected type of to be string", k))
		return
	}

	// Pattern from Azure API: ^[A-Za-z0-9]{1,60}$ (letters and numbers only, 1-60 characters)
	if !regexp.MustCompile(`^[A-Za-z0-9]{1,60}$`).MatchString(v) {
		if len(v) == 0 {
			errors = append(errors, fmt.Errorf("%q cannot be empty", k))
		} else if len(v) > 60 {
			errors = append(errors, fmt.Errorf("%q length should be less than or equal to 60 characters, got %d", k, len(v)))
		} else {
			errors = append(errors, fmt.Errorf("%q can only contain letters and numbers, got %v", k, v))
		}
		return
	}

	return
}
