package validate

import (
	"fmt"
	"regexp"
)

func ExportName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespace", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain letters and digits: %q", k, name))
	}

	if len(name) < 3 || len(name) > 24 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 3 and 24 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}
