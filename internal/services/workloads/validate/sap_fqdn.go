package validate

import (
	"fmt"
)

func SAPFQDN(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 2 || len(value) > 34 {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 34 characters in length", k))
		return warnings, errors
	}

	return warnings, errors
}
