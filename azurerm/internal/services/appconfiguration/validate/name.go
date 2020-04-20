package validate

import (
	"fmt"
	"regexp"
)

func AppConfigurationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{5,50}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and must be between 5-50 chars", k))
	}

	return warnings, errors
}
