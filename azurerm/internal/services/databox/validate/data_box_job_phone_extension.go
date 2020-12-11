package validate

import (
	"fmt"
	"regexp"
)

func DataBoxJobPhoneExtension(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\d]{0,4}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain only numeric values with a maximum length of 4 digits", k))
	}

	return warnings, errors
}
