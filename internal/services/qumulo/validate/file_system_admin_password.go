// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ValidatePasswordComplexity(input interface{}, key string) (warnings []string, errors []error) {
	password, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected `%q` to be a string", key))
		return warnings, errors
	}

	complexityMatch := 0
	re := regexp.MustCompile(`[a-z]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[A-Z]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[0-9]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[\W_]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	if complexityMatch < 3 {
		errors = append(errors, fmt.Errorf("`%q` did not meet minimum password complexity requirements. A password must contain at least 3 of the 4 following conditions: a lower case character, a upper case character, a digit and/or a special character. Got %q", key, password))
		return warnings, errors
	}

	if len(password) < 8 {
		errors = append(errors, fmt.Errorf("`%q` must be at least 8 characters long. Got %q(%d characters)", key, password, len(password)))
		return warnings, errors
	}

	return warnings, errors
}
