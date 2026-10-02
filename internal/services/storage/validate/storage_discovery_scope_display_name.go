// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageDiscoveryScopeDisplayName() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(4, 64),
		validation.StringMatch(
			regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9 -]*[a-zA-Z]$`),
			"must start and end with a letter and can only contain letters, numbers, spaces, and hyphens",
		),
		validation.StringDoesNotMatch(
			regexp.MustCompile(`  |--`),
			"cannot contain consecutive spaces or hyphens",
		),
	)
}
