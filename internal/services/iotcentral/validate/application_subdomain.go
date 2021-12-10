package validate

import (
	"fmt"
	"regexp"
)

func ApplicationSubdomain(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile(`^[a-z\d][a-z\d-]{0,61}[a-z\d]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("test: %s, %q may only contain alphanumeric characters and dashes, length between 2-63", k, v))
	}
	return warnings, errors
}
