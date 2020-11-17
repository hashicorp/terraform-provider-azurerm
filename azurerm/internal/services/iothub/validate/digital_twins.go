package validate

import (
	"fmt"
	"regexp"
)

func DigitaltwinsName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 3 {
		errors = append(errors, fmt.Errorf("length should be greater than %d, got %q", 3, v))
		return
	}

	if len(v) > 63 {
		errors = append(errors, fmt.Errorf("length should be less than %d, got %q", 63, v))
		return
	}

	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("expected value of %s does not match regular expression, got %v", k, v))
		return
	}
	return
}
