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

func sortSharedImageVersions(values []compute.GalleryImageVersion) ([]compute.GalleryImageVersion, []error) {
	errors := make([]error, 0)
	sort.Slice(values, func(i, j int) bool {
		if values[i].Name == nil || values[j].Name == nil {
			return false
		}

		verA, err := version.NewVersion(*values[i].Name)
		if err != nil {
			errors = append(errors, err)
			return false
		}
		verA = version.Must(verA, err)

		verB, err := version.NewVersion(*values[j].Name)
		if err != nil {
			errors = append(errors, err)
			return false
		}
		verB = version.Must(verB, err)
		return verA.LessThan(verB)
	})

	if len(errors) > 0 {
		return values, errors
	}
	return values, nil
}
