// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// StorageDiscoveryWorkspaceName validates the name of a Storage Discovery Workspace
// consistent with the Azure portal. Rules: 4-64 chars; letters, hyphens, and numbers only;
// no consecutive hyphens; cannot start or end with hyphen; must start with a letter.
func StorageDiscoveryWorkspaceName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if len(input) < 4 || len(input) > 64 {
		errors = append(errors, fmt.Errorf("workspace name must be 4-64 characters long"))
		return warnings, errors
	}

	// Must start with letter, end with letter/number, no consecutive hyphens
	if !regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9]|-[a-zA-Z0-9]){2,62}[a-zA-Z0-9]$`).MatchString(input) {
		errors = append(errors, fmt.Errorf("workspace name can only contain letters, hyphens, and numbers; hyphens cannot be consecutive; cannot start or end with a hyphen; must start with a letter"))
		return warnings, errors
	}

	if strings.Contains(input, "--") {
		errors = append(errors, fmt.Errorf("workspace name cannot contain consecutive hyphens"))
		return warnings, errors
	}

	return warnings, errors
}
