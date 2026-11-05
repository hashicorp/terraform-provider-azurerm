// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_virtual_hard_disk

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualharddisks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,62}[a-zA-Z0-9]$`),
				"name must begin and end with an alphanumeric character, be between 2 and 64 characters in length and can only contain alphanumeric characters, hyphens, periods or underscores.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"disk_size_in_gb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"block_size_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_file_format": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualharddisks.PossibleValuesForDiskFileFormat(), false),
		},

		"dynamic_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"hyperv_generation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualharddisks.PossibleValuesForHyperVGeneration(), false),
		},

		"logical_sector_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"physical_sector_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"storage_path_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: storagecontainers.ValidateStorageContainerID,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
