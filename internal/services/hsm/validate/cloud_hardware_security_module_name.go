// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ValidateCloudHsmClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// A Cloud HSM name must be between 3-23 alphanumeric characters.
	// The name must begin with a letter or digit, end with a letter or digit,
	// and not contain consecutive hyphens
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,21}[a-zA-Z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 23 alphanumeric characters. It must begin with a letter or digit, end with a letter or digit", k))
		return
	}

	// No consecutive hyphens
	if regexp.MustCompile("(--)").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must not contain any consecutive hyphens", k))
	}

	return
}
