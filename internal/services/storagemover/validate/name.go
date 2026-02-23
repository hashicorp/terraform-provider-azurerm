// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// StorageMoverResourceName validates Storage Mover resource names per the REST API pattern.
// Pattern: ^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$ (1-64 chars, start with alphanumeric, letters/numbers/underscore/hyphen)
func StorageMoverResourceName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`).MatchString(input) {
		errors = append(errors, fmt.Errorf("name must be between 1 and 64 characters in length, begin with a letter or number, and may contain letters, numbers, dashes and underscores"))
	}

	return warnings, errors
}
