package managedapplication

import (
	"fmt"
	"regexp"
)

func ValidateManagedAppDefinitionName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters or numbers.", k))
	}

	return warnings, errors
}

func ValidateManagedAppDefinitionDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{4,60}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 60 characters in length.", k))
	}

	return warnings, errors
}

func ValidateManagedAppDefinitionDescription(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{0,200}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q should not exceed 200 characters in length.", k))
	}

	return warnings, errors
}
