// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importOrchestratedVirtualMachineScaleSet(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{}, err
	}

	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	options := virtualmachinescalesets.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	_, err = client.Get(ctx, *id, options)
	if err != nil {
		return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return []*pluginsdk.ResourceData{d}, nil
}

func importVirtualMachineScaleSet(osType virtualmachinescalesets.OperatingSystemTypes, resourceType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
		options := virtualmachinescalesets.DefaultGetOperationOptions()
		options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
		vm, err := client.Get(ctx, *id, options)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if vm.Model == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `model` was nil", id)
		}

		if vm.Model.Properties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties` was nil", id)
		}

		if vm.Model.Properties.VirtualMachineProfile == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties.virtualMachineProfile` was nil", id)
		}

		if vm.Model.Properties.VirtualMachineProfile.OsProfile == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties.virtualMachineProfile.osProfile` was nil", id)
		}

		isCorrectOS := false
		hasSshKeys := false
		if profile := vm.Model.Properties.VirtualMachineProfile.OsProfile; profile != nil {
			if profile.LinuxConfiguration != nil && osType == virtualmachinescalesets.OperatingSystemTypesLinux {
				isCorrectOS = true

				if profile.LinuxConfiguration.Ssh != nil && profile.LinuxConfiguration.Ssh.PublicKeys != nil {
					hasSshKeys = len(*profile.LinuxConfiguration.Ssh.PublicKeys) > 0
				}
			}

			if profile.WindowsConfiguration != nil && osType == virtualmachinescalesets.OperatingSystemTypesWindows {
				isCorrectOS = true
			}
		}

		if !isCorrectOS {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("the %q resource only supports %s Virtual Machine Scale Sets", resourceType, string(osType))
		}

		if !hasSshKeys {
			d.Set("admin_password", "ignored-as-imported")
		}

		var updatedExtensions []map[string]interface{}
		if vm.Model.Properties.VirtualMachineProfile.ExtensionProfile != nil {
			if extensionsProfile := vm.Model.Properties.VirtualMachineProfile.ExtensionProfile; extensionsProfile != nil {
				for _, v := range *extensionsProfile.Extensions {
					if v.Properties != nil {
						v.Properties.ProtectedSettings = nil
					}
				}
				updatedExtensions, err = flattenVirtualMachineScaleSetExtensions(extensionsProfile, d)
				if err != nil {
					return []*pluginsdk.ResourceData{}, fmt.Errorf("could not read VMSS extensions data for %s", id)
				}
			}
		}
		d.Set("extension", updatedExtensions)

		return []*pluginsdk.ResourceData{d}, nil
	}
}
