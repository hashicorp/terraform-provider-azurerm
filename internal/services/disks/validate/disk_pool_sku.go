// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DiskPoolSku() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice(
		[]string{
			"Basic_B1",
			"Standard_S1",
			"Premium_P1",
		}, false,
	)
}
