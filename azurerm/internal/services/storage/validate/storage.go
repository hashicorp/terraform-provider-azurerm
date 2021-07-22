package validate

import (
	"fmt"
	"regexp"
)

func StorageShareDirectoryName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[A-Za-z0-9\-_]+(/[A-Za-z0-9\-_]+)*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must contain only uppercase and lowercase alphanumeric characters, numbers, hyphens and underscores, and can be nested multiple levels", k))
	}

	return warnings, errors
}
