package validate

import "fmt"

func ApplicationDefinitionDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 4 || len(value) > 60 {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 60 characters in length.", k))
	}

	return warnings, errors
}
