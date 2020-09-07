package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func StorageShareDirectoryName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// File share names can contain only uppercase and lowercase letters, numbers, and hyphens,
	// However they can be nested (e.g. foo/bar) with at most one level.
	if !regexp.MustCompile(`^[A-Za-z0-9-]+(/[A-Za-z0-9-]+)?$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must contain only uppercase and lowercase alphanumeric characters, numbers and hyphens, and can be nested one level", k))
	}

	// The name must begin and end with a letter or a number.
	start := regexp.MustCompile(`^[a-zA-Z]`)
	end := regexp.MustCompile(`[a-zA-Z0-9]$`)
	parts := strings.Split(value, "/")
	for _, p := range parts {
		if !start.MatchString(p) || !end.MatchString(p) {
			errors = append(errors, fmt.Errorf("%s must start and end with a letter and end only with a number or letter", k))
			break
		}
	}

	return warnings, errors
}
