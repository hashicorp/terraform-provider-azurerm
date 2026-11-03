// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_virtual_network_appliance

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkvirtualappliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_palo_alto_virtual_network_appliance"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkvirtualappliances.ValidateNetworkVirtualApplianceID
}
