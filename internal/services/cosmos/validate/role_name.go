// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func RoleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-z0-9]{1,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 63 characters in length and only contain lower case letters and numbers", k))
	}

	return warnings, errors
}
