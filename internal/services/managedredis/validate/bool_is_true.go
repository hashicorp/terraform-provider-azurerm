// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "fmt"

func BoolIsTrue(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(bool)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be bool", k))
		return warnings, errors
	}

	if !v {
		errors = append(errors, fmt.Errorf("expected %s to be true, got %t", k, v))
		return warnings, errors
	}

	return warnings, errors
}
