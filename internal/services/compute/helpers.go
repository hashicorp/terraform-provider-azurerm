package compute

import (
	"sort"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

func sortSharedImageVersions(values []compute.GalleryImageVersion) []compute.GalleryImageVersion {
	sort.Slice(values, func(i, j int) bool {
		verA, _ := version.NewVersion(*values[i].Name)
		verB, _ := version.NewVersion(*values[j].Name)
		return verA.LessThan(verB)
	})
	return values
}
