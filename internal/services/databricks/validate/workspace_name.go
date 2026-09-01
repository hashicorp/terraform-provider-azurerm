// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func WorkspaceName(i interface{}, k string) ([]string, []error) {
	// NOTE: Portal limits the name to 30 characters but the API allows up to 64. In order to facilitate imports of
	// workspaces created through the API, use the higher limit. Note that once the workspace is created, the portal
	// handles workspace names >30 characters just fine
	return validation.All(
		validation.StringLenBetween(3, 64),
		validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9_-]*$"), "can contain only alphanumeric characters, underscores, and hyphens"),
	)(i, k)
}
