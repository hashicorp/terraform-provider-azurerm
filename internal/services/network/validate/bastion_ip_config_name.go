package validate

import (
	"fmt"
	"regexp"
)

func BastionIPConfigName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[A-Za-z0-9][a-zA-Z0-9_.-]+[a-zA-Z0-9_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens. %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 characters: %q", k, value))
	}

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
