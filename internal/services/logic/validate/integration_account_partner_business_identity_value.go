// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IntegrationAccountPartnerBusinessIdentityValue() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(0, 128),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-() ._]+$`), "contains only letters, numbers, dots, parentheses, hyphens and underscores"),
	)
}
