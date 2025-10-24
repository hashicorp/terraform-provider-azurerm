// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"unicode"
)

func ResourceAnchorName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	for _, r := range v {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return []string{}, append(errors, fmt.Errorf("%v must contain only letters and numbers", k))
		}
	}

	if len(v) > 24 {
		return []string{}, append(errors, fmt.Errorf("%v must be 24 characters max", k))
	}

	return []string{}, []error{}
}
