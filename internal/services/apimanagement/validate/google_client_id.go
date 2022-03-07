package validate

import (
	"fmt"
	"regexp"
)

func GoogleClientID(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[A-Za-z0-9-]+\.apps\.googleusercontent\.com$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must start with an identifier containing alphanumeric characters and hyphens and end with '.apps.googleusercontent.com'", k))
	}
	return warnings, errors
}
