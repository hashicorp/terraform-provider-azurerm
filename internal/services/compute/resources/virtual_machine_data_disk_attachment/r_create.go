// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_data_disk_attachment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineDataDiskAttachmentCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedVirtualMachineId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(parsedVirtualMachineId.VirtualMachineName, helpers.VirtualMachineResourceName)
	defer locks.UnlockByName(parsedVirtualMachineId.VirtualMachineName, helpers.VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, *parsedVirtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			return fmt.Errorf("%s was not found", parsedVirtualMachineId)
		}
		return fmt.Errorf("retrieving %s: %+v", parsedVirtualMachineId, err)
	}

	if virtualMachine.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", parsedVirtualMachineId)
	}
	if virtualMachine.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", parsedVirtualMachineId)
	}
	if virtualMachine.Model.Properties.StorageProfile == nil {
		return fmt.Errorf("retrieving %s: `storageprofile` was nil", parsedVirtualMachineId)
	}

	managedDiskId := d.Get("managed_disk_id").(string)
	managedDisk, err := retrieveDataDiskAttachmentManagedDisk(d, meta, managedDiskId)
	if err != nil {
		return fmt.Errorf("retrieving Managed Disk %q: %+v", managedDiskId, err)
	}

	if managedDisk.Sku == nil {
		return fmt.Errorf("unable to determine Storage Account Type for Managed Disk %q: %+v", managedDiskId, err)
	}

	name := *managedDisk.Name
	resourceId := fmt.Sprintf("%s/dataDisks/%s", parsedVirtualMachineId.ID(), name)
	lun := int32(d.Get("lun").(int))
	caching := d.Get("caching").(string)
	createOption := virtualmachines.DiskCreateOptionTypes(d.Get("create_option").(string))
	writeAcceleratorEnabled := d.Get("write_accelerator_enabled").(bool)

	expandedDisk := virtualmachines.DataDisk{
		Name:         pointer.To(name),
		Caching:      pointer.ToEnum[virtualmachines.CachingTypes](caching),
		CreateOption: createOption,
		Lun:          int64(lun),
		ManagedDisk: &virtualmachines.ManagedDiskParameters{
			Id:                 pointer.To(managedDiskId),
			StorageAccountType: pointer.ToEnum[virtualmachines.StorageAccountTypes](string(*managedDisk.Sku.Name)),
		},
		WriteAcceleratorEnabled: pointer.To(writeAcceleratorEnabled),
	}

	// there are ways to provision a VM without a StorageProfile and/or DataDisks
	if virtualMachine.Model.Properties.StorageProfile == nil {
		virtualMachine.Model.Properties.StorageProfile = &virtualmachines.StorageProfile{}
	}

	if virtualMachine.Model.Properties.StorageProfile.DataDisks == nil {
		virtualMachine.Model.Properties.StorageProfile.DataDisks = pointer.To(make([]virtualmachines.DataDisk, 0))
	}

	disks := *virtualMachine.Model.Properties.StorageProfile.DataDisks

	existingIndex := -1
	for i, disk := range disks {
		if *disk.Name == name {
			existingIndex = i
			break
		}
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			if existingIndex != -1 {
				return tf.ImportAsExistsError("azurerm_virtual_machine_data_disk_attachment", resourceId)
			}

			disks = append(disks, expandedDisk)
		}
	} else {
		if existingIndex == -1 {
			return fmt.Errorf("unable to find Disk %q attached to Virtual Machine %q ", name, parsedVirtualMachineId.String())
		}

		disks[existingIndex] = expandedDisk
	}

	virtualMachine.Model.Properties.StorageProfile.DataDisks = &disks

	// fixes #2485
	virtualMachine.Model.Identity = nil
	// fixes #1600
	virtualMachine.Model.Resources = nil
	// fixes #24145
	virtualMachine.Model.Properties.ApplicationProfile = nil

	// if there's too many disks we get a 409 back with:
	//   `The maximum number of data disks allowed to be attached to a VM of this size is 1.`
	// which we're intentionally not wrapping, since the errors good.
	// TODO: implement `CallbackThenPoll` on Create, requires updated resource ID implementing `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, *parsedVirtualMachineId, *virtualMachine.Model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s with Disk %q: %+v", parsedVirtualMachineId, name, err)
	}

	d.SetId(resourceId)
	return resourceVirtualMachineDataDiskAttachmentRead(d, meta)
}
