package validate

import (
	"fmt"
	"regexp"
)

func RuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z_0-9.-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only word characters, numbers, underscores, periods, and hyphens allowed in %q: %q",
			k, value))
	}

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 80 characters: %q", k, value))
	}

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be an empty string: %q", k, value))
	}
	if !regexp.MustCompile(`[a-zA-Z0-9_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must end with a word character, number, or underscore: %q", k, value))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must start with a word character or number: %q", k, value))
	}

	return warnings, errors
}
