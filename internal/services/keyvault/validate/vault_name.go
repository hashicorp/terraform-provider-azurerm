// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func VaultName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and must be between 3-24 chars", k))
	}

	if matched2 := regexp.MustCompile(`^[a-zA-Z].*[a-zA-Z0-9]$`).Match([]byte(value)); !matched2 {
		errors = append(errors, fmt.Errorf("%q must start with a letter and end with a letter or number", k))
	}

	if strings.Contains(value, "--") {
		errors = append(errors, fmt.Errorf("%q cannot contain consecutive hyphens (\"--\")", k))
	}

	return warnings, errors
}
