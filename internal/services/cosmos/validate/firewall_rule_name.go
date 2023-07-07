// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FirewallRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,}[a-zA-Z0-9_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must consist of letters, digits, underscores, periods and hyphens. The first character must be a letter or digit, and the last character must be a letter, a digit or an underscore", k))
	}

	return warnings, errors
}
