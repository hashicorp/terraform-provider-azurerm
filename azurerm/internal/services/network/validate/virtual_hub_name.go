package validate

import (
	"fmt"
	"regexp"
)

func VirtualHubName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^.{1,256}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 256 characters in length.", k))
	}

	return warnings, errors
}
