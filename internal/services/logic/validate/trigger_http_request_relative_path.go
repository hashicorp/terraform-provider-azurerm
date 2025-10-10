// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func TriggerHttpRequestRelativePath(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9_/}{]+$").MatchString(value) {
		errors = append(errors, fmt.Errorf("%s can only contain alphanumeric characters, underscores, forward slashes and curly braces", k))
	}

	return warnings, errors
}
