package validate

import (
	"fmt"
	"regexp"
)

// GenericRFC3986Compliance - This validation rule is to act as a general minimum bar that
// all azure resource names should meet instead of falling back to just NoEmptyStrings
// if they do not have any known validation rules.
func GenericRFC3986Compliance(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	// The minimum generic rules are:
	// 1. Must not be empty.
	// 2. Must be between 1 and 80 characters.
	// 3. The attribute must:
	//    a) begin with a letter or number
	//    b) end with a letter, number or underscore
	//    c) may contain only letters, numbers, underscores, periods, or hyphens.

	if len(v) == 1 {
		if matched := regexp.MustCompile(`^([a-zA-Z\d])`).Match([]byte(v)); !matched {
			errors = append(errors, fmt.Errorf("%s must begin with a letter or number", k))
		}
	} else {
		if matched := regexp.MustCompile(`^([a-zA-Z\d])([a-zA-Z\d-\_\.]{0,78})([a-zA-Z\d\_])$`).Match([]byte(v)); !matched {
			errors = append(errors, fmt.Errorf("%s must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens", k))
		}
	}

	return warnings, errors
}
