// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func StringDoesNotContainUpperCaseLetter(input interface{}, k string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if strings.ToLower(v) != v {
		errors = append(errors, fmt.Errorf("expected value of %s to not contain any uppercase letter", k))
		return warnings, errors
	}

	return warnings, errors
}
