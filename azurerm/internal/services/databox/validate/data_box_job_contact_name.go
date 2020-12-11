package validate

import (
	"fmt"
	"regexp"
)

func DataBoxJobContactName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\S][\s\S]{1,32}[\S]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 34 characters in length", k))
	}

	return warnings, errors
}
