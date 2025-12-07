// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FileSystemName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,13}[a-zA-Z0-9]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("`%q` must be between 2 and 15 characters in length, must not begin or end with a hyphen and may only contain alphanumeric characters and hyphens", k))
	}

	return warnings, errors
}
