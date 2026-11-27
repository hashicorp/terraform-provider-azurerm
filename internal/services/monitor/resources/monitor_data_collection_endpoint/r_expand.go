// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
)

func expandCreateForDataCollectionEndpoint(state DataCollectionEndpoint) datacollectionendpoints.DataCollectionEndpointResource {
	return datacollectionendpoints.DataCollectionEndpointResource{
		Kind:     expandDataCollectionEndpointKind(state.Kind),
		Location: location.Normalize(state.Location),
		Name:     pointer.To(state.Name),
		Properties: &datacollectionendpoints.DataCollectionEndpoint{
			Description: pointer.To(state.Description),
			NetworkAcls: &datacollectionendpoints.NetworkRuleSet{
				PublicNetworkAccess: expandDataCollectionEndpointPublicNetworkAccess(state.PublicNetworkAccessEnabled),
			},
		},
		Tags: tags.Expand(state.Tags),
	}
}
