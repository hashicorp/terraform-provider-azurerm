// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SqlPoolName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can't contain '<,>,*,%,&,:,\,/,?,@,-' or control characters
	// 2. can't end with '.' or ' '
	// 3. The value must be between 1 and 60 characters long
	if !regexp.MustCompile(`^[^<>*%&:\\\/?@-]{0,59}[^\s.<>*%&:\\\/?@-]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only letters, numbers or underscore, The value must be between 1 and 60 characters long", k))
		return
	}

	return warnings, errors
}
