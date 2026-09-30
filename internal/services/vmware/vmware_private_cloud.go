// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenPrivateCloudManagementCluster(input *privateclouds.CommonClusterProperties) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"size":  input.ClusterSize,
			"id":    input.ClusterId,
			"hosts": pluginsdk.FlattenSlice(input.Hosts),
		},
	}
}

func flattenPrivateCloudCircuit(input *privateclouds.Circuit) []any {
	if input == nil {
		return make([]any, 0)
	}

	expressRouteId := pointer.From(input.ExpressRouteID)
	expressRoutePrivatePeeringId := pointer.From(input.ExpressRoutePrivatePeeringID)
	primarySubnet := pointer.From(input.PrimarySubnet)
	secondarySubnet := pointer.From(input.SecondarySubnet)
	return []any{
		map[string]any{
			"express_route_id":                 expressRouteId,
			"express_route_private_peering_id": expressRoutePrivatePeeringId,
			"primary_subnet_cidr":              primarySubnet,
			"secondary_subnet_cidr":            secondarySubnet,
		},
	}
}
