// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FleetName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)
	if len(name) < 3 || len(name) > 50 {
		errors = append(errors, fmt.Errorf("length of %q must be between 3 to 50 (inclusive)", k))
	}

	if !regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must consist of lower case letters, digits, and hyphens. The first and last character must be a letter or digit", k))
	}

	return warnings, errors
}
