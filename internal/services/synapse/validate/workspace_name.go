package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func WorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can contain only lowercase letters, numbers or hyphens
	// 2. must start and end with a lowercase letter or number
	// 3. must not contain the string '-ondemand'
	// 4. The value must be between 1 and 50 characters long

	if !regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,48}[a-z0-9])?$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must start and end with a letter or number, can contain only lowercase letters, numbers or hyphens, and be between 1 and 50 characters long", k))
		return
	}
	if strings.Contains(v, "-ondemand") {
		errors = append(errors, fmt.Errorf("%s must not contain the string '-ondemand'", k))
		return
	}
	return warnings, errors
}
