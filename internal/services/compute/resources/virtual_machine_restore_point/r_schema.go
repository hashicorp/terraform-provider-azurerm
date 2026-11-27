// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_restore_point

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Virtual Machine Restore Point.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"virtual_machine_restore_point_collection_id": {
			Description:  "The ID of the Virtual Machine Restore Point Collection the Virtual Machine Restore Point will be associated with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: restorepointcollections.ValidateRestorePointCollectionID,
		},

		"crash_consistency_mode_enabled": {
			Description: "Whether the Consistency Mode of the Virtual Machine Restore Point is set to `CrashConsistent`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"excluded_disks": {
			Description: "A list of disks that will be excluded from the Virtual Machine Restore Point.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateManagedDiskID,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
