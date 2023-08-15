// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func BotName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 4 {
		errors = append(errors, fmt.Errorf("length should be greater than %d", 4))
		return
	}

	if len(v) > 42 {
		errors = append(errors, fmt.Errorf("length should be less than %d", 42))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must start with a letter or digit and may only contain alphanumeric characters, underscores and dashes", k))
		return
	}

	return
}
