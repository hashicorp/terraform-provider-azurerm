// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FunctionAppFunctionName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z](([-_0-9a-zA-Z-]{0,126})[-_0-9a-zA-Z])?$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must start with a letter, may only contain alphanumeric characters, dashes, underscore and up to 128 characters in length", key))
	}

	return warnings, errors
}
