package validate

import (
	"fmt"
	"regexp"
)

func AgentPoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"alpha numeric characters only are allowed in %q: %q", k, value))
	}

	if 3 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 3 characters: %q", k, value))
	}

	if len(value) > 20 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 20 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
