// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
)

func flattenUefiSettings(input *galleryimageversions.GalleryImageVersionUefiSettings) []any {
	results := make([]any, 0)

	if input == nil {
		return results
	}

	results = append(results, map[string]any{
		"signature_template_names": pointer.From(input.SignatureTemplateNames),
		"additional_signatures":    flattenAdditionalSignatures(input.AdditionalSignatures),
	})

	return results
}

func flattenAdditionalSignatures(input *galleryimageversions.UefiKeySignatures) []any {
	results := make([]any, 0)

	if input == nil {
		return results
	}

	result := make(map[string]any)
	result["db"] = flattenUefiKeyList(input.Db)
	result["dbx"] = flattenUefiKeyList(input.Dbx)
	result["kek"] = flattenUefiKeyList(input.Kek)
	result["pk"] = flattenUefiKey(input.Pk)

	return append(results, result)
}

func flattenUefiKeyList(input *[]galleryimageversions.UefiKey) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		if item := flattenUefiKey(&v); len(item) > 0 {
			results = append(results, item[0])
		}
	}

	return results
}

func flattenUefiKey(input *galleryimageversions.UefiKey) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	result := make(map[string]any)
	if input.Value != nil && len(*input.Value) > 0 {
		result["values_base64"] = *input.Value
	}
	result["type"] = pointer.From(input.Type)

	return append(results, result)
}

func flattenSharedImageVersionTargetRegions(input *[]galleryimageversions.TargetRegion) []any {
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

			diskEncryptionSetId := ""
			if v.Encryption != nil && v.Encryption.OsDiskImage != nil && v.Encryption.OsDiskImage.DiskEncryptionSetId != nil {
				diskEncryptionSetId = *v.Encryption.OsDiskImage.DiskEncryptionSetId
			}
			output["disk_encryption_set_id"] = diskEncryptionSetId

			output["exclude_from_latest_enabled"] = pointer.From(v.ExcludeFromLatest)

			results = append(results, output)
		}
	}

	return results
}
