// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"unicode"
)

func ResourceAnchorName(i interface{}, k string) (warnings []string, errors []error) {
	var validationErrors []error

	v, ok := i.(string)
	if !ok {
		validationErrors = append(validationErrors, fmt.Errorf("expected type of %s to be string", k))
		return []string{}, validationErrors
	}

	hasInvalidChar := false
	for _, r := range v {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			hasInvalidChar = true
			break
		}
	}
	if hasInvalidChar {
		validationErrors = append(validationErrors, fmt.Errorf("%v must contain only letters and numbers", k))
	}

	if len(v) > 24 {
		validationErrors = append(validationErrors, fmt.Errorf("%v must be 24 characters max", k))
	}

	return []string{}, validationErrors
}
