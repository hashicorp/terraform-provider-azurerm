package validate

import (
	"fmt"
	"regexp"
)

func PrivateLinkConfigurationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[A-Za-z][A-Za-z_.-]+[A-Za-z_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter, end with a letter or underscore, and may contain only letters, underscores, periods or hyphens", k))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 characters: %q", k, value))
	}

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
