package validate

import (
	"fmt"
	"regexp"
)

func CustomIpPrefixName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if !regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_])$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods or hyphens", k))
	}

	return nil, errors
}
