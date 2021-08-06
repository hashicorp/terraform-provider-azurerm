package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func WindowsComputerNameFull(i interface{}, k string) (warnings []string, errors []error) {
	// Windows computer name cannot be more than 15 characters long
	return windowsComputerName(i, k, 15)
}

func WindowsComputerNamePrefix(i interface{}, k string) (warnings []string, errors []error) {
	// Windows computer name prefix cannot be more than 9 characters long
	return windowsComputerName(i, k, 9)
}

func windowsComputerName(i interface{}, k string, maxLength int) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't!", k))
		return
	}

	// The value must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
	}

	if strings.HasSuffix(v, "-") {
		errors = append(errors, fmt.Errorf("%q cannot end with dash", k))
	}

	// A windows computer name can only contain alphanumeric characters and hyphens
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	// Windows computer name cannot contain only numbers
	if matched := regexp.MustCompile(`^\d+$`).Match([]byte(v)); matched {
		errors = append(errors, fmt.Errorf("%q cannot contain only numbers", k))
	}

	return warnings, errors
}
