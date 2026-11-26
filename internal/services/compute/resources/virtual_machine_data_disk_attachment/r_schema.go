// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_data_disk_attachment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineDataDiskAttachmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_disk_id": {
			Description:      "The ID of an existing Managed Disk which should be attached.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     commonids.ValidateManagedDiskID,
			DiffSuppressFunc: suppress.CaseDifference,
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

		"caching": {
			Description:  "The caching requirements for this Data Disk.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForCachingTypes(), false),
		},

		"create_option": {
			Description: "The Create Option of the Data Disk, such as `Empty` or `Attach`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(virtualmachines.DiskCreateOptionTypesAttach),
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.DiskCreateOptionTypesAttach),
				string(virtualmachines.DiskCreateOptionTypesEmpty),
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
