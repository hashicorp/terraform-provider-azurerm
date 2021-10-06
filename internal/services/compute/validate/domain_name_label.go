package validate

import (
	"fmt"
	"regexp"
)

func OrchestratedDomainNameLabel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^[a-z][a-z0-9-]{1,61}[a-z0-9]$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 1 - 63 characters long, start with a lower case letter, end with a lower case letter or number and contains only a-z, 0-9 and hyphens", k))
	}
	return
}
