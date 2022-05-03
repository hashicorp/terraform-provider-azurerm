package validate

import (
	"fmt"
	"regexp"
)

func FunctionAppFunctionName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z](([0-9a-zA-Z-]{0,30})[0-9a-zA-Z])?$`).Match([]byte(v)); !matched {
		warnings = append(warnings, fmt.Sprintf("%q up to version 4.x of Azure Functions Core Tools, the function name will be truncated to 32 characters", key))
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z](([0-9a-zA-Z-]{0,126})[0-9a-zA-Z])?$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must start with a letter, may only contain alphanumeric characters and dashes and up to 128 characters in length", key))
	}

	return warnings, errors
}
