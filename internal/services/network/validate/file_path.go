// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FilePath(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^(.)+.cap$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must end with extension name '.cap'", k))
	}

	return warnings, errors
}
