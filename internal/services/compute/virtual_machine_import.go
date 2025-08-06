// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importVirtualMachine(osType virtualmachines.OperatingSystemTypes, resourceType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := virtualmachines.ParseVirtualMachineID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Compute.VirtualMachinesClient
		vm, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if model := vm.Model; model == nil || model.Properties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties` was nil", id)
		}

		isCorrectOS := false
		if profile := vm.Model.Properties.StorageProfile; profile != nil {
			if profile.OsDisk != nil && profile.OsDisk.OsType != nil && *profile.OsDisk.OsType == osType {
				isCorrectOS = true
			}

			if profile.OsDisk.Vhd != nil {
				return []*pluginsdk.ResourceData{}, fmt.Errorf("the %q resource only supports Managed Disks - please use the `azurerm_virtual_machine` resource for Unmanaged Disks", resourceType)
			}
		}

		if !isCorrectOS {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("the %q resource only supports %s Virtual Machines", resourceType, string(osType))
		}

		// we don't support VM's without an OS Profile / attach
		if vm.Model.Properties.OsProfile == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("the %q resource doesn't support attaching OS Disks - please use the `azurerm_virtual_machine` resource instead", resourceType)
		}

		hasSshKeys := false
		if osType == virtualmachines.OperatingSystemTypesLinux {
			if linux := vm.Model.Properties.OsProfile.LinuxConfiguration; linux != nil {
				if linux.Ssh != nil && linux.Ssh.PublicKeys != nil {
					hasSshKeys = len(*linux.Ssh.PublicKeys) > 0
				}
			}
		}

		if !hasSshKeys {
			d.Set("admin_password", "ignored-as-imported")
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
