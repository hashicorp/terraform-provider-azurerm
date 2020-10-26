package validate

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

func LogAnalyticsGenericName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if len(v) < 4 {
		errors = append(errors, fmt.Errorf("%q length should be greater than or equal to %d characters in length", k, 4))
		return
	}

	if len(v) > 63 {
		errors = append(errors, fmt.Errorf("%q length should be less than or equal %d characters in length", k, 63))
		return
	}

	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("the %q is invalid, the %q must begin with an alphanumeric character, end with an alphanumeric character and may only contain alphanumeric characters or hyphens, got %q", k, k, v))
		return
	}
	return
}

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

	_, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a base64 encoded string", k))
		return
	}

	return
}
