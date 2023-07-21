// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FirewallRuleName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. Can't contain '<,>,*,%,&,:,\,/,?'.
	// 2. Can't end with '.'
	// 2. The value must be between 1 and 128 characters long

	if !regexp.MustCompile(`^[^<>*%&:\\/?]{0,127}[^.<>*%&:\\/?]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can't contain '<,>,*,%%,&,:,\\,/,?', can't end with '.', and must be between 1 and 128 characters long", k))
		return
	}

	return warnings, errors
}
