// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sharedImageVersionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageVersionName,
		},

		"gallery_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageGalleryName,
		},

		"image_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"target_region": {
			// This needs to be a `TypeList` due to the `StateFunc` on the nested property `name`
			// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						StateFunc:        location.StateFunc,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},

					"regional_replica_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
					},

					"exclude_from_latest_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					// The Service API doesn't support to update `storage_account_type`. So it has to recreate the resource for updating `storage_account_type`.
					// However, `ForceNew` cannot be used since resource would be recreated while adding or removing `target_region`.
					// And `CustomizeDiff` also cannot be used since it doesn't support in a `Set`.
					// So currently terraform would directly return the error message from Service API while updating this property. If this property needs to be updated, please recreate this pluginsdk.
					"storage_account_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(galleryimageversions.PossibleValuesForStorageAccountType(), false),
						Default:      string(galleryimageversions.StorageAccountTypeStandardLRS),
					},
				},
			},
		},

		"blob_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			RequiredWith: []string{"storage_account_id"},
			ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"blob_uri"},
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"end_of_life_date": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"os_disk_snapshot_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
			// TODO -- add a validation function when snapshot has its own validation function
		},

		"managed_image_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				images.ValidateImageID,
				commonids.ValidateVirtualMachineID,
			),
			ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
		},

		"replication_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryimageversions.PossibleValuesForReplicationMode(), false),
			Default:      galleryimageversions.ReplicationModeFull,
		},

		"exclude_from_latest": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"deletion_of_replicated_locations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}
