// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ComputeClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The portal says: It can include letters, digits and dashes. It must start with a letter, end with a letter or digit, and be between 3 and 32 characters in length.
	// If you provide invalid name, the rest api will return an error with the following regex.
	if matched := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{1,30}[a-zA-Z0-9]$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 3 and 32 characters, may only include alphanumeric characters and '-' and must start with a letter, end with a letter or digit", k))
	}
	return
}
