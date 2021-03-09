package validate

import (
	"fmt"
	"regexp"
)

func ManagedHardwareSecurityModuleName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return warnings, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules are :
	// 1. can contain only alphanumeric characters.
	// 2. The first character must be a letter.
	// 3. The last character must be a letter or number
	// 4. The value must be between 3 and 24 characters long
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d]{2,23}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must begin with a letter, end with a letter or number, contain only alphanumeric characters. The value must be between 3 and 24 characters long", k))
	}

	return warnings, errors
}
