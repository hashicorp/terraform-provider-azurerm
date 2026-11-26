// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func imageSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the image.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"zone_resilient": {
			Description:   "Whether zone resiliency is enabled.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			Default:       false,
			ConflictsWith: []string{"source_virtual_machine_id"},
		},

		"hyper_v_generation": {
			Description:  "The Hyper-V Generation Type of the Virtual Machine created from the image as `V1`, `V2`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(images.HyperVGenerationTypesVOne),
			ValidateFunc: validation.StringInSlice(images.PossibleValuesForHyperVGenerationTypes(), false),
		},

		"source_virtual_machine_id": {
			Description:  "The Virtual Machine ID from which to create the image.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"os_disk": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"source_virtual_machine_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"os_type": {
						Description:  "The type of operating system contained in the Virtual Machine image.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForOperatingSystemTypes(), false),
					},

					"os_state": {
						Description:  "The state of the operating system contained in the blob. Currently, the only value is Generalized.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForOperatingSystemStateTypes(), false),
					},

					"managed_disk_id": {
						Description:      "The ID of the Managed Disk resource that you want to use to create the image.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true, // azignore:AZS007 - pre-existing violation
						ValidateFunc:     commonids.ValidateManagedDiskID,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"blob_uri": {
						Description:  "The URI in Azure storage of the blob that you want to use to create the image.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"caching": {
						Description:  "The caching mode as `ReadWrite`, `ReadOnly`, or `None`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(images.CachingTypesNone),
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForCachingTypes(), false),
					},

					"size_gb": {
						Description:  "The size of the image to be created.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.NoZeroValues,
					},

					"disk_encryption_set_id": {
						Description:  "The ID of the Disk Encryption Set which should be used to encrypt this disk.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
					},

					"storage_type": {
						Description:  "The type of Storage Disk to use.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForStorageAccountTypes(), false),
					},
				},
			},
		},

		"data_disk": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"source_virtual_machine_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"lun": {
						Description: "The logical unit number of the data disk.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
					},

					"managed_disk_id": {
						Description:  "The ID of the Managed Disk resource that you want to use to create the image.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateManagedDiskID,
					},

					"blob_uri": {
						Description:  "The URI in Azure storage of the blob that you want to use to create the image.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"caching": {
						Description:  "The caching mode as `ReadWrite`, `ReadOnly`, or `None`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(images.CachingTypesNone),
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForCachingTypes(), false),
					},

					"size_gb": {
						Description:  "The size of the image to be created. The target size can't be smaller than the source size.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.NoZeroValues,
					},

					"disk_encryption_set_id": {
						Description:  "The ID of the Disk Encryption Set which should be used to encrypt this disk.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
					},

					"storage_type": {
						Description:  "The type of Storage Disk to use.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForStorageAccountTypes(), false),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
