package validate

import (
	"fmt"
	"regexp"
)

func VolumePath(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-\da-zA-Z]{0,79}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length and start with letters and contains only letters, numbers or hyphens.", k))
	}

	return warnings, errors
}
