package validate

import (
	"fmt"
	"regexp"
)

func ApplicationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-\da-zA-Z]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters, numbers or hyphens.", k))
	}

	return warnings, errors
}
