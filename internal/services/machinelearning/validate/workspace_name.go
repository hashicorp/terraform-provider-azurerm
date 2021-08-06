package validate

import (
	"fmt"
	"regexp"
)

func WorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The portal says: The workspace name must be between 3 and 33 characters. The name may only include alphanumeric characters and '-'.
	// If you provide invalid name, the rest api will return an error with the following regex.
	if matched := regexp.MustCompile(`^[a-zA-Z0-9][\w-]{2,32}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 3 and 33 characters, and may only include alphanumeric characters and '-'.", k))
	}
	return
}
