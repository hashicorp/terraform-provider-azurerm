// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
)

func FlattenSharedImageVersionDataSourceTargetRegions(input *[]galleryimageversions.TargetRegion) []any {
	results := make([]any, 0)

	if input != nil {
		for _, v := range *input {
			output := make(map[string]any)

			output["name"] = location.Normalize(v.Name)

			if v.RegionalReplicaCount != nil {
				output["regional_replica_count"] = int(*v.RegionalReplicaCount)
			}

			if v.StorageAccountType != nil {
				output["storage_account_type"] = string(*v.StorageAccountType)
			}

			results = append(results, output)
		}
	}

	return results
}
