// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DataStoreName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]{0,254}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 1 and 255 characters, and may only include alphanumeric characters and '-'.", k))
	}
	return
}
