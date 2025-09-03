// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func LogAnalyticsWorkspaceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$").MatchString(value) {
		errors = append(errors, fmt.Errorf("%s can only contain alphabet, number, and '-' character. You can not use '-' as the start and end of the name", k))
	}

	length := len(value)
	if length > 63 || 4 > length {
		errors = append(errors, fmt.Errorf("%s can only be between 4 and 63 letters", k))
	}

	return warnings, errors
}
