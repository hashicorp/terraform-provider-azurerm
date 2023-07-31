// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
)

func IntegerPositive(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", i))
		return
	}
	if v <= 0 {
		errors = append(errors, fmt.Errorf("expected %s to be positive, got %d", k, v))
		return
	}
	return
}
