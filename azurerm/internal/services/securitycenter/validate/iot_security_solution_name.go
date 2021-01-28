package validate

import (
	"fmt"
	"regexp"
)

func IotSecuritySolutionName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	// The name attribute rules are :
	// 1. can only contain letter, digit, '-', '.' or '_'

	if !regexp.MustCompile(`^([-a-zA-Z0-9_.])+$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can only contain letter, digit, '-', '.' or '_'", v))
	}

	return
}
