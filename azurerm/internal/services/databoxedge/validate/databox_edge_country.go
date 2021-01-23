package validate

import (
	"fmt"
	"regexp"
)

func DataboxEdgeCountry(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z]{2,3}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 3 characters in length", k))
	}

	return warnings, errors
}
