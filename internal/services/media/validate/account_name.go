package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func AccountName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("Expected %q to be a string but it wasn't!", k))
		return
	}

	// The value must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	if matched := regexp.MustCompile(`^[-a-z0-9]{3,24}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must be 3 - 24 characters long, contain only lowercase letters and numbers.", k))
	}

	return warnings, errors
}
