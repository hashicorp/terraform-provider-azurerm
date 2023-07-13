// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
)

func RolePassword(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if 8 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 8 character: %q", k, value))
	}

	if len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
