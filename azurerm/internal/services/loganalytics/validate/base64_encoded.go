package validate

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func IsBase64Encoded(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if len(strings.TrimSpace(v)) < 1 {
		errors = append(errors, fmt.Errorf("%q must not be an empty string", k))
		return
	}

	if _, err := base64.StdEncoding.DecodeString(v); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a base64 encoded string", k))
		return
	}

	return
}
