package validate

import (
	"fmt"
	"strings"
)

func PolicySetDefinitionReferenceID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The service returns error when definition reference id is too long
	// error: The policy definition reference id must not exceed '128' characters.
	// By my additional test, the definition reference id cannot contain the following characters: %^#/\?.
	if len(v) > 128 {
		errors = append(errors, fmt.Errorf("%s must not exceed '128' characters", k))
		return
	}
	const invalidCharacters = `%^#/\?`
	if strings.ContainsAny(v, invalidCharacters) {
		errors = append(errors, fmt.Errorf("%s cannot contain the following characters: %s", k, invalidCharacters))
	}

	return warnings, errors
}
