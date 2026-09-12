// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageDiscoveryWorkspaceName() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(4, 64),
		validation.StringMatch(
			regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9]|-[a-zA-Z0-9]){2,62}[a-zA-Z0-9]$`),
			"can only contain letters, hyphens, and numbers; hyphens cannot be consecutive; cannot start or end with a hyphen; must start with a letter",
		),
	)
}
