package validate

import (
	"fmt"
	"regexp"
)

func DedicatedHardwareSecurityModuleName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))

		return
	}

	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{1,22}[a-zA-Z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 24 alphanumeric characters. It must begin with a letter, end with a letter or digit.", k))

		return
	}

	// No consecutive hyphens
	if regexp.MustCompile("(--)").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must not contain any consecutive hyphens", k))
	}

	return
}
