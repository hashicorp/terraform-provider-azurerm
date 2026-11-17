// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package extended_location_custom_location

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = Resource{}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_extended_location_custom_location"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return customlocations.ValidateCustomLocationID
}
