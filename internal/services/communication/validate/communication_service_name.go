package validate

import (
	"fmt"
	"regexp"
)

func CommunicationServiceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !regexp.MustCompile(`^(([a-zA-Z])|([a-zA-Z][0-9a-zA-Z-]{0,62}[0-9a-zA-Z]))$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length and start with letters and contain only letters, numbers and hyphens. And it cannot end with hyphen.", k))
		return
	}

	return
}
