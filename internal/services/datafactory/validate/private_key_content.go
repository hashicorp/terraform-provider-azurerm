package validate

import (
	"fmt"
	"strings"
)

func PrivateSSHKey(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(value, "-----BEGIN") {
		errors = append(errors, fmt.Errorf("private SSH Key must be in OpenSSH format %q: %q", k, value))
	}

	return warnings, errors
}
