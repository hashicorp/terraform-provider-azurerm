// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fabric_capacity

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fabriccapacities.ValidateCapacityID
}

func (r Resource) ResourceType() string {
	return "azurerm_fabric_capacity"
}
