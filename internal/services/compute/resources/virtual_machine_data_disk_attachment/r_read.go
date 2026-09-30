// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_data_disk_attachment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineDataDiskAttachmentRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataDiskID(d.Id())
	if err != nil {
		return err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)

	virtualMachine, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			log.Printf("[DEBUG] %s was not found therefore Data Disk Attachment cannot exist - removing from state", virtualMachineId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	var disk *virtualmachines.DataDisk
	if model := virtualMachine.Model; model != nil {
		if props := model.Properties; props != nil {
			if profile := props.StorageProfile; profile != nil {
				if dataDisks := profile.DataDisks; dataDisks != nil {
					for _, dataDisk := range *dataDisks {
						// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
						if *dataDisk.Name == id.Name {
							disk = &dataDisk
							break
						}
					}
				}
			}
		}
	}

	if disk == nil {
		log.Printf("[DEBUG] %s was not found on Virtual Machine %q  - removing from state", id, id.VirtualMachineName)
		d.SetId("")
		return nil
	}

	d.Set("virtual_machine_id", virtualMachineId.ID())
	d.Set("caching", pointer.FromEnum(disk.Caching))
	d.Set("create_option", string(disk.CreateOption))
	d.Set("write_accelerator_enabled", disk.WriteAcceleratorEnabled)

	if managedDisk := disk.ManagedDisk; managedDisk != nil {
		d.Set("managed_disk_id", managedDisk.Id)
	}

	d.Set("lun", int(disk.Lun))

	return nil
}

func retrieveDataDiskAttachmentManagedDisk(d *pluginsdk.ResourceData, meta any, id string) (*disks.Disk, error) {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedId, err := commonids.ParseManagedDiskID(id)
	if err != nil {
		return nil, fmt.Errorf("parsing Managed Disk ID %q: %+v", id, err)
	}

	resp, err := client.Get(ctx, *parsedId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("managed disk %q was not found", parsedId.String())
		}

		return nil, fmt.Errorf("making Read request on Azure Managed Disk %q : %+v", parsedId.String(), err)
	}

	return resp.Model, nil
}
