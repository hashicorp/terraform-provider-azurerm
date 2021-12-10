package validate

import (
	"fmt"
	"regexp"
)

func VirtualHubConnectionName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-_.\da-zA-Z]{0,78}[_\da-zA-Z]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters and begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.", k))
	}

	return warnings, errors
}
