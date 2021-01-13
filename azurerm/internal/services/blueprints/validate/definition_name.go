package validate

import (
	"fmt"
	"regexp"
)

func DefinitionName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^[A-Za-z0-9-_]{1,48}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s can include letters, numbers, underscores or dashes. Spaces and other special characters are not allowed.", k))
	}

	return warnings, errors
}
