// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// lintignore:V011,V001 // the length check is combined with character/format rules
func StorageQueueName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q", k,
		))
	}

	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q cannot start with a hyphen", k))
	}

	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q cannot end with a hyphen", k))
	}

	if len(value) > 63 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 63 characters", k,
		))
	}

	if len(value) < 3 {
		errors = append(errors, fmt.Errorf(
			"%q must be at least 3 characters", k,
		))
	}

	return warnings, errors
}
