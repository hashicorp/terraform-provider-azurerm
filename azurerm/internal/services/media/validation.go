package media

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateMediaServicesAccountName(i interface{}, k string) (warnings []string, errors []error) {
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

func ValidateStreamingEnpointName(i interface{}, k string) (warnings []string, errors []error) {
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

	const maxLength = 32
	// Streaming endpoint name can be 1-32 characters in length
	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q cannot exceed 32 characters.", k))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9]+(-*[a-zA-Z0-9])*$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q can only contain alphanumeric characters and hyphens. Must not begin or end with hyphen.", k))
	}

	return warnings, errors
}
