// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DatadogMonitorsName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_-]{2,32}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 32 characters in length, can only contain alphanumeric characters, underscore and hyphen symbols", k))
	}

	return
}
