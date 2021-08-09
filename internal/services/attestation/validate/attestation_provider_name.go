package validate

import (
	"fmt"
	"regexp"
)

func AttestationProviderName(i interface{}, k string) (warning []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if !regexp.MustCompile(`^[a-z\d]{3,24}\z`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must be between 3 and 24 characters in length and use numbers and lower-case letters only.", k))
	}

	return
}
