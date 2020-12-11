package validate

import (
	"fmt"
	"regexp"
)

func DataBoxJobDiskPassKey(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 12 || len(value) > 32 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 12 and 32 characters in length, received: %q", k, value))
	}

	if !regexp.MustCompile(`^.*[\d]+.*[^a-zA-Z0-9 ]+.*|.*[^a-zA-Z0-9 ]+.*[\d]+.*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be alphanumeric, contain at least one special character and at least one number", k))
	}

	return warnings, errors
}
