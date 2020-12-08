package validate

import (
	"fmt"
	"regexp"
)

func AppServiceEnvironmentName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z]{0,61}[0-9a-zA-Z]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes up to 60 characters in length, and must start and end in an alphanumeric", k))
	}

	return warnings, errors
}
