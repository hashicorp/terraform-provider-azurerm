package compute

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func importVirtualMachine(osType compute.OperatingSystemTypes, resourceType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.VirtualMachineID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Compute.VMClient
		vm, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("Error retrieving Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if vm.VirtualMachineProperties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("Error retrieving Virtual Machine %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
		}

		isCorrectOS := false
		if profile := vm.VirtualMachineProperties.StorageProfile; profile != nil {
			if profile.OsDisk != nil && profile.OsDisk.OsType == osType {
				isCorrectOS = true
			}

			if profile.OsDisk.Vhd != nil {
				return []*pluginsdk.ResourceData{}, fmt.Errorf("The %q resource only supports Managed Disks - please use the `azurerm_virtual_machine` resource for Unmanaged Disks", resourceType)
			}
		}

		if !isCorrectOS {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("The %q resource only supports %s Virtual Machines", resourceType, string(osType))
		}

		// we don't support VM's without an OS Profile / attach
		if vm.VirtualMachineProperties.OsProfile == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("The %q resource doesn't support attaching OS Disks - please use the `azurerm_virtual_machine` resource instead", resourceType)
		}

		hasSshKeys := false
		if osType == compute.Linux {
			if linux := vm.VirtualMachineProperties.OsProfile.LinuxConfiguration; linux != nil {
				if linux.SSH != nil && linux.SSH.PublicKeys != nil {
					hasSshKeys = len(*linux.SSH.PublicKeys) > 0
				}
			}
		}

		if !hasSshKeys {
			d.Set("admin_password", "ignored-as-imported")
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
