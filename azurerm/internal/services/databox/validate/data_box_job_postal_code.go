package validate

import (
	"fmt"
	"regexp"
)

func DataBoxJobPostalCode(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{1,9}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 9 characters in length", k))
	}

	return warnings, errors
}
