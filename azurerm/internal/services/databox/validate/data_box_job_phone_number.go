package validate

import (
	"fmt"
	"regexp"
)

func DataBoxJobPhoneNumber(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[+][\d]{2,}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must begin with '+' and must be at least 2 digits in length", k))
	}

	return warnings, errors
}
