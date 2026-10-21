// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_connectivity_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectivityconfigurations"
)

func flattenDeleteExistingPeering(input *connectivityconfigurations.DeleteExistingPeering) bool {
	if input == nil {
		return false
	}
	return *input == connectivityconfigurations.DeleteExistingPeeringTrue
}

func flattenConnectedGroupAddressOverlap(input connectivityconfigurations.ConnectedGroupAddressOverlap) bool {
	return input != connectivityconfigurations.ConnectedGroupAddressOverlapDisallowed
}

func flattenConnectivityConfIsGlobal(input *connectivityconfigurations.IsGlobal) bool {
	if input == nil {
		return false
	}
	return *input == connectivityconfigurations.IsGlobalTrue
}

func flattenPeeringEnforcement(input connectivityconfigurations.PeeringEnforcement) bool {
	return input == connectivityconfigurations.PeeringEnforcementEnforced
}

func flattenConnectivityGroupItemModel(inputList []connectivityconfigurations.ConnectivityGroupItem) []ConnectivityGroupItemModel {
	outputList := make([]ConnectivityGroupItemModel, 0, len(inputList))
	for _, input := range inputList {
		output := ConnectivityGroupItemModel{
			GroupConnectivity: input.GroupConnectivity,
			UseHubGateway:     flattenUseHubGateWay(input.UseHubGateway),
			GlobalMeshEnabled: flattenConnectivityConfIsGlobal(input.IsGlobal),
			NetworkGroupId:    input.NetworkGroupId,
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenUseHubGateWay(input *connectivityconfigurations.UseHubGateway) bool {
	if input == nil {
		return false
	}
	return *input == connectivityconfigurations.UseHubGatewayTrue
}

func flattenHubModel(inputList *[]connectivityconfigurations.Hub) []HubModel {
	if inputList == nil {
		return []HubModel{}
	}

	outputList := make([]HubModel, 0, len(*inputList))
	for _, input := range *inputList {
		output := HubModel{}

		if input.ResourceId != nil {
			output.ResourceId = *input.ResourceId
		}

		if input.ResourceType != nil {
			output.ResourceType = *input.ResourceType
		}

		outputList = append(outputList, output)
	}

	return outputList
}
