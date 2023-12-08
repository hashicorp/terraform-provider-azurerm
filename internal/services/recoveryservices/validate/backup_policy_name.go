// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func BackupPolicyName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-\da-zA-Z]{2,149}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 150 characters in length and start with letters and contains only letters, numbers or hyphens.", k))
	}

	return warnings, errors
}
