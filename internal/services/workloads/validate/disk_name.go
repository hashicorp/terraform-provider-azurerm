package validate

import (
	"fmt"
)

func DiskName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if value == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return warnings, errors
	}

	if len(value) < 1 || len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length", k))
		return warnings, errors
	}

	return warnings, errors
}
