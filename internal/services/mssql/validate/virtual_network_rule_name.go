package validate

import (
	"fmt"
	"regexp"
)

func VirtualNetworkRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Cannot be empty
	if len(value) == 0 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be an empty string: %q", k, value))
	}

	// Cannot be shorter than 2 characters
	if len(value) == 1 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be shorter than 2 characters: %q", k, value))
	}

	// Cannot be more than 64 characters
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	// Must only contain alphanumeric characters, underscores, periods or hyphens
	if !regexp.MustCompile(`^[A-Za-z0-9-\._]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q can only contain alphanumeric characters, underscores, periods and hyphens: %q",
			k, value))
	}

	// Cannot end in a hyphen
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen: %q", k, value))
	}

	// Cannot end in a period
	if regexp.MustCompile(`\.$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a period: %q", k, value))
	}

	// Cannot start with a period, underscore or hyphen
	if regexp.MustCompile(`^[\._-]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with a period, underscore or hyphen: %q", k, value))
	}

	// There are multiple returns in the case that there is more than one invalid
	// case applied to the name.
	return warnings, errors
}
