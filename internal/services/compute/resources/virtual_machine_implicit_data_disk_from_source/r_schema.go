// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_implicit_data_disk_from_source

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Data Disk.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_machine_id": {
			Description:  "The ID of the Virtual Machine to which the Data Disk should be attached.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"lun": {
			Description:  "The Logical Unit Number of the Data Disk, which needs to be unique within the Virtual Machine.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"create_option": {
			Description: "The Create Option of the Data Disk.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.DiskCreateOptionTypesCopy),
			}, false),
		},

		"disk_size_gb": {
			Description:  "The size of the Data Disk in gigabytes.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 1023),
		},

		"source_resource_id": {
			Description:  "The ID of the source resource which this Data Disk was created from.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(snapshots.ValidateSnapshotID, commonids.ValidateManagedDiskID),
		},

		"caching": {
			Description: "The caching requirements for this Data Disk.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.CachingTypesReadOnly),
				string(virtualmachines.CachingTypesReadWrite),
			}, false),
		},

		"write_accelerator_enabled": {
			Description: "Whether Write Accelerator is enabled on the disk. This can only be enabled on `Premium_LRS` Managed Disks with no caching and [M-Series VMs](https://learn.microsoft.com/azure/virtual-machines/workloads/sap/how-to-enable-write-accelerator).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldSize, newSize := metadata.ResourceDiff.GetChange("disk_size_gb"); newSize.(int) < oldSize.(int) {
				if err := metadata.ResourceDiff.ForceNew("disk_size_gb"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
