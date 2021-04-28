package validate

import (
	"fmt"
)

func StorageBlobIndexTagName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if len(value) == 0 || len(value) > 128 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 128 characters: %q", k, value))
	}
	return warnings, errors
}

func StorageBlobIndexTagValue(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if len(value) > 256 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 0 and 256 characters: %q", k, value))
	}
	return warnings, errors
}
