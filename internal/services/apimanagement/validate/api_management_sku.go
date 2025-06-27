// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApimSkuName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^Consumption_0$|^Basic_(1|2)$|^BasicV2_([1-9]|10)|^Developer_1$|^Premium_([1-9][0-9]{0,1})$|^PremiumV2_([1-9]|1[0-9]|2[0-9]|30)$|^Standard_[1-4]$|^StandardV2_([1-9]|10)$`),
		`This is not a valid Api Management sku name.`,
	)
}
