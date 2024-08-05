// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func NetworkInterfaceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z][0-9a-zA-Z_.-]{0,62}[0-9a-zA-Z_]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 64 characters in length. It must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods or hyphens.", k))
	}

	return warnings, errors
}
