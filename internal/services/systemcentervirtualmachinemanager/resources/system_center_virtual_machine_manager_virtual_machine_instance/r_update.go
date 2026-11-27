// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := virtualmachineinstances.VirtualMachineInstanceUpdate{
				Properties: &virtualmachineinstances.VirtualMachineInstanceUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("infrastructure") {
				parameters.Properties.InfrastructureProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfileForUpdate(model.Infrastructure)
			}

			if metadata.ResourceData.HasChange("hardware") {
				parameters.Properties.HardwareProfile = expandSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfileForUpdate(model.Hardware, metadata.ResourceData)
			}

			if metadata.ResourceData.HasChange("network_interface") {
				parameters.Properties.NetworkProfile = &virtualmachineinstances.NetworkProfileUpdate{
					NetworkInterfaces: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfacesForUpdate(model.NetworkInterfaces),
				}
			}

			if metadata.ResourceData.HasChange("storage_disk") {
				parameters.Properties.StorageProfile = &virtualmachineinstances.StorageProfileUpdate{
					Disks: expandSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisksForUpdate(model.StorageDisks, metadata.ResourceData),
				}
			}

			if metadata.ResourceData.HasChange("system_center_virtual_machine_manager_availability_set_ids") {
				availabilitySets, err := expandSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds)
				if err != nil {
					return err
				}
				parameters.Properties.AvailabilitySets = availabilitySets
			}

			needToRestart := metadata.ResourceData.HasChanges("hardware", "network_interface", "storage_disk")

			if needToRestart {
				if err := client.StopThenPoll(ctx, commonids.NewScopeID(id.Scope), virtualmachineinstances.StopVirtualMachineOptions{
					SkipShutdown: pointer.To(virtualmachineinstances.SkipShutdownFalse),
				}); err != nil {
					return fmt.Errorf("stopping %s: %+v", *id, err)
				}
			}

			if err := client.UpdateThenPoll(ctx, commonids.NewScopeID(id.Scope), parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if needToRestart {
				if err := client.StartThenPoll(ctx, commonids.NewScopeID(id.Scope)); err != nil {
					return fmt.Errorf("starting %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
