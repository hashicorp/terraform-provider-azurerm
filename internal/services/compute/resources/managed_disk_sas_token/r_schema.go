// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_disk_sas_token

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func managedDiskSasTokenSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_disk_id": {
			Description:  "The ID of an existing Managed Disk which should be exported.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateManagedDiskID,
		},

		// unable to provide upper value of 4294967295 as it's not compatible with 32-bit (overflow errors)
		"duration_in_seconds": {
			Description:  "The duration for which the export should be allowed.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(30),
		},

		"access_level": {
			Description: "The level of access required on the disk. Supported are Read, Write.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(disks.AccessLevelRead),
				string(disks.AccessLevelWrite),
			}, false),
		},

		"sas_url": {
			Description: "The computed Shared Access Signature (SAS) of the Managed Disk.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
