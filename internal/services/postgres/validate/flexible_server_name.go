package validate

import (
	"fmt"
	"regexp"
)

func FlexibleServerName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", 3, v))
		return
	}

	if len(v) > 63 {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", 63, v))
		return
	}

	if !regexp.MustCompile(`^[a-z0-9]([a-z0-9-]+[a-z0-9])?$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must only contains numbers, lowercase characters and '-', got %v", k, v))
		return
	}
	return
}
