package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func expandIDsToSubResources(input []interface{}) *[]compute.SubResource {
	ids := make([]compute.SubResource, 0)

	for _, v := range input {
		ids = append(ids, compute.SubResource{
			ID: utils.String(v.(string)),
		})
	}

	return &ids
}

func flattenSubResourcesToIDs(input *[]compute.SubResource) []interface{} {
	ids := make([]interface{}, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.ID == nil {
			continue
		}

		ids = append(ids, *v.ID)
	}

	return ids
}
