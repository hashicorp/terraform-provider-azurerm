package validate

import (
	"fmt"
	"regexp"
)

func ApplicationVersion(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-._\da-zA-Z]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q can contain any combination of alphanumeric characters, hyphens, underscores, and periods: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
