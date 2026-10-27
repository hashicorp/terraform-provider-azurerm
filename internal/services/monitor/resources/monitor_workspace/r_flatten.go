// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_workspace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *azuremonitorworkspaces.AccountId, model *azuremonitorworkspaces.AzureMonitorWorkspaceResource) error {
	state := WorkspaceResourceModel{
		Name:              id.AccountName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Tags = pointer.From(model.Tags)
		state.Location = location.Normalize(model.Location)

		if properties := model.Properties; properties != nil {
			publicNetworkAccess := true
			if properties.PublicNetworkAccess != nil {
				publicNetworkAccess = azuremonitorworkspaces.PublicNetworkAccessEnabled == *properties.PublicNetworkAccess
			}
			state.PublicNetworkAccessEnabled = publicNetworkAccess

			if properties.Metrics != nil && properties.Metrics.PrometheusQueryEndpoint != nil {
				state.QueryEndpoint = *properties.Metrics.PrometheusQueryEndpoint
			}

			if properties.DefaultIngestionSettings != nil {
				if properties.DefaultIngestionSettings.DataCollectionEndpointResourceId != nil {
					state.DefaultDataCollectionEndpointId = *properties.DefaultIngestionSettings.DataCollectionEndpointResourceId
				}
				if properties.DefaultIngestionSettings.DataCollectionRuleResourceId != nil {
					state.DefaultDataCollectionRuleId = *properties.DefaultIngestionSettings.DataCollectionRuleResourceId
				}
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
