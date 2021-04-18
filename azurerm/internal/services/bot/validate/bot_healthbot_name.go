package validate

import (
	"fmt"
	"regexp"
)

func HealthbotName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if len(v) < 2 {
		errors = append(errors, fmt.Errorf("length should be greater than %d", 2))
		return
	}
	if len(v) > 64 {
		errors = append(errors, fmt.Errorf("length should be less than %d", 64))
		return
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("expected value of %s not match regular expression, got %v", k, v))
		return
	}
	return
}
