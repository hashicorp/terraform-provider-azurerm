package validate

import (
	"fmt"
	"regexp"
)

func ImportExportJobName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. must be between 2 and 64 characters long.
	// 2. must start with a letter, and can contain only letters, numbers, and hyphens.

	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d-]{1,63}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter, contain only letters, numbers and hyphens. The length should be between 2 and 64 characters long.", k))
	}

	return warnings, errors
}

func ImportExportJobEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^(([^<>()\[\]\.,;:\s@\"]+(\.[^<>()\[\]\.,;:\s@\"]+)*)|(\".+\"))@(([^<>()[\]\.,;:\s@\"]+\.)+[^<>()[\]\.,;:\s@\"]{2,})$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be with @", k))
	}

	return warnings, errors
}

func ImportExportJobPhone(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(`%q must align to regex "^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$"`, k))
	}

	return warnings, errors
}
