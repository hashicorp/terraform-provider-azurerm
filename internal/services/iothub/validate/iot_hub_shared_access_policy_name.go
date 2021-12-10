package validate

import (
	"fmt"
	"regexp"
)

func IotHubSharedAccessPolicyName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules are :
	// 1. must not be empty.
	// 2. must not exceed 64 characters in length.
	// 3. can only contain alphanumeric characters, exclamation marks, periods, underscores and hyphens

	if !regexp.MustCompile(`[a-zA-Z0-9!._-]{1,64}`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must not be empty, and must not exceed 64 characters in length, and can only contain alphanumeric characters, exclamation marks, periods, underscores and hyphens", k))
	}

	return nil, errors
}
