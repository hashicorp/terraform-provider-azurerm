package validate

import (
	"fmt"
)

func AdminUsername(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if value == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return warnings, errors
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q may contain at most 64 characters", k))
		return warnings, errors
	}

	return warnings, errors
}
