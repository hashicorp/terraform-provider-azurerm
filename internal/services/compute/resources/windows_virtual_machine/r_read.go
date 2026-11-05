// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWindowsVirtualMachineRead(d *pluginsdk.ResourceData, meta any) error {
	metaClient := meta.(*clients.Client)
	client := metaClient.Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(metaClient.StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	options := virtualmachines.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachines.InstanceViewTypesUserData)
	resp, err := client.Get(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Windows %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Windows %s: %+v", id, err)
	}

	return resourceWindowsVirtualMachineFlatten(ctx, metaClient, d, id, resp.Model, true)
}
