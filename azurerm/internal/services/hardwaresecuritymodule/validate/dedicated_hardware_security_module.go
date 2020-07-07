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

	if !regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q can contain only letters, numbers, and hyphens. It must be between 3 and 24 characters long.", k))

		return
	}

	return
}
