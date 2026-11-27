// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	systemcentervirtualmachinemanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := systemcentervirtualmachinemanagerClient.From(metadata.Client).VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{
				ScopedResourceId: id.Scope,
			}

			if model := resp.Model; model != nil {
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)

				if props := model.Properties; props != nil {
					state.Hardware = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(props.HardwareProfile)
					state.Infrastructure = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(props.InfrastructureProfile)
					state.OperatingSystem = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(props.OsProfile, metadata.ResourceData.Get("operating_system.0.admin_password").(string))
					state.SystemCenterVirtualMachineManagerAvailabilitySetIds = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(props.AvailabilitySets)

					if v := props.NetworkProfile; v != nil {
						state.NetworkInterfaces = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(v.NetworkInterfaces)
					}

					if v := props.StorageProfile; v != nil {
						state.StorageDisks = flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisks(v.Disks)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

type Hardware struct {
	CpuCount                    int64 `tfschema:"cpu_count"`
	DynamicMemoryMaxInMb        int64 `tfschema:"dynamic_memory_max_in_mb"`
	DynamicMemoryMinInMb        int64 `tfschema:"dynamic_memory_min_in_mb"`
	LimitCpuForMigrationEnabled bool  `tfschema:"limit_cpu_for_migration_enabled"`
	MemoryInMb                  int64 `tfschema:"memory_in_mb"`
}

type Infrastructure struct {
	CheckpointType                                          string `tfschema:"checkpoint_type"`
	SystemCenterVirtualMachineManagerCloudId                string `tfschema:"system_center_virtual_machine_manager_cloud_id"`
	SystemCenterVirtualMachineManagerInventoryItemId        string `tfschema:"system_center_virtual_machine_manager_inventory_item_id"`
	SystemCenterVirtualMachineManagerTemplateId             string `tfschema:"system_center_virtual_machine_manager_template_id"`
	SystemCenterVirtualMachineManagerVirtualMachineServerId string `tfschema:"system_center_virtual_machine_manager_virtual_machine_server_id"`
}

type NetworkInterface struct {
	Name             string `tfschema:"name"`
	Ipv4AddressType  string `tfschema:"ipv4_address_type"`
	Ipv6AddressType  string `tfschema:"ipv6_address_type"`
	MacAddressType   string `tfschema:"mac_address_type"`
	VirtualNetworkId string `tfschema:"virtual_network_id"`
}

type OperatingSystem struct {
	ComputerName  string `tfschema:"computer_name"`
	AdminPassword string `tfschema:"admin_password"`
}

type StorageDisk struct {
	Bus                  int64  `tfschema:"bus"`
	BusType              string `tfschema:"bus_type"`
	DiskSizeGB           int64  `tfschema:"disk_size_gb"`
	Lun                  int64  `tfschema:"lun"`
	Name                 string `tfschema:"name"`
	StorageQoSPolicyName string `tfschema:"storage_qos_policy_name"`
	TemplateDiskId       string `tfschema:"template_disk_id"`
	VhdType              string `tfschema:"vhd_type"`
}
