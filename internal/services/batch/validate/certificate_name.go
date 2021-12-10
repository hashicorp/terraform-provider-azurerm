package validate

import (
	"fmt"
	"regexp"
)

func CertificateName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w]+-[\w]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"must be made up of algorithm and thumbprint separated by a dash in %q: %q", k, value))
	}

	return warnings, errors
}
