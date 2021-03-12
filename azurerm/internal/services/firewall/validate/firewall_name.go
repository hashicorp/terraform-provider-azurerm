package validate

import (
	"fmt"
	"regexp"
)

func FirewallName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// From the Portal:
	// The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.
	if matched := regexp.MustCompile(`^[0-9a-zA-Z]([0-9a-zA-Z._-]{0,}[0-9a-zA-Z_])?$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.", k))
	}

	return warnings, errors
}
