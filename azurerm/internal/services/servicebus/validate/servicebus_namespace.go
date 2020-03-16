package validate

import (
	"fmt"
	"regexp"
)

// validation
func ServiceBusNamespaceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}").MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must be a GUID", k))
	}

	if !regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{4,48}[a-zA-Z0-9]$").MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must contain only letters, numbers, and hyphens. The namespace must start with a letter, and it must end with a letter or number and be between 6 and 50 characters long", k))
	}

// Ask KT about this vs breaking them out separately
	// 	if strings.HasSuffix(v, "-") || strings.HasSuffix(v, "-sb") || strings.HasSuffix(v, "-mgmt") {
	// 	errors = append(errors, fmt.Errorf("%q cannot end with a hyphen, -sb, or -mgmt", k, value))
	// }

	// Cannot end in a hyphen
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen: %q", k, value))
	}

	// Cannot end in -sb
	if regexp.MustCompile(`-sb$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with -sb: %q", k, value))
	}

	// Cannot end in -mgmt
	if regexp.MustCompile(`-mgmt$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with -mgmt: %q", k, value))
	}

	return warnings, errors
}