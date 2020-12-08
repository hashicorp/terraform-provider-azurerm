package validate

import (
	"fmt"
	"regexp"
)

func ActionRuleName(i interface{}, k string) (warning []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if !regexp.MustCompile(`^([a-zA-Z\d])[a-zA-Z\d-_]*$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s should begin with a letter or number, contain only letters, numbers, underscores and hyphens.", k))
	}

	return
}
