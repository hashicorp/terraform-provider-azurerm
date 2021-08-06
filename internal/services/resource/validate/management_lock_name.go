package validate

import (
	"fmt"
	"regexp"
)

func ManagementLockName(v interface{}, k string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`[A-Za-z0-9-_]`).MatchString(input) {
		errors = append(errors, fmt.Errorf("%s can only consist of alphanumeric characters, dashes and underscores", k))
	}

	if len(input) >= 260 {
		errors = append(errors, fmt.Errorf("%s can only be a maximum of 260 characters", k))
	}

	return warnings, errors
}
