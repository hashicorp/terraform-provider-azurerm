// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func LocalUserName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-z0-9]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"Username must be between 3 and 64 characters in length, use numbers and lower-case letters only: %q", value))
	}

	return warnings, errors
}
