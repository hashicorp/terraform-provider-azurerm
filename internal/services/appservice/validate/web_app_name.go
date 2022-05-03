package validate

import (
	"fmt"
	"regexp"
)

func WebAppName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,32}$`).Match([]byte(value)); !matched {
		warnings = append(warnings, fmt.Sprintf("%q up to version 4.x of Azure Functions Core Tools, the function name will be truncated to 32 characters", k))
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,60}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and up to 60 characters in length", k))
	}

	return warnings, errors
}
