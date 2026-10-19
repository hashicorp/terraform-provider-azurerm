// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"zone_resilient": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Default:       false,
			ForceNew:      true,
			ConflictsWith: []string{"source_virtual_machine_id"},
		},

		"hyper_v_generation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(images.HyperVGenerationTypesVOne),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(images.PossibleValuesForHyperVGenerationTypes(), false),
		},

		"source_virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"os_disk": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ForceNew:      true,
			ConflictsWith: []string{"source_virtual_machine_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"os_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForOperatingSystemTypes(), false),
					},

					"os_state": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForOperatingSystemStateTypes(), false),
					},

					"managed_disk_id": {
						Type:             pluginsdk.TypeString,
						Computed:         true, // azignore:AZS007 - pre-existing violation
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     commonids.ValidateManagedDiskID,
					},

					"blob_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"caching": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(images.CachingTypesNone),
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForCachingTypes(), false),
					},

					"size_gb": {
						Type:         pluginsdk.TypeInt,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.NoZeroValues,
					},

					"disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
					},

					"storage_type": {
						Type:         pluginsdk.TypeString,
						Description:  "The type of storage disk",
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
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"managed_disk_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateManagedDiskID,
					},

					"blob_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"caching": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(images.CachingTypesNone),
						ValidateFunc: validation.StringInSlice(images.PossibleValuesForCachingTypes(), false),
					},

					"size_gb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.NoZeroValues,
					},

					"disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
					},

					"storage_type": {
						Type:         pluginsdk.TypeString,
						Description:  "The type of storage disk",
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
