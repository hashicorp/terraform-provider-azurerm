// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_connectivity_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.ConnectivityConfigurations

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			state := ManagerConnectivityConfigurationModel{
				Name:                                id.ConnectivityConfigurationName,
				NetworkManagerId:                    connectivityconfigurations.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
				AppliesToGroups:                     flattenConnectivityGroupItemModel(properties.AppliesToGroups),
				ConnectivityTopology:                properties.ConnectivityTopology,
				DeleteExistingPeeringEnabled:        flattenDeleteExistingPeering(properties.DeleteExistingPeering),
				GlobalMeshEnabled:                   flattenConnectivityConfIsGlobal(properties.IsGlobal),
				Hub:                                 flattenHubModel(properties.Hubs),
				ConnectedGroupAddressOverlapEnabled: true,
				ConnectedGroupPrivateEndpointsScale: string(connectivityconfigurations.ConnectedGroupPrivateEndpointsScaleStandard),
			}

			if properties.ConnectivityCapabilities != nil {
				state.ConnectedGroupAddressOverlapEnabled = flattenConnectedGroupAddressOverlap(properties.ConnectivityCapabilities.ConnectedGroupAddressOverlap)
				state.ConnectedGroupPrivateEndpointsScale = string(properties.ConnectivityCapabilities.ConnectedGroupPrivateEndpointsScale)
				state.PeeringEnforcementEnabled = flattenPeeringEnforcement(properties.ConnectivityCapabilities.PeeringEnforcement)
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}
