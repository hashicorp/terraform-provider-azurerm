// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

var schedulerNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$`)

func SchedulerName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 3 || len(value) > 63 {
		errors = append(errors, fmt.Errorf("property `%s` must be between 3 and 63 characters, got %d", k, len(value)))
		return warnings, errors
	}

	if !schedulerNameRegex.MatchString(value) {
		errors = append(errors, fmt.Errorf("property `%s` must start and end with alphanumeric characters and can contain hyphens", k))
	}

	return warnings, errors
}
