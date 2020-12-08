package validate

import (
	"fmt"
)

func SqlFilter(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// SqlFilters can not be empty
	if len(v) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		return warnings, errors
	}

	// SqlFilters have a maximum length of 1024
	if len(v) > 1024 {
		errors = append(errors, fmt.Errorf("%q is of length %d, which exceeds the maximum length of 1024", k, len(v)))
		return
	}

	return warnings, errors
}
