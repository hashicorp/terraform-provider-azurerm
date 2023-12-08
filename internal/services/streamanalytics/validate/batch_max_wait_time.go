// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func BatchMaxWaitTime(input interface{}, key string) (warnings []string, errors []error) {
	value, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", key))
		return
	}

	if value == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", key))
	}

	if matched := regexp.MustCompile(`[0-9]{2}:[0-9]{2}:[0-9]{2}`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must have the following format hh:mm:ss", key))
	}

	return warnings, errors
}
