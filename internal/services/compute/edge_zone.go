package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandEdgeZone(input string) *compute.ExtendedLocation {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &compute.ExtendedLocation{
		Name: utils.String(normalized),
		Type: compute.ExtendedLocationTypesEdgeZone,
	}
}

func flattenEdgeZone(input *compute.ExtendedLocation) string {
	if input == nil || input.Type != compute.ExtendedLocationTypesEdgeZone || input.Name == nil {
		return ""
	}
	return edgezones.NormalizeNilable(input.Name)
}
