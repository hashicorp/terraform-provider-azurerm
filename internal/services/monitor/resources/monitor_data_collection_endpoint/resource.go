// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_monitor_data_collection_endpoint"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datacollectionendpoints.ValidateDataCollectionEndpointID
}
