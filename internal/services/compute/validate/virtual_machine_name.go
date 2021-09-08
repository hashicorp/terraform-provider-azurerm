package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func VirtualMachineName(i interface{}, k string) (warnings []string, errors []error) {
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

	const maxLength = 80
	// VM name can be 1-80 characters in length
	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, dots, dashes and underscores", k))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9]`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must begin with an alphanumeric character", k))
	}

	if matched := regexp.MustCompile(`\w$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must end with an alphanumeric character or underscore", k))
	}

	// Portal: Virtual machine name cannot contain only numbers.
	if matched := regexp.MustCompile(`^\d+$`).Match([]byte(v)); matched {
		errors = append(errors, fmt.Errorf("%q cannot contain only numbers", k))
	}

	return warnings, errors
}
