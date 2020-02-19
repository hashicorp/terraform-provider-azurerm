package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func StorageShareDirectoryName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// File share names can contain only uppercase and lowercase letters, numbers, and hyphens,
	// and must begin and end with a letter or a number.
	// however they can be nested (e.g. foo/bar)
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`).MatchString(value) && !regexp.MustCompile(`^[A-Za-z0-9]{1,}/[A-Za-z0-9]{1,}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must contain only uppercase and lowercase alphanumeric characters, numbers and hyphens. It must start and end with a letter and end only with a number or letter", k))
	}

	// The name cannot contain two consecutive hyphens.
	if strings.Contains(value, "--") {
		errors = append(errors, fmt.Errorf("%s cannot contain two concecutive hyphens", k))
	}

	return warnings, errors
}
