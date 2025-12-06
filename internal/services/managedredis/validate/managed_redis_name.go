// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ManagedRedisClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q expected type of to be string", k))
		return warnings, errors
	}

	if len(v) < 3 {
		errors = append(errors, fmt.Errorf("%q length should be greater than %d, got %q", k, 3, v))
		return warnings, errors
	}

	if len(v) > 63 {
		errors = append(errors, fmt.Errorf("%q length should be less than %d, got %q", k, 63, v))
		return warnings, errors
	}

	if strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must not contain any consecutive hyphens, got %q", k, v))
		return warnings, errors
	}

	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q can only contain letters, numbers and hyphens. The first and last characters must each be a letter or a number, got %v", k, v))
		return warnings, errors
	}

	return warnings, errors
}
