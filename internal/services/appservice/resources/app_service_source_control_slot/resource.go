// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control_slot

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_app_service_source_control_slot" // TODO - Does this name fit the new convention?
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	// This is a meta resource with a 1:1 relationship with the slot it's pointed at, so we use the same ID
	return webapps.ValidateSlotID
}
