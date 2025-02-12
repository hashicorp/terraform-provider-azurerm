package validate

import (
	"fmt"
	"regexp"
)

func SSHPrivateKey(i interface{}, k string) (warnings []string, errors []error) {
	value, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}
	pattern := `(?s)-----BEGIN (RSA|DSA) PRIVATE KEY-----.*?-----END (RSA|DSA) PRIVATE KEY-----`

	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf("invalid private key, %q must be in OpenSSH PEM format", k))
	}

	return warnings, errors
}
