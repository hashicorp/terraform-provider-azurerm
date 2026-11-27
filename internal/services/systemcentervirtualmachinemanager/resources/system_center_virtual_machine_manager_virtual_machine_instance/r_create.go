// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	systemcentervirtualmachinemanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := systemcentervirtualmachinemanagerClient.From(metadata.Client).VirtualMachineInstances

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewSystemCenterVirtualMachineManagerVirtualMachineInstanceID(model.ScopedResourceId)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, commonids.NewScopeID(id.Scope))
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := virtualmachineinstances.VirtualMachineInstance{
				ExtendedLocation: virtualmachineinstances.ExtendedLocation{
					Type: pointer.To("customLocation"),
					Name: pointer.To(model.CustomLocationId),
				},
				Properties: &virtualmachineinstances.VirtualMachineInstanceProperties{
					InfrastructureProfile: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(model.Infrastructure),
					HardwareProfile:       expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(model.Hardware, metadata.ResourceData),
					OsProfile:             expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(model.OperatingSystem),
				},
			}

			if v := model.NetworkInterfaces; v != nil {
				parameters.Properties.NetworkProfile = &virtualmachineinstances.NetworkProfile{
					NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(v),
				}
			}

			if v := model.StorageDisks; v != nil {
				parameters.Properties.StorageProfile = &virtualmachineinstances.StorageProfile{
					Disks: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForCreate(v, metadata.ResourceData),
				}
			}

			if v := model.SystemCenterVirtualMachineManagerAvailabilitySetIds; v != nil {
				availabilitySets, err := expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(v)
				if err != nil {
					return err
				}
				parameters.Properties.AvailabilitySets = availabilitySets
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, commonids.NewScopeID(id.Scope), parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForCreate(input []Infrastructure) *virtualmachineinstances.InfrastructureProfile {
	if len(input) == 0 {
		return nil
	}

	infrastructureProfile := input[0]

	result := virtualmachineinstances.InfrastructureProfile{}

	if v := infrastructureProfile.CheckpointType; v != "" {
		result.CheckpointType = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerCloudId; v != "" {
		result.CloudId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerInventoryItemId; v != "" {
		result.InventoryItemId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerTemplateId; v != "" {
		result.TemplateId = pointer.To(v)
	}

	if v := infrastructureProfile.SystemCenterVirtualMachineManagerVirtualMachineServerId; v != "" {
		result.VMmServerId = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForCreate(input []Hardware, d *pluginsdk.ResourceData) *virtualmachineinstances.HardwareProfile {
	if len(input) == 0 {
		return nil
	}

	hardwareProfile := input[0]

	result := virtualmachineinstances.HardwareProfile{}

	// As TF always sets bool value to false when it isn't set, so it has to use d.GetRawConfig() to determine whether it is set in the tf config
	if v := d.GetRawConfig().AsValueMap()["hardware"].AsValueSlice()[0].AsValueMap()["limit_cpu_for_migration_enabled"]; !v.IsNull() {
		result.LimitCPUForMigration = pointer.ToEnum[virtualmachineinstances.LimitCPUForMigration](strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))
	}

	if v := hardwareProfile.CpuCount; v != 0 {
		result.CpuCount = pointer.To(v)
	}

	dynamicMemoryEnabled := hardwareProfile.DynamicMemoryMaxInMb != 0 || hardwareProfile.DynamicMemoryMinInMb != 0

	result.DynamicMemoryEnabled = pointer.ToEnum[virtualmachineinstances.DynamicMemoryEnabled](strconv.FormatBool(dynamicMemoryEnabled))

	if v := hardwareProfile.DynamicMemoryMaxInMb; v != 0 {
		result.DynamicMemoryMaxMB = pointer.To(v)
	}

	if v := hardwareProfile.DynamicMemoryMinInMb; v != 0 {
		result.DynamicMemoryMinMB = pointer.To(v)
	}

	if v := hardwareProfile.MemoryInMb; v != 0 {
		result.MemoryMB = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForCreate(input []NetworkInterface) *[]virtualmachineinstances.NetworkInterface {
	result := make([]virtualmachineinstances.NetworkInterface, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterface := virtualmachineinstances.NetworkInterface{
			Name: pointer.To(v.Name),
		}

		if ipv4AddressType := v.Ipv4AddressType; ipv4AddressType != "" {
			networkInterface.IPv4AddressType = pointer.ToEnum[virtualmachineinstances.AllocationMethod](ipv4AddressType)
		}

		if ipv6AddressType := v.Ipv6AddressType; ipv6AddressType != "" {
			networkInterface.IPv6AddressType = pointer.ToEnum[virtualmachineinstances.AllocationMethod](ipv6AddressType)
		}

		if macAddressType := v.MacAddressType; macAddressType != "" {
			networkInterface.MacAddressType = pointer.ToEnum[virtualmachineinstances.AllocationMethod](macAddressType)
		}

		if vnetId := v.VirtualNetworkId; vnetId != "" {
			networkInterface.VirtualNetworkId = pointer.To(vnetId)
		}

		result = append(result, networkInterface)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input []OperatingSystem) *virtualmachineinstances.OsProfileForVMInstance {
	if len(input) == 0 {
		return nil
	}

	osProfile := input[0]

	result := virtualmachineinstances.OsProfileForVMInstance{}

	if v := osProfile.ComputerName; v != "" {
		result.ComputerName = pointer.To(v)
	}

	if v := osProfile.AdminPassword; v != "" {
		result.AdminPassword = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForCreate(input []StorageDisk, d *pluginsdk.ResourceData) *[]virtualmachineinstances.VirtualDisk {
	result := make([]virtualmachineinstances.VirtualDisk, 0)
	if len(input) == 0 {
		return &result
	}

	for k, v := range input {
		virtualDisk := virtualmachineinstances.VirtualDisk{}

		// As API allows zero value for this property, so TF has to use d.GetRawConfig() to determine whether it's set in tf config
		if bus := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["bus"]; !bus.IsNull() {
			virtualDisk.Bus = pointer.To(v.Bus)
		}

		if lun := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["lun"]; !lun.IsNull() {
			virtualDisk.Lun = pointer.To(v.Lun)
		}

		if busType := v.BusType; busType != "" {
			virtualDisk.BusType = pointer.To(busType)
		}

		if diskSizeGB := v.DiskSizeGB; diskSizeGB != 0 {
			virtualDisk.DiskSizeGB = pointer.To(diskSizeGB)
		}

		if name := v.Name; name != "" {
			virtualDisk.Name = pointer.To(name)
		}

		if templateDiskId := v.TemplateDiskId; templateDiskId != "" {
			virtualDisk.TemplateDiskId = pointer.To(templateDiskId)
		}

		if vhdType := v.VhdType; vhdType != "" {
			virtualDisk.VhdType = pointer.To(vhdType)
		}

		if storageQosPolicyName := v.StorageQoSPolicyName; storageQosPolicyName != "" {
			virtualDisk.StorageQoSPolicy = &virtualmachineinstances.StorageQosPolicyDetails{
				Name: pointer.To(storageQosPolicyName),
			}
		}

		result = append(result, virtualDisk)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input []string) (*[]virtualmachineinstances.AvailabilitySetListItem, error) {
	result := make([]virtualmachineinstances.AvailabilitySetListItem, 0)
	if len(input) == 0 {
		return &result, nil
	}

	for _, v := range input {
		availabilitySetId, err := availabilitysets.ParseAvailabilitySetID(v)
		if err != nil {
			return nil, err
		}

		result = append(result, virtualmachineinstances.AvailabilitySetListItem{
			Id:   pointer.To(availabilitySetId.ID()),
			Name: pointer.To(availabilitySetId.AvailabilitySetName),
		})
	}

	return &result, nil
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForUpdate(input []Infrastructure) *virtualmachineinstances.InfrastructureProfileUpdate {
	if len(input) == 0 {
		return nil
	}

	infrastructureProfile := input[0]

	result := virtualmachineinstances.InfrastructureProfileUpdate{}

	if v := infrastructureProfile.CheckpointType; v != "" {
		result.CheckpointType = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(input []Hardware, d *pluginsdk.ResourceData) *virtualmachineinstances.HardwareProfileUpdate {
	result := virtualmachineinstances.HardwareProfileUpdate{}

	if len(input) == 0 {
		return &result
	}

	hardwareProfile := input[0]

	// As TF always sets bool value to false when it isn't set, so it has to use d.GetRawConfig() to determine whether it is set in the tf config
	if v := d.GetRawConfig().AsValueMap()["hardware"].AsValueSlice()[0].AsValueMap()["limit_cpu_for_migration_enabled"]; !v.IsNull() {
		result.LimitCPUForMigration = pointer.ToEnum[virtualmachineinstances.LimitCPUForMigration](strconv.FormatBool(hardwareProfile.LimitCpuForMigrationEnabled))
	}

	if v := hardwareProfile.CpuCount; v != 0 {
		result.CpuCount = pointer.To(v)
	}

	dynamicMemoryEnabled := hardwareProfile.DynamicMemoryMaxInMb != 0 || hardwareProfile.DynamicMemoryMinInMb != 0

	result.DynamicMemoryEnabled = pointer.ToEnum[virtualmachineinstances.DynamicMemoryEnabled](strconv.FormatBool(dynamicMemoryEnabled))

	if v := hardwareProfile.DynamicMemoryMaxInMb; v != 0 {
		result.DynamicMemoryMaxMB = pointer.To(v)
	}

	if v := hardwareProfile.DynamicMemoryMinInMb; v != 0 {
		result.DynamicMemoryMinMB = pointer.To(v)
	}

	if v := hardwareProfile.MemoryInMb; v != 0 {
		result.MemoryMB = pointer.To(v)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForUpdate(input []NetworkInterface) *[]virtualmachineinstances.NetworkInterfaceUpdate {
	result := make([]virtualmachineinstances.NetworkInterfaceUpdate, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterface := virtualmachineinstances.NetworkInterfaceUpdate{
			Name: pointer.To(v.Name),
		}

		if ipv4AddressType := v.Ipv4AddressType; ipv4AddressType != "" {
			networkInterface.IPv4AddressType = pointer.ToEnum[virtualmachineinstances.AllocationMethod](ipv4AddressType)
		}

		if ipv6AddressType := v.Ipv6AddressType; ipv6AddressType != "" {
			networkInterface.IPv6AddressType = pointer.ToEnum[virtualmachineinstances.AllocationMethod](ipv6AddressType)
		}

		if macAddressType := v.MacAddressType; macAddressType != "" {
			networkInterface.MacAddressType = pointer.ToEnum[virtualmachineinstances.AllocationMethod](macAddressType)
		}

		if vnetId := v.VirtualNetworkId; vnetId != "" {
			networkInterface.VirtualNetworkId = pointer.To(vnetId)
		}

		result = append(result, networkInterface)
	}

	return &result
}

func expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForUpdate(input []StorageDisk, d *pluginsdk.ResourceData) *[]virtualmachineinstances.VirtualDiskUpdate {
	result := make([]virtualmachineinstances.VirtualDiskUpdate, 0)
	if len(input) == 0 {
		return &result
	}

	for k, v := range input {
		virtualDisk := virtualmachineinstances.VirtualDiskUpdate{}

		// As API allows zero value for this property, so TF has to use d.GetRawConfig() to determine whether it's set in tf config
		if bus := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["bus"]; !bus.IsNull() {
			virtualDisk.Bus = pointer.To(v.Bus)
		}

		if lun := d.GetRawConfig().AsValueMap()["storage_disk"].AsValueSlice()[k].AsValueMap()["lun"]; !lun.IsNull() {
			virtualDisk.Lun = pointer.To(v.Lun)
		}

		if busType := v.BusType; busType != "" {
			virtualDisk.BusType = pointer.To(busType)
		}

		if diskSizeGB := v.DiskSizeGB; diskSizeGB != 0 {
			virtualDisk.DiskSizeGB = pointer.To(diskSizeGB)
		}

		if name := v.Name; name != "" {
			virtualDisk.Name = pointer.To(name)
		}

		if vhdType := v.VhdType; vhdType != "" {
			virtualDisk.VhdType = pointer.To(vhdType)
		}

		if storageQosPolicyName := v.StorageQoSPolicyName; storageQosPolicyName != "" {
			virtualDisk.StorageQoSPolicy = &virtualmachineinstances.StorageQosPolicyDetails{
				Name: pointer.To(storageQosPolicyName),
			}
		}

		result = append(result, virtualDisk)
	}

	return &result
}
