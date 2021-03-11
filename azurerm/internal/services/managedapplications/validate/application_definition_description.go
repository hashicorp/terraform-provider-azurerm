package validate

import (
	"fmt"
)

func ApplicationDefinitionDescription(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 200 {
		errors = append(errors, fmt.Errorf("%q should not exceed 200 characters in length.", k))
	}

	return warnings, errors
}
