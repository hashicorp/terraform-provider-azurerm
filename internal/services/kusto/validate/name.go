package validate

import (
	"fmt"
	"regexp"
)

func DataConnectionName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain letters, digits, whitespaces, dashes and dots: %q", k, name))
	}

	if len(name) > 40 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 1 and 40 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func EntityName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain letters, digits, underscores, spaces, dashes and dots: %q", k, name))
	}

	if len(name) > 1024 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 1 and 1024 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func ClusterName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if !regexp.MustCompile(`^[a-z][a-z0-9]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter and may only contain alphanumeric characters: %q", k, name))
	}

	if len(name) < 4 || len(name) > 22 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 4 and 22 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func DatabaseName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, whitespaces, dashes and dots: %q", k, name))
	}

	if len(name) > 260 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 4 and 22 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func DatabasePrincipalAssignmentName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, whitespaces, dashes and dots: %q", k, name))
	}

	if len(name) > 260 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 4 and 22 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}
