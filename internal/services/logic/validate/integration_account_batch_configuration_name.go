// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IntegrationAccountBatchConfigurationName() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(0, 20),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9]+$`), "contains only letters and numbers"),
	)
}
