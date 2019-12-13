package netapp

import (
	"fmt"
	"regexp"
)

func ValidateNetAppAccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-_\da-zA-Z]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppPoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-_\da-zA-Z]{2,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and start with letters or numbers and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppVolumeName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-_\da-zA-Z]{0,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length and start with letters and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppVolumeVolumePath(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-\da-zA-Z]{0,79}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length and start with letters and contains only letters, numbers or hyphens.", k))
	}

	return warnings, errors
}
