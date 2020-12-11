package validate

import (
	"fmt"
	"regexp"
)

func DataBoxJobName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-\da-zA-Z]{1,22}[\da-zA-Z]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 24 characters in length, begin and end with an alphanumeric character, can only contain alphanumeric characters and hyphens", k))
	}

	return warnings, errors
}
