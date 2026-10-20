// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_resource_bridge_appliance

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_arc_resource_bridge_appliance"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appliances.ValidateApplianceID
}
