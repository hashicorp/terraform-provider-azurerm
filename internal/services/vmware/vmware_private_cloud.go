// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
)

func flattenPrivateCloudManagementCluster(input *privateclouds.CommonClusterProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"size":  input.ClusterSize,
			"id":    input.ClusterId,
			"hosts": helpers.FlattenStringSlice(input.Hosts),
		},
	}
}

func flattenPrivateCloudCircuit(input *privateclouds.Circuit) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	expressRouteId := pointer.From(input.ExpressRouteID)
	expressRoutePrivatePeeringId := pointer.From(input.ExpressRoutePrivatePeeringID)
	primarySubnet := pointer.From(input.PrimarySubnet)
	secondarySubnet := pointer.From(input.SecondarySubnet)
	return []interface{}{
		map[string]interface{}{
			"express_route_id":                 expressRouteId,
			"express_route_private_peering_id": expressRoutePrivatePeeringId,
			"primary_subnet_cidr":              primarySubnet,
			"secondary_subnet_cidr":            secondarySubnet,
		},
	}
}
