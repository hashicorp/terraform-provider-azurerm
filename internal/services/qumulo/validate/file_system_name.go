package validate

import (
	"fmt"
	"regexp"
)

func FileSystemName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,13}[a-zA-Z0-9]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and hyphen symbol and can't start with hyphen or end with hyphen, must be between 2-15 chars", k))
	}

	return warnings, errors
}
