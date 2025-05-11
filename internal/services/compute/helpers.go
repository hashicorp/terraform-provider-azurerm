// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"sort"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/go-version"
)

func expandIDsToSubResources(input []interface{}) *[]virtualmachinescalesets.SubResource {
	ids := make([]virtualmachinescalesets.SubResource, 0)

	for _, v := range input {
		ids = append(ids, virtualmachinescalesets.SubResource{
			Id: pointer.To(v.(string)),
		})
	}

	return &ids
}

func flattenSubResourcesToIDs(input *[]virtualmachinescalesets.SubResource) []interface{} {
	ids := make([]interface{}, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}

func flattenSubResourcesToStringIDs(input *[]virtualmachinescalesets.SubResource) []string {
	ids := make([]string, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}

func sortSharedImageVersions(values []galleryimageversions.GalleryImageVersion) ([]galleryimageversions.GalleryImageVersion, []error) {
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
