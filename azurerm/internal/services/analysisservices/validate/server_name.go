package validate

import (
	"fmt"
	"regexp"
)

func ServerName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-z][0-9a-z]{2,62}$`).Match([]byte(value)) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter, be lowercase alphanumeric, and be between 3 and 63 characters in length", k))
	}

	return warnings, errors
}
