// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func BackupPolicyMySQLFlexibleServerName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be 3 - 150 characters long, contain only letters, numbers and hyphens", k))
	}

	return warnings, errors
}
