// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func CostAnomalyAlertName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules (as far as we know) are :
	// 1. can contain only lowercase letters, numbers and hyphens.

	if !regexp.MustCompile(`^([a-z\d-]*)$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must contain only lowercase letters, numbers and hyphens.", k))
	}

	return nil, errors
}
