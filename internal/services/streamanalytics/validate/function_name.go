// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FunctionName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z0-9-]{3,63}$`).MatchString(input) {
		errors = append(errors, fmt.Errorf("%s contain only letters, numbers and hyphens. The value must be between 3 and 63 characters long", input))
	}

	return warnings, errors
}
