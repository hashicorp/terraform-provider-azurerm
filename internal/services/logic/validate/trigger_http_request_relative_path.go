package validate

import (
	"fmt"
	"regexp"
)

func TriggerHttpRequestRelativePath(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9_/}{]+$").MatchString(value) {
		errors = append(errors, fmt.Errorf("Relative Path can only contain alphanumeric characters, underscores, forward slashes and curly braces."))
	}

	return warnings, errors
}
