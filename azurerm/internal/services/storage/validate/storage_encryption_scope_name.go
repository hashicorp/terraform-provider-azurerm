package validate

import (
	"fmt"
	"regexp"
)

func StorageEncryptionScopeName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile("^[0-9a-zA-Z]{4,63}$").MatchString(input) {
		errors = append(errors, fmt.Errorf("storage encryption scope name %q must be alphanumeric, and between 4 to 63 characters", input))
	}

	return warnings, errors
}
