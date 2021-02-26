package validate

import (
	"fmt"
	"regexp"
)

func ApplicationDefinitionName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[^\W_]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters or numbers.", k))
	}

	return warnings, errors
}
