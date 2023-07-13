// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

const (
	minLength = 1
	maxLength = 63
)

func FlexibleServerDatabaseName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < minLength {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", minLength, v))
		return
	}

	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", maxLength, v))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z-_]`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter got %v", k, v))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must only contains numbers, characters and `-`, `_`, got %v", k, v))
		return
	}

	return
}
