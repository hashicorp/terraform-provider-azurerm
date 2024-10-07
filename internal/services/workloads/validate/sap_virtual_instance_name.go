// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SAPVirtualInstanceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[A-Z][A-Z0-9][A-Z0-9]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must have three characters. The first character must be alphabetic. The second and third characters must be alphanumerical. All alphabetic characters should be uppercase.", k))
		return warnings, errors
	}

	return warnings, errors
}
