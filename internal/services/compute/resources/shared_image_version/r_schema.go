// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sharedImageVersionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The version number for this Image Version, such as `1.0.0`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageVersionName,
		},

		"gallery_name": {
			Description:  "The name of the Shared Image Gallery in which the Shared Image exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageGalleryName,
		},

		"image_name": {
			Description:  "The name of the Shared Image within the Shared Image Gallery in which this Version should be created.",
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
						Description:      "The Azure Region in which this Image Version should exist.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						StateFunc:        location.StateFunc,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},

					"regional_replica_count": {
						Description: "The number of replicas of the Image Version to be created per region.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"disk_encryption_set_id": {
						Description:  "The ID of the Disk Encryption Set to encrypt the Image Version in the target region.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
					},

					"exclude_from_latest_enabled": {
						Description: "Whether this Shared Image Version should be excluded when querying for the `latest` version.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					// The Service API doesn't support to update `storage_account_type`. So it has to recreate the resource for updating `storage_account_type`.
					// However, `ForceNew` cannot be used since resource would be recreated while adding or removing `target_region`.
					// And `CustomizeDiff` also cannot be used since it doesn't support in a `Set`.
					// So currently terraform would directly return the error message from Service API while updating this property. If this property needs to be updated, please recreate this pluginsdk.
					"storage_account_type": {
						Description:  "The Storage Account type for the image version. You can store all of your image version replicas in Zone Redundant Storage by specifying `Standard_ZRS`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(galleryimageversions.PossibleValuesForStorageAccountType(), false),
						Default:      string(galleryimageversions.StorageAccountTypeStandardLRS),
					},
				},
			},
		},

		"blob_uri": {
			Description:  "URI of the Azure Storage Blob used to create the Image Version.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			RequiredWith: []string{"storage_account_id"},
			ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
		},

		"storage_account_id": {
			Description:  "The ID of the Storage Account where the Blob exists.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"blob_uri"},
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"end_of_life_date": {
			Description:      "The end of life date in RFC3339 format of the Image Version.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"os_disk_snapshot_id": {
			Description:  "The ID of the OS disk snapshot which should be used for this Shared Image Version.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
			// TODO -- add a validation function when snapshot has its own validation function
		},

		"managed_image_id": {
			Description: "The ID of the Managed Image or Virtual Machine ID which should be used for this Shared Image Version.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				images.ValidateImageID,
				commonids.ValidateVirtualMachineID,
			),
			ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
		},

		"replication_mode": {
			Description:  "Mode to be used for replication.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryimageversions.PossibleValuesForReplicationMode(), false),
			Default:      galleryimageversions.ReplicationModeFull,
		},

		"exclude_from_latest": {
			Description: "Should this Image Version be excluded from the `latest` filter? If set to `true` this Image Version won't be returned for the `latest` version.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"deletion_of_replicated_locations_enabled": {
			Description: "Whether this Shared Image Version can be deleted from the Azure Regions this is replicated to.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"uefi": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"signature_template_names": {
						Description: "A list of UEFI signature template names to include in the UEFI settings.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(galleryimageversions.PossibleValuesForUefiSignatureTemplateName(), false),
						},
					},
					"additional_signatures": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"db": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									ForceNew:     true,
									Elem:         helpers.UefiKeySchema(),
									AtLeastOneOf: []string{"uefi.0.additional_signatures.0.db", "uefi.0.additional_signatures.0.dbx", "uefi.0.additional_signatures.0.kek", "uefi.0.additional_signatures.0.pk"},
								},
								"dbx": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									ForceNew:     true,
									Elem:         helpers.UefiKeySchema(),
									AtLeastOneOf: []string{"uefi.0.additional_signatures.0.db", "uefi.0.additional_signatures.0.dbx", "uefi.0.additional_signatures.0.kek", "uefi.0.additional_signatures.0.pk"},
								},
								"kek": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									ForceNew:     true,
									Elem:         helpers.UefiKeySchema(),
									AtLeastOneOf: []string{"uefi.0.additional_signatures.0.db", "uefi.0.additional_signatures.0.dbx", "uefi.0.additional_signatures.0.kek", "uefi.0.additional_signatures.0.pk"},
								},
								"pk": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									ForceNew:     true,
									MaxItems:     1,
									Elem:         helpers.UefiKeySchema(),
									AtLeastOneOf: []string{"uefi.0.additional_signatures.0.db", "uefi.0.additional_signatures.0.dbx", "uefi.0.additional_signatures.0.kek", "uefi.0.additional_signatures.0.pk"},
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
