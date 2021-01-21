package validate

import (
	"fmt"
	"regexp"
)

func DataboxEdgeStreetAddress(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{1,35}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 35 characters in length", k))
	}

	return warnings, errors
}
