// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func EndpointName(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return warnings, errors
	}

	if !regexp.MustCompile(`^[0-9a-zA-Z][-_0-9a-zA-Z]{0,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length, begin with a letter or number, and may contain letters, numbers, dashes and underscore", k))
	}

	return warnings, errors
}
