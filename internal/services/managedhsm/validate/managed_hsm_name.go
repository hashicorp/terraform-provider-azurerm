// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ManagedHardwareSecurityModuleName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return warnings, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules are :
	// Must be a 3-24 character string, containing only 0-9, a-z. A-Z, and -
	// The name must begin with a letter, end with a letter or digit, and not contain consecutive hyphens.
	if strings.Contains(v, "--") || !regexp.MustCompile(`^[a-zA-Z][-a-zA-Z\d]{1,22}[a-zA-Z\d]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must begin with a letter, end with a letter or number, contain only alphanumeric characters. The value must be between 3 and 24 characters long", k))
	}

	return warnings, errors
}
