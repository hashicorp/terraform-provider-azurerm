package validate

import (
	"fmt"
	"regexp"
)

func AccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-z0-9]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("lowercase letters and numbers only are allowed in %q: %q", k, value))
	}

	if 3 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 3 characters: %q", k, value))
	}

	if len(value) > 24 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 24 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
