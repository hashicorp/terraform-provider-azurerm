package validate

import (
	"fmt"
	"regexp"
)

func SynapseBigDataPoolName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. must contain only letters or numbers.
	// 2. must start with a letter.
	// 3. must be between 1 and 15 characters long

	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d]{0,14}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only letters or numbers, must start with a letter, and be between 1 and 15 characters long", k))
		return
	}

	return warnings, errors
}
