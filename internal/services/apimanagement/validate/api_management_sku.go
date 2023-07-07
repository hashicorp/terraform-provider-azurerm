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
		regexp.MustCompile(`^Consumption_0$|^Basic_(1|2)$|^Developer_1$|^Premium_([1-9][0-9]{0,1})$|^Standard_[1-4]$`),
		`This is not a valid Api Management sku name.`,
	)
}
