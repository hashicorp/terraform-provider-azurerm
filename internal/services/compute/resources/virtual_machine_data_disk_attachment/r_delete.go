// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_data_disk_attachment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineDataDiskAttachmentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataDiskID(d.Id())
	if err != nil {
		return err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)

	locks.ByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			return fmt.Errorf("%s was not found", virtualMachineId)
		}

		return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
	}

	if virtualMachine.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", virtualMachineId)
	}
	if virtualMachine.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", virtualMachineId)
	}
	if virtualMachine.Model.Properties.StorageProfile == nil {
		return fmt.Errorf("retrieving %s: `storageprofile` was nil", virtualMachineId)
	}

	dataDisks := make([]virtualmachines.DataDisk, 0)
	for _, dataDisk := range *virtualMachine.Model.Properties.StorageProfile.DataDisks {
		// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
		if *dataDisk.Name != id.Name {
			dataDisks = append(dataDisks, dataDisk)
		}
	}

	virtualMachine.Model.Properties.StorageProfile.DataDisks = &dataDisks

	// fixes #2485
	virtualMachine.Model.Identity = nil
	// fixes #1600
	virtualMachine.Model.Resources = nil
	// fixes #24145
	virtualMachine.Model.Properties.ApplicationProfile = nil

	if err := client.CreateOrUpdateThenPoll(ctx, virtualMachineId, *virtualMachine.Model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("removing %s from Virtual Machine %q : %+v", id, id.VirtualMachineName, err)
	}

	return nil
}
