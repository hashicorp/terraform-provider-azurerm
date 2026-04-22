// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func NestedItemName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 127 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 127 characters in length, got %d", k, len(value)))
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	return warnings, errors
}
