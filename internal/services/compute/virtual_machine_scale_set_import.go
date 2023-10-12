// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func importOrchestratedVirtualMachineScaleSet(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
	id, err := commonids.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{}, err
	}

	client := meta.(*clients.Client).Compute.VMScaleSetClient
	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	_, err = client.Get(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return []*pluginsdk.ResourceData{d}, nil
}

func importVirtualMachineScaleSet(osType compute.OperatingSystemTypes, resourceType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := commonids.ParseVirtualMachineScaleSetID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Compute.VMScaleSetClient
		// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
		vm, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, compute.ExpandTypesForGetVMScaleSetsUserData)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if vm.VirtualMachineScaleSetProperties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties` was nil", id)
		}

		if vm.VirtualMachineScaleSetProperties.VirtualMachineProfile == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties.virtualMachineProfile` was nil", id)
		}

		if vm.VirtualMachineScaleSetProperties.VirtualMachineProfile.OsProfile == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties.virtualMachineProfile.osProfile` was nil", id)
		}

		isCorrectOS := false
		hasSshKeys := false
		if profile := vm.VirtualMachineScaleSetProperties.VirtualMachineProfile.OsProfile; profile != nil {
			if profile.LinuxConfiguration != nil && osType == compute.OperatingSystemTypesLinux {
				isCorrectOS = true

				if profile.LinuxConfiguration.SSH != nil && profile.LinuxConfiguration.SSH.PublicKeys != nil {
					hasSshKeys = len(*profile.LinuxConfiguration.SSH.PublicKeys) > 0
				}
			}

			if profile.WindowsConfiguration != nil && osType == compute.OperatingSystemTypesWindows {
				isCorrectOS = true
			}
		}

		if !isCorrectOS {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("The %q resource only supports %s Virtual Machine Scale Sets", resourceType, string(osType))
		}

		if !hasSshKeys {
			d.Set("admin_password", "ignored-as-imported")
		}

		var updatedExtensions []map[string]interface{}
		if vm.VirtualMachineScaleSetProperties.VirtualMachineProfile.ExtensionProfile != nil {
			if extensionsProfile := vm.VirtualMachineScaleSetProperties.VirtualMachineProfile.ExtensionProfile; extensionsProfile != nil {
				for _, v := range *extensionsProfile.Extensions {
					v.ProtectedSettings = ""
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
