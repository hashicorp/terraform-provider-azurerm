package network

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandEdgeZone(input string) *network.ExtendedLocation {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &network.ExtendedLocation{
		Name: utils.String(normalized),
		Type: network.ExtendedLocationTypesEdgeZone,
	}
}

func flattenEdgeZone(input *network.ExtendedLocation) string {
	if input == nil || input.Type != network.ExtendedLocationTypesEdgeZone || input.Name == nil {
		return ""
	}
	return edgezones.NormalizeNilable(input.Name)
}
