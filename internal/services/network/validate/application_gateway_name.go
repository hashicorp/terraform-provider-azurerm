package validate

import (
	"fmt"
	"regexp"
)

func ApplicationGatewayName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	// Validate name: 1-80 chars, begin with letter/number, end with letter/number/underscore,
	// and may contain letters, numbers, underscores, periods, or hyphens in the middle
	pattern := `^[a-zA-Z\d]$|^[a-zA-Z\d][a-zA-Z\d-_.]{0,78}[a-zA-Z\d_]$`
	if matched := regexp.MustCompile(pattern).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore (_), and may contain only alphanumeric characters, underscores (_), hyphens (-), and periods (.)", k))
	}

	return warnings, errors
}
