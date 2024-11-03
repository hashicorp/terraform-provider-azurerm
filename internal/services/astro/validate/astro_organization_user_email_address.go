package validate

import (
	"fmt"
	"regexp"
)

func UserEmailAddress(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	regex := `^[A-Za-z0-9._%+-]+@(?:[A-Za-z0-9-]+\.)+[A-Za-z]{2,}$`
	if matched := regexp.MustCompile(regex).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must be a valid email address that matches the regular expression `%s`", k, regex))
	}

	return warnings, errors
}
