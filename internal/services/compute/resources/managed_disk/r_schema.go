// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func managedDiskSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"storage_account_type": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringInSlice(disks.PossibleValuesForDiskStorageAccountTypes(), false),
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"create_option": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(disks.DiskCreateOptionCopy),
				string(disks.DiskCreateOptionEmpty),
				string(disks.DiskCreateOptionFromImage),
				string(disks.DiskCreateOptionImport),
				string(disks.DiskCreateOptionImportSecure),
				string(disks.DiskCreateOptionRestore),
				string(disks.DiskCreateOptionUpload),
			}, false),
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"logical_sector_size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.IntInSlice([]int{
				512,
				4096,
			}),
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"optimized_frequent_attach_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"performance_plus_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"source_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"source_resource_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true, // Not supported by disk update
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"image_reference_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"gallery_image_reference_id"},
		},

		"gallery_image_reference_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.AsGeneratedID(galleryimageversions.ParseImageVersionIDInsensitively),
			ConflictsWith: []string{"image_reference_id"},
		},

		"os_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForOperatingSystemTypes(), false),
		},

		"disk_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because Azure computes disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntBetween(0, 65536),
		},

		"upload_size_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_iops_read_write": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because Azure assigns IOPS based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_mbps_read_write": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because Azure assigns throughput based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_iops_read_only": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because Azure assigns read-only IOPS based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_mbps_read_only": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because Azure assigns read-only throughput based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
			//    https://github.com/Azure/azure-rest-api-specs/issues/8132
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
			ConflictsWith:    []string{"secure_vm_disk_encryption_set_id"},
		},

		"encryption_settings": helpers.EncryptionSettingsSchema(),

		"network_access_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      disks.NetworkAccessPolicyAllowAll,
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForNetworkAccessPolicy(), false),
		},
		"disk_access_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
			//    https://github.com/Azure/azure-rest-api-specs/issues/14192
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     diskaccesses.ValidateDiskAccessID,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"max_shares": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(2, 10),
		},

		"trusted_launch_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"secure_vm_disk_encryption_set_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
			ConflictsWith: []string{"disk_encryption_set_id"},
		},

		"security_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(disks.DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey),
				string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey),
				string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey),
			}, false),
		},

		"hyper_v_generation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true, // Not supported by disk update
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForHyperVGeneration(), false),
		},

		"on_demand_bursting_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
