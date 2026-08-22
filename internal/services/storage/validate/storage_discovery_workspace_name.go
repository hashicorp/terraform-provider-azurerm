// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	stderrors "errors"
	"regexp"
	"strings"
)

var storageDiscoveryWorkspaceNamePattern = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9]|-[a-zA-Z0-9]){2,62}[a-zA-Z0-9]$`)

func StorageDiscoveryWorkspaceName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if len(input) < 4 || len(input) > 64 {
		errors = append(errors, stderrors.New("workspace name must be 4-64 characters long"))
		return warnings, errors
	}

	if !storageDiscoveryWorkspaceNamePattern.MatchString(input) {
		errors = append(errors, stderrors.New("workspace name can only contain letters, hyphens, and numbers; hyphens cannot be consecutive; cannot start or end with a hyphen; must start with a letter"))
		return warnings, errors
	}

	if strings.Contains(input, "--") {
		errors = append(errors, stderrors.New("workspace name cannot contain consecutive hyphens"))
		return warnings, errors
	}

	return warnings, errors
}
