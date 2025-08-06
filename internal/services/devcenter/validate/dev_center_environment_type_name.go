// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DevCenterEnvironmentTypeName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-_.]{2,62}$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must start with an alphanumeric character, may contain alphanumeric characters, dashes, underscores or periods and must be between 3 and 63 characters long", k))
	}

	return warnings, errors
}
