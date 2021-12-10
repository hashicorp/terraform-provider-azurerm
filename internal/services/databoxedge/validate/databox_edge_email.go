package validate

import (
	"fmt"
	"regexp"
)

func DataboxEdgeEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^(([^<>()\[\]\.,;:\s@\"]+(\.[^<>()\[\]\.,;:\s@\"]+)*)|(\".+\"))@(([^<>()[\]\.,;:\s@\"]+\.)+[^<>()[\]\.,;:\s@\"]{2,})$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain an at symbol and have a domain of at least two characters in length", k))
	}

	return warnings, errors
}
