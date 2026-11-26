// Copyright IBM Corp. 2014, 2026
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
			Description: "The name of the Managed Disk.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"storage_account_type": {
			Description:      "The type of storage to use for the Managed Disk.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringInSlice(disks.PossibleValuesForDiskStorageAccountTypes(), false),
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"create_option": {
			Description: "The method to use when creating the Managed Disk. * `ImportSecure` - Securely import a VHD file in to the Managed Disk (VHD specified with `source_uri`). * `Empty` - Create an empty Managed Disk. * `Copy` - Copy an existing Managed Disk or snapshot (specified with `source_resource_id`). * `FromImage` - Copy a Platform Image (specified with `image_reference_id`) * `Restore` - Set by Azure Backup or Site Recovery on a restored disk (specified with `source_resource_id`). * `Upload` - Upload a VHD disk with the help of SAS URL (to be used with `upload_size_bytes`).",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
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
			Description: "Logical Sector Size.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			ValidateFunc: validation.IntInSlice([]int{
				512,
				4096,
			}),
		},

		"optimized_frequent_attach_enabled": {
			Description: "Whether this Managed Disk should be optimized for frequent disk attachments (where a disk is attached/detached more than 5 times in a day).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"performance_plus_enabled": {
			Description: "Whether Performance Plus is enabled for this Managed Disk.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"source_uri": {
			Description: "URI to a valid VHD file to be used when `create_option` is `Import` or `ImportSecure`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"source_resource_id": {
			Description: "The ID of an existing Managed Disk or Snapshot to copy when `create_option` is `Copy` or the recovery point to restore when `create_option` is `Restore`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"storage_account_id": {
			Description:  "The ID of the Storage Account where the `source_uri` is located.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true, // Not supported by disk update
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"image_reference_id": {
			Description:   "ID of an existing platform/marketplace disk image to copy when `create_option` is `FromImage`. This field cannot be specified if gallery_image_reference_id is specified.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"gallery_image_reference_id"},
		},

		"gallery_image_reference_id": {
			Description:   "ID of a Gallery Image Version to copy when `create_option` is `FromImage`. This field cannot be specified if image_reference_id is specified.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.AsGeneratedID(galleryimageversions.ParseImageVersionIDInsensitively),
			ConflictsWith: []string{"image_reference_id"},
		},

		"os_type": {
			Description:  "A value when the source of an `Import`, `ImportSecure` or `Copy` operation targets a source that contains an operating system.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForOperatingSystemTypes(), false),
		},

		"disk_size_gb": {
			Description: "(Optional, Required for a new Managed Disk) Specifies the size of the Managed Disk to create in gigabytes. If `create_option` is `Copy` or `FromImage`, then the value must be equal to or greater than the source's size. The size can only be increased.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because Azure computes disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntBetween(0, 65536),
		},

		"upload_size_bytes": {
			Description:  "The size of the Managed Disk to create in bytes. Source disk size could be calculated with `ls -l` or `wc -c`. More information can be found at [Copy a Managed Disk](https://learn.microsoft.com/azure/virtual-machines/linux/disks-upload-vhd-to-managed-disk-cli#copy-a-managed-disk).",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_iops_read_write": {
			Description: "The number of IOPS allowed for this disk; only settable for UltraSSD disks and PremiumV2 disks. One operation can transfer between 4k and 256k bytes.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because Azure assigns IOPS based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_mbps_read_write": {
			Description: "The bandwidth allowed for this disk; only settable for UltraSSD disks and PremiumV2 disks. MBps means millions of bytes per second.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because Azure assigns throughput based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_iops_read_only": {
			Description: "The number of IOPS allowed across all VMs mounting the shared disk as read-only; only settable for UltraSSD disks and PremiumV2 disks with shared disk enabled. One operation can transfer between 4k and 256k bytes.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because Azure assigns read-only IOPS based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_mbps_read_only": {
			Description: "The bandwidth allowed across all VMs mounting the shared disk as read-only; only settable for UltraSSD disks and PremiumV2 disks with shared disk enabled. MBps means millions of bytes per second.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because Azure assigns read-only throughput based on disk size when not specified
			Computed:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_encryption_set_id": {
			Description:  "The ID of a Disk Encryption Set which should be used to encrypt this Managed Disk.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
			// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
			//    https://github.com/Azure/azure-rest-api-specs/issues/8132
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith:    []string{"secure_vm_disk_encryption_set_id"},
		},

		"encryption_settings": helpers.EncryptionSettingsSchema(),

		"network_access_policy": {
			Description:  "Policy for accessing the disk via network.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      disks.NetworkAccessPolicyAllowAll,
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForNetworkAccessPolicy(), false),
		},
		"disk_access_id": {
			Description:  "The ID of the Disk Access resource for using Private Endpoints on disks.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: diskaccesses.ValidateDiskAccessID,
			// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
			//    https://github.com/Azure/azure-rest-api-specs/issues/14192
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"public_network_access_enabled": {
			Description: "Whether it is allowed to access the disk via public network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tier": {
			Description: "The disk performance tier to use. This feature is currently supported only for premium SSDs.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"max_shares": {
			Description:  "The maximum number of VMs that can attach to the disk at the same time. Value greater than one indicates a disk that can be mounted on multiple VMs at the same time.",
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(2, 10),
		},

		"trusted_launch_enabled": {
			Description: "Whether Trusted Launch is enabled for the Managed Disk.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"secure_vm_disk_encryption_set_id": {
			Description:   "The ID of the Disk Encryption Set which should be used to Encrypt this OS Disk when the Virtual Machine is a Confidential VM.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
			ConflictsWith: []string{"disk_encryption_set_id"},
		},

		"security_type": {
			Description: "Security Type of the Managed Disk when it is used for a Confidential VM.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(disks.DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey),
				string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey),
				string(disks.DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey),
			}, false),
		},

		"hyper_v_generation": {
			Description:  "The HyperV Generation of the Disk when the source of an `Import` or `Copy` operation targets a source that contains an operating system. For `ImportSecure` it must be set to `V2`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true, // Not supported by disk update
			ValidateFunc: validation.StringInSlice(disks.PossibleValuesForHyperVGeneration(), false),
		},

		"on_demand_bursting_enabled": {
			Description: "Whether On-Demand Bursting is enabled for the Managed Disk.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
