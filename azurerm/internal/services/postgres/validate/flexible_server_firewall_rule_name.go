package validate

import (
	"fmt"
	"regexp"
)

func FlexibleServerFirewallRuleName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", 1, v))
		return
	}

	if len(v) > 128 {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", 128, v))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must only contains numbers, characters and `-`, `_`, got %v", k, v))
		return
	}
	return
}
