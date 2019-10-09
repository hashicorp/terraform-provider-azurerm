package validate

import (
	"fmt"
	"strings"
)

// NoEmptyStrings validates that the string is not just whitespace characters (equal to [\r\n\t\f\v ])
func NoEmptyStrings(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("%q must not be empty", k)}
	}

	return nil, nil
}

// PrivateLinkEnpointRequestMessage validates that the Private Link Enpoint Request Message is less than 140 characters
func PrivateLinkEnpointRequestMessage(i interface{}, k string) (_ []string, errors []error) {
	return stringMaxLength(140)(i, k)
}

func stringMaxLength(maxLength int) func(i interface{}, k string) (_ []string, errors []error) {
	return func(i interface{}, k string) (_ []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
		}

		if len(v) > maxLength {
			return nil, []error{fmt.Errorf("%q must not be longer than %d characters, got %d", k, maxLength, len(v))}
		}

		if strings.TrimSpace(v) == "" {
			return nil, []error{fmt.Errorf("%q must not be empty", k)}
		}

		return
	}
}
