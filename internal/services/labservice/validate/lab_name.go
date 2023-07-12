// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func LabName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[0-9a-zA-Z-_.]{1,60}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, periods, underscores and hyphens", k))
	}

	return warnings, errors
}
