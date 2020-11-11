package validate

import (
	"fmt"
	"strings"
)

func SubscriptionAliasName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The value must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	const maxLength = 128
	// subscription alias name can be 1-128 characters in length
	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
	}

	// subscription alias name cannot contain the following characters: '/', '\\', '?','#'
	if strings.ContainsAny(v, "/\\?#") {
		errors = append(errors, fmt.Errorf("%q cannot contain '/', '\\', '?' or '#'", k))
	}

	return
}
