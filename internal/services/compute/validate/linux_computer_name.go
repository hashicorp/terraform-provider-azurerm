package validate

import (
	"fmt"
	"strings"
)

func LinuxComputerNameFull(i interface{}, k string) (warnings []string, errors []error) {
	// Linux host name cannot exceed 64 characters in length
	return LinuxComputerName(i, k, 64, false)
}

func LinuxComputerNamePrefix(i interface{}, k string) (warnings []string, errors []error) {
	// Linux host name prefix cannot exceed 58 characters in length
	return LinuxComputerName(i, k, 58, true)
}

func LinuxComputerName(i interface{}, k string, maxLength int, allowDashSuffix bool) (warnings []string, errors []error) {
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

	if strings.HasPrefix(v, "_") {
		errors = append(errors, fmt.Errorf("%q cannot begin with an underscore", k))
	}

	if strings.HasSuffix(v, ".") {
		errors = append(errors, fmt.Errorf("%q cannot end with a period", k))
	}

	if !allowDashSuffix && strings.HasSuffix(v, "-") {
		errors = append(errors, fmt.Errorf("%q cannot end with a dash", k))
	}

	// Linux host name cannot contain the following characters
	specialCharacters := `\/"[]:|<>+=;,?*@&~!#$%^()_{}'`
	if strings.ContainsAny(v, specialCharacters) {
		errors = append(errors, fmt.Errorf("%q cannot contain the special characters: `%s`", k, specialCharacters))
	}

	return warnings, errors
}
