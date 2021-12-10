package validate

import (
	"fmt"
	"regexp"
)

// PoolName validates the name of a Batch pool
func PoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"any combination of alphanumeric characters including hyphens and underscores are allowed in %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
