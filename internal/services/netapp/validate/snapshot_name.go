package validate

import (
	"fmt"
	"regexp"
)

func SnapshotName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-_\da-zA-Z]{3,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 64 characters in length and start with letters or numbers and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}
