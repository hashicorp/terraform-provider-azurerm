package validate

import (
	"fmt"
	"regexp"
)

func SqlAdministratorLoginName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can contain only letters or numbers.
	// 2. must start with letter
	// 3. The value must be between 1 and 128 characters long

	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d]{0,127}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only letters or numbers, must start with a letter, The value must be between 1 and 128 characters long", k))
		return
	}

	return warnings, errors
}
