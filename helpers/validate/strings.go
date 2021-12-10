package validate

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Base64EncodedString validates that the string is base64 encoded
func Base64EncodedString(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("%q must not be empty", k)}
	}

	if _, err := base64.StdEncoding.DecodeString(v); err != nil {
		return nil, []error{fmt.Errorf("%q must be a valid base64 encoded string", k)}
	}

	return nil, nil
}

// LowerCasedString validates that the string is lower-cased
func LowerCasedString(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("%q must not be empty", k)}
	}

	if strings.ToLower(v) != v {
		return nil, []error{fmt.Errorf("%q must be a lower-cased string", k)}
	}

	if strings.ContainsAny(v, " ") {
		return nil, []error{fmt.Errorf("%q cannot contain whitespace", k)}
	}

	return nil, nil
}
