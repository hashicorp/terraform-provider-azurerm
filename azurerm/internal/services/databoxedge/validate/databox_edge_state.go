package validate

import (
	"fmt"
	"regexp"
)

func DataboxEdgeState(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{1,30}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 30 characters in length", k))
	}

	return warnings, errors
}
