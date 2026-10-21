// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_connectivity_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.ConnectivityConfigurations
			networkManagerId, err := connectivityconfigurations.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := connectivityconfigurations.NewConnectivityConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			conf := connectivityconfigurations.ConnectivityConfiguration{
				Properties: &connectivityconfigurations.ConnectivityConfigurationProperties{
					AppliesToGroups:       expandConnectivityGroupItemModel(model.AppliesToGroups),
					ConnectivityTopology:  model.ConnectivityTopology,
					DeleteExistingPeering: expandDeleteExistingPeering(model.DeleteExistingPeeringEnabled),
					IsGlobal:              expandConnectivityConfIsGlobal(model.GlobalMeshEnabled),
					Hubs:                  expandHubModel(model.Hub),
					ConnectivityCapabilities: &connectivityconfigurations.ConnectivityConfigurationPropertiesConnectivityCapabilities{
						ConnectedGroupAddressOverlap:        expandConnectedGroupAddressOverlap(model.ConnectedGroupAddressOverlapEnabled),
						ConnectedGroupPrivateEndpointsScale: connectivityconfigurations.ConnectedGroupPrivateEndpointsScale(model.ConnectedGroupPrivateEndpointsScale),
						PeeringEnforcement:                  expandPeeringEnforcement(model.PeeringEnforcementEnabled),
					},
				},
			}

			if model.Description != "" {
				conf.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, conf); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandDeleteExistingPeering(input bool) *connectivityconfigurations.DeleteExistingPeering {
	output := connectivityconfigurations.DeleteExistingPeeringFalse
	if input {
		output = connectivityconfigurations.DeleteExistingPeeringTrue
	}
	return &output
}

func expandConnectedGroupAddressOverlap(input bool) connectivityconfigurations.ConnectedGroupAddressOverlap {
	if input {
		return connectivityconfigurations.ConnectedGroupAddressOverlapAllowed
	}

	return connectivityconfigurations.ConnectedGroupAddressOverlapDisallowed
}

func expandConnectivityConfIsGlobal(input bool) *connectivityconfigurations.IsGlobal {
	output := connectivityconfigurations.IsGlobalFalse
	if input {
		output = connectivityconfigurations.IsGlobalTrue
	}
	return &output
}

func expandPeeringEnforcement(input bool) connectivityconfigurations.PeeringEnforcement {
	if input {
		return connectivityconfigurations.PeeringEnforcementEnforced
	}

	return connectivityconfigurations.PeeringEnforcementUnenforced
}

func expandConnectivityGroupItemModel(inputList []ConnectivityGroupItemModel) []connectivityconfigurations.ConnectivityGroupItem {
	outputList := make([]connectivityconfigurations.ConnectivityGroupItem, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := connectivityconfigurations.ConnectivityGroupItem{
			GroupConnectivity: input.GroupConnectivity,
			IsGlobal:          expandConnectivityConfIsGlobal(input.GlobalMeshEnabled),
			NetworkGroupId:    input.NetworkGroupId,
			UseHubGateway:     expandUseHubGateWay(input.UseHubGateway),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func expandUseHubGateWay(input bool) *connectivityconfigurations.UseHubGateway {
	output := connectivityconfigurations.UseHubGatewayFalse
	if input {
		output = connectivityconfigurations.UseHubGatewayTrue
	}
	return &output
}

func expandHubModel(inputList []HubModel) *[]connectivityconfigurations.Hub {
	outputList := make([]connectivityconfigurations.Hub, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := connectivityconfigurations.Hub{
			ResourceId:   pointer.To(input.ResourceId),
			ResourceType: pointer.To(input.ResourceType),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}
