package databox

import (
	"fmt"
	"regexp"
)

func ValidateDataBoxJobName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-\da-zA-Z]{1,22}[\da-zA-Z]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 24 characters in length, and it must begin and end with an alphanumeric and can only contain alphanumeric characters and hyphens", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobContactName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{3,34}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 34 characters in length", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobPhoneNumber(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[+][\d]{2,}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must begin with + and may contain only at least 2 numbers", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be with @", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobPhoneExtension(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\d]{0,4}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must have maximum 4 characters in length and can only contain numbers", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobStreetAddress(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{1,35}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 35 characters in length", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobPostCode(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{1,9}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 9 characters in length", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobCity(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{2,30}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 30 characters in length", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobCompanyName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{2,35}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 35 characters in length", k))
	}

	return warnings, errors
}

func ValidateDataBoxJobDiskPassKey(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 12 || len(value) > 32 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 12 and 32 characters: %q", k, value))
	}

	if !regexp.MustCompile(`^.*[\d]+.*[^a-zA-Z0-9 ]+.*|.*[^a-zA-Z0-9 ]+.*[\d]+.*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be alphanumeric, contain at least one special character and atleast one number", k))
	}

	return warnings, errors
}
