package validate

import (
	"fmt"
	"regexp"
)

func AccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-_\da-zA-Z]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}
