package validate

import (
	"fmt"
)

func AdministratorLoginPassword(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if v == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return warnings, errors
	}

	if len(v) < 8 {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", 8, v))
		return warnings, errors
	}

	if len(v) > 256 {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", 256, v))
		return warnings, errors
	}

	return warnings, errors
}
