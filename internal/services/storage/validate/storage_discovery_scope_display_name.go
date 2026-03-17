// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// StorageDiscoveryScopeDisplayName validates the scope display_name per Azure portal:
// 4-64 characters; letters, numbers, spaces, and hyphens only; cannot start or end with
// a number, space, or hyphen; no consecutive spaces or hyphens.
func StorageDiscoveryScopeDisplayName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if len(input) < 4 || len(input) > 64 {
		errors = append(errors, fmt.Errorf("scope display name must be between 4 and 64 characters"))
		return warnings, errors
	}

	if strings.HasPrefix(input, " ") || strings.HasSuffix(input, " ") {
		errors = append(errors, fmt.Errorf("scope display name cannot start or end with a space"))
		return warnings, errors
	}

	if strings.HasPrefix(input, "-") || strings.HasSuffix(input, "-") {
		errors = append(errors, fmt.Errorf("scope display name cannot start or end with a hyphen"))
		return warnings, errors
	}

	if regexp.MustCompile(`^[0-9]`).MatchString(input) || regexp.MustCompile(`[0-9]$`).MatchString(input) {
		errors = append(errors, fmt.Errorf("scope display name cannot start or end with a number"))
		return warnings, errors
	}

	if strings.Contains(input, "  ") || strings.Contains(input, "--") {
		errors = append(errors, fmt.Errorf("scope display name cannot contain consecutive spaces or hyphens"))
		return warnings, errors
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\- ]+$`).MatchString(input) {
		errors = append(errors, fmt.Errorf("scope display name can only contain letters, numbers, spaces, and hyphens"))
		return warnings, errors
	}

	return warnings, errors
}
