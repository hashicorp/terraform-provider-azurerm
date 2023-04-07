package validate

import (
	"fmt"
)

func AvailabilitySetName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 1 || len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length", k))
		return warnings, errors
	}

	return warnings, errors
}
