// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func AccessConnectorName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q type to be string", k))
		return warnings, errors
	}

	// The Azure Portal shows the following validation criteria:

	// 1) Cannot be empty
	if len(v) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		// Treating this as a special case and returning early to match Azure Portal behaviour.
		return warnings, errors
	}

	// 2) Must be at least 1 characters:
	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("%q must be at least 1 character: %q", k, v))
	}

	// 3) The value must have a length of at most 64
	if len(v) > 64 {
		errors = append(errors, fmt.Errorf("%q must be no more than 64 characters: %q", k, v))
	}

	// 4) Only alphanumeric characters, underscores, and hyphens are allowed.
	if !regexp.MustCompile("^[a-zA-Z0-9_-]*$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q can contain only alphanumeric characters, underscores, and hyphens: %q", k, v))
	}

	return warnings, errors
}
