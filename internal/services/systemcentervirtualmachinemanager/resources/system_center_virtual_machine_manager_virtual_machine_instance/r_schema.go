// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scoped_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
		},

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}, "Custom Location"),

		"infrastructure": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"checkpoint_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Disabled",
							"Production",
							"ProductionOnly",
							"Standard",
						}, false),
					},

					"system_center_virtual_machine_manager_cloud_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: clouds.ValidateCloudID,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_template_id"},
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},

					"system_center_virtual_machine_manager_inventory_item_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: inventoryitems.ValidateInventoryItemID,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},

					"system_center_virtual_machine_manager_template_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: virtualmachinetemplates.ValidateVirtualMachineTemplateID,
						RequiredWith: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id"},
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},

					"system_center_virtual_machine_manager_virtual_machine_server_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: vmmservers.ValidateVMmServerID,
						AtLeastOneOf: []string{"infrastructure.0.system_center_virtual_machine_manager_cloud_id", "infrastructure.0.system_center_virtual_machine_manager_inventory_item_id", "infrastructure.0.system_center_virtual_machine_manager_template_id", "infrastructure.0.system_center_virtual_machine_manager_virtual_machine_server_id"},
					},
				},
			},
		},

		"hardware": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 64),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"dynamic_memory_max_in_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(32, 1048576),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"dynamic_memory_min_in_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(32, 1048576),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"limit_cpu_for_migration_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},

					"memory_in_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(32, 1048576),
						AtLeastOneOf: []string{"hardware.0.cpu_count", "hardware.0.dynamic_memory_max_in_mb", "hardware.0.dynamic_memory_min_in_mb", "hardware.0.limit_cpu_for_migration_enabled", "hardware.0.memory_in_mb"},
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: networkValidate.NetworkInterfaceName,
					},

					"virtual_network_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: virtualnetworks.ValidateVirtualNetworkID,
					},

					"ipv4_address_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForAllocationMethod(), false),
					},

					"ipv6_address_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForAllocationMethod(), false),
					},

					"mac_address_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualmachineinstances.PossibleValuesForAllocationMethod(), false),
					},
				},
			},
		},

		"operating_system": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceComputerName,
						AtLeastOneOf: []string{"operating_system.0.computer_name", "operating_system.0.admin_password"},
					},

					"admin_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"operating_system.0.computer_name", "operating_system.0.admin_password"},
					},
				},
			},
		},

		"storage_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bus": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 3),
					},

					"bus_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"IDE",
							"SCSI",
						}, false),
					},

					"disk_size_gb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"lun": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 63),
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDiskName,
					},

					"storage_qos_policy_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"template_disk_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsUUID,
					},

					"vhd_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Dynamic",
							"Fixed",
						}, false),
					},
				},
			},
		},

		"system_center_virtual_machine_manager_availability_set_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: availabilitysets.ValidateAvailabilitySetID,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
