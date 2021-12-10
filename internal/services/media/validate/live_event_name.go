package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func LiveEventName(i interface{}, k string) (warnings []string, errors []error) {
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

	// Live Event name can be 1-32 characters in length
	const maxLength = 32
	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q cannot exceed 32 characters.", k))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9]+(-*[a-zA-Z0-9])*$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q can contain letters, numbers, and hyphens (but the first and last character must be a letter or number).", k))
	}

	return warnings, errors
}
