package validate

import (
	"fmt"
	"regexp"
)

func DataBoxJobEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^(([^<>()\[\]\.,;:\s@\"]+(\.[^<>()\[\]\.,;:\s@\"]+)*)|(\".+\"))@(([^<>()[\]\.,;:\s@\"]+\.)+[^<>()[\]\.,;:\s@\"]{2,})$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be with @", k))
	}

	return warnings, errors
}
