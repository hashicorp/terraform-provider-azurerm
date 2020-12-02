package validate

import (
	"fmt"
	"regexp"
)

func PowerBIEmbeddedName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-z][a-z0-9]{3,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 64 characters in length and contains only lowercase letters or numbers.", k))
	}

	return warnings, errors
}
