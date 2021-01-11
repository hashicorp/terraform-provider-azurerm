package validate

import (
	"fmt"
	"regexp"
)

func KeyVaultCertificateIssuerName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	return warnings, errors
}
