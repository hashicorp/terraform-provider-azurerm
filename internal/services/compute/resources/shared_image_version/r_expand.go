// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandSharedImageVersionTargetRegions(d *pluginsdk.ResourceData) (*[]galleryimageversions.TargetRegion, error) {
	vs := d.Get("target_region").([]any)
	results := make([]galleryimageversions.TargetRegion, 0)

	for _, v := range vs {
		input := v.(map[string]any)

		name := input["name"].(string)
		regionalReplicaCount := input["regional_replica_count"].(int)
		storageAccountType := input["storage_account_type"].(string)
		diskEncryptionSetId := input["disk_encryption_set_id"].(string)
		excludeFromLatest := input["exclude_from_latest_enabled"].(bool)

		output := galleryimageversions.TargetRegion{
			Name:                 name,
			ExcludeFromLatest:    pointer.To(excludeFromLatest),
			RegionalReplicaCount: pointer.To(int64(regionalReplicaCount)),
			StorageAccountType:   pointer.ToEnum[galleryimageversions.StorageAccountType](storageAccountType),
		}

		if diskEncryptionSetId != "" {
			if d.Get("replication_mode").(string) == string(galleryimageversions.ReplicationModeShallow) {
				return nil, fmt.Errorf("`disk_encryption_set_id` cannot be used when `replication_mode` is `Shallow`")
			}

			output.Encryption = &galleryimageversions.EncryptionImages{
				OsDiskImage: &galleryimageversions.OSDiskImageEncryption{
					DiskEncryptionSetId: pointer.To(diskEncryptionSetId),
				},
			}
		}

		results = append(results, output)
	}

	return &results, nil
}

func expandUefiSettings(input []any) *galleryimageversions.GalleryImageVersionUefiSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	result := &galleryimageversions.GalleryImageVersionUefiSettings{}

	if templateNamesSet, ok := v["signature_template_names"].(*pluginsdk.Set); ok {
		result.SignatureTemplateNames = expandSignatureTemplateNames(templateNamesSet.List())
	}

	if additionalSignatures, ok := v["additional_signatures"].([]any); ok {
		result.AdditionalSignatures = expandAdditionalSignatures(additionalSignatures)
	}

	return result
}

func expandSignatureTemplateNames(input []any) *[]galleryimageversions.UefiSignatureTemplateName {
	if len(input) == 0 {
		return nil
	}

	result := make([]galleryimageversions.UefiSignatureTemplateName, 0)
	for _, v := range input {
		result = append(result, galleryimageversions.UefiSignatureTemplateName(v.(string)))
	}
	return &result
}

func expandAdditionalSignatures(input []any) *galleryimageversions.UefiKeySignatures {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	result := &galleryimageversions.UefiKeySignatures{}

	if db, ok := v["db"].([]any); ok {
		result.Db = expandUefiKeyList(db)
	}

	if dbx, ok := v["dbx"].([]any); ok {
		result.Dbx = expandUefiKeyList(dbx)
	}

	if kek, ok := v["kek"].([]any); ok {
		result.Kek = expandUefiKeyList(kek)
	}

	if pk, ok := v["pk"].([]any); ok {
		result.Pk = expandUefiKey(pk)
	}

	return result
}

func expandUefiKeyList(input []any) *[]galleryimageversions.UefiKey {
	if len(input) == 0 {
		return nil
	}

	result := make([]galleryimageversions.UefiKey, 0)
	for _, v := range input {
		if item := expandUefiKey([]any{v}); item != nil {
			result = append(result, *item)
		}
	}
	return &result
}

func expandUefiKey(input []any) *galleryimageversions.UefiKey {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	data := input[0].(map[string]any)

	certData := make([]string, 0)
	if certList, ok := data["values_base64"].([]any); ok {
		for _, item := range certList {
			if str, ok := item.(string); ok {
				certData = append(certData, str)
			}
		}
	}

	return &galleryimageversions.UefiKey{
		Type:  pointer.ToEnum[galleryimageversions.UefiKeyType](data["type"].(string)),
		Value: &certData,
	}
}
