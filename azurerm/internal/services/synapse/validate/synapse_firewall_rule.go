package validate

import (
	"fmt"
	"regexp"
)

func SynapseFirewallRuleName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can contain only letters, numbers, underscore and hythen.
	// 2. The value must be between 1 and 128 characters long

	if !regexp.MustCompile(`^[a-zA-Z\d-_]{1,128}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only letters, numbers, underscore and hythen, and be between 1 and 128 characters long", k))
		return
	}

	return warnings, errors
}
