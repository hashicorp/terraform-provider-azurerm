// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func LabUsername(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^.{1,20}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only be up to 20 characters in length", k))
	}

	return warnings, errors
}
