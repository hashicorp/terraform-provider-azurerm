package vmware

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func flattenPrivateCloudManagementCluster(input privateclouds.ManagementCluster) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"size":  input.ClusterSize,
			"id":    input.ClusterId,
			"hosts": utils.FlattenStringSlice(input.Hosts),
		},
	}
}

func flattenPrivateCloudCircuit(input *privateclouds.Circuit) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var expressRouteId string
	if input.ExpressRouteID != nil {
		expressRouteId = *input.ExpressRouteID
	}
	var expressRoutePrivatePeeringId string
	if input.ExpressRoutePrivatePeeringID != nil {
		expressRoutePrivatePeeringId = *input.ExpressRoutePrivatePeeringID
	}
	var primarySubnet string
	if input.PrimarySubnet != nil {
		primarySubnet = *input.PrimarySubnet
	}
	var secondarySubnet string
	if input.SecondarySubnet != nil {
		secondarySubnet = *input.SecondarySubnet
	}
	return []interface{}{
		map[string]interface{}{
			"express_route_id":                 expressRouteId,
			"express_route_private_peering_id": expressRoutePrivatePeeringId,
			"primary_subnet_cidr":              primarySubnet,
			"secondary_subnet_cidr":            secondarySubnet,
		},
	}
}
