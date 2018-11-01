package azure

import (
	"fmt"
	"regexp"
)

func ValidateResourceID(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseAzureResourceID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
	}

	return
}

//true for a resource ID or an empty string
func ValidateResourceIDOrEmpty(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" {
		return
	}

	return ValidateResourceID(i, k)
}

//true for a resource ID or an empty string
func ValidateServiceName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	//First, second, and last characters must be a letter or number with a total length between 3 to 50 lowercase characters.
	r := regexp.MustCompile("^[a-z0-9]{2}[-a-z0-9]{0,47}[a-z0-9]{1}$")
	if !r.MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be 3 - 50 characters in length", k))
		errors = append(errors, fmt.Errorf("%q first, second, and last characters must be a lowercase letter or number", k))
		errors = append(errors, fmt.Errorf("%q can only contain lowercase letters, numbers and hyphens", k))
	}

	//No consecutive dashes.
	r = regexp.MustCompile("(--)")
	if r.MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must not contain any consecutive hyphens", k))
	}

	return
}
