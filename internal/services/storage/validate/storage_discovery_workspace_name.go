// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// StorageDiscoveryWorkspaceName validates the name of a Storage Discovery Workspace
// Name must be between 3 and 63 characters, can only contain lowercase letters, numbers, and hyphens,
// and must start and end with a letter or number.
func StorageDiscoveryWorkspaceName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`).MatchString(input) {
		errors = append(errors, fmt.Errorf("name (%q) can only consist of lowercase letters, numbers, and hyphens, must be between 3 and 63 characters long, and must start and end with a letter or number", input))
	}

	return warnings, errors
}
