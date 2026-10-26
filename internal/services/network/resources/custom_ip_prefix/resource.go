// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_ip_prefix

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = Resource{}

type Resource struct {
	client *customipprefixes.CustomIPPrefixesClient
}

func (Resource) ResourceType() string {
	return "azurerm_custom_ip_prefix"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return customipprefixes.ValidateCustomIPPrefixID
}
