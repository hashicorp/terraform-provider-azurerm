package validate

import "fmt"

func ApplicationDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 1024 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 1024 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
