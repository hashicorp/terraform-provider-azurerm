// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IntegrationAccountPartnerName() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(0, 80),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-().]+$`), "contains only letters, numbers, dots, parentheses and hyphens"),
	)
}
