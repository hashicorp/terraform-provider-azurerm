// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_connectivity_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.ConnectivityConfigurations

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			if metadata.ResourceData.HasChange("applies_to_group") {
				properties.AppliesToGroups = expandConnectivityGroupItemModel(model.AppliesToGroups)
			}

			if metadata.ResourceData.HasChange("connected_group_address_overlap_enabled") {
				if properties.ConnectivityCapabilities == nil {
					properties.ConnectivityCapabilities = &connectivityconfigurations.ConnectivityConfigurationPropertiesConnectivityCapabilities{}
				}
				properties.ConnectivityCapabilities.ConnectedGroupAddressOverlap = expandConnectedGroupAddressOverlap(model.ConnectedGroupAddressOverlapEnabled)
			}

			if metadata.ResourceData.HasChange("connected_group_private_endpoints_scale") {
				if properties.ConnectivityCapabilities == nil {
					properties.ConnectivityCapabilities = &connectivityconfigurations.ConnectivityConfigurationPropertiesConnectivityCapabilities{}
				}
				properties.ConnectivityCapabilities.ConnectedGroupPrivateEndpointsScale = connectivityconfigurations.ConnectedGroupPrivateEndpointsScale(model.ConnectedGroupPrivateEndpointsScale)
			}

			if metadata.ResourceData.HasChange("peering_enforcement_enabled") {
				if properties.ConnectivityCapabilities == nil {
					properties.ConnectivityCapabilities = &connectivityconfigurations.ConnectivityConfigurationPropertiesConnectivityCapabilities{}
				}
				properties.ConnectivityCapabilities.PeeringEnforcement = expandPeeringEnforcement(model.PeeringEnforcementEnabled)
			}

			if metadata.ResourceData.HasChange("connectivity_topology") {
				properties.ConnectivityTopology = model.ConnectivityTopology
			}

			if metadata.ResourceData.HasChange("delete_existing_peering_enabled") {
				properties.DeleteExistingPeering = expandDeleteExistingPeering(model.DeleteExistingPeeringEnabled)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("hub") {
				properties.Hubs = expandHubModel(model.Hub)
			}

			if metadata.ResourceData.HasChange("global_mesh_enabled") {
				properties.IsGlobal = expandConnectivityConfIsGlobal(model.GlobalMeshEnabled)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
