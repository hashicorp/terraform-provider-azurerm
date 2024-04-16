// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func WorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)

	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't", k))
		return
	}

	// The value must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	const minLength = 3
	const maxLength = 64

	// Workspace name can be 3-64 characters in length
	if len(v) > maxLength || len(v) < minLength {
		errors = append(errors, fmt.Errorf("%q must be between %d-%d characters, got %d", k, minLength, maxLength, len(v)))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, dots, dashes and underscores", k))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9]`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must begin with an alphanumeric character", k))
	}

	if matched := regexp.MustCompile(`\w$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must end with an alphanumeric character or underscore", k))
	}

	return warnings, errors
}
