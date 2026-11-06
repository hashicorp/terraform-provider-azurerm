// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_fleet

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_compute_fleet"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleets.ValidateFleetID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &fleets.FleetId{}
}
