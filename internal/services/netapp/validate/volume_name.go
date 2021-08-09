package validate

import (
	"fmt"
	"regexp"
)

func VolumeName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-_\da-zA-Z]{0,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length and start with letters and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}
