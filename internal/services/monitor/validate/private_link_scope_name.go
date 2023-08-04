// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func PrivateLinkScopeName(i interface{}, k string) (warning []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("length should be greater than %d", 1))
		return
	}

	if len(v) > 255 {
		errors = append(errors, fmt.Errorf("length should be less and equal than %d", 255))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9()-_.]*[a-zA-Z0-9_-]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q only allows alphanumeric characters, periods, underscores, hyphens and parenthesis and cannot end in a period", k))
		return
	}

	return
}
