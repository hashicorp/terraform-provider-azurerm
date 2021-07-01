package validate

import (
	"fmt"
	"regexp"
)

func ExpressRouteConnectionName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_])$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods or hyphens", k))
	}

	return warnings, errors
}
