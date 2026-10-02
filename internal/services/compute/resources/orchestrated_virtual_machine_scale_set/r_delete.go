// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orchestrated_virtual_machine_scale_set

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceOrchestratedVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	resp, err := client.Get(ctx, *id, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving Orchestrated %s: %w", id, err)
	}

	// Sometimes VMSS's aren't fully deleted when the `Delete` call returns - as such we'll try to scale the cluster
	// to 0 nodes first, then delete the cluster - which should ensure there's no Network Interfaces kicking around
	// and work around this Azure API bug:
	// Original Error: Code="InUseSubnetCannotBeDeleted" Message="Subnet internal is in use by
	// /{nicResourceID}/|providers|Microsoft.Compute|virtualMachineScaleSets|acctestvmss-190923101253410278|virtualMachines|0|networkInterfaces|example/ipConfigurations/internal and cannot be deleted.
	// In order to delete the subnet, delete all the resources within the subnet. See aka.ms/deletesubnet.
	if resp.Model != nil && resp.Model.Sku != nil {
		resp.Model.Sku.Capacity = pointer.To(int64(0))

		log.Printf("[DEBUG] Scaling instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		update := virtualmachinescalesets.VirtualMachineScaleSetUpdate{
			Sku: resp.Model.Sku,
		}
		if err := client.UpdateThenPoll(ctx, *id, update, virtualmachinescalesets.DefaultUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating number of instances in Orchestrated %s to scale to 0: %w", id, err)
		}
		log.Printf("[DEBUG] Scaled instances to 0 prior to deletion - this helps avoids networking issues within Azure")
	} else {
		log.Printf("[DEBUG] Unable to scale instances to `0` since the `sku` block is nil - trying to delete anyway")
	}

	// @ArcturusZhang (mimicking from windows_virtual_machine_pluginsdk.go): sending `nil` here omits this value from being sent
	// which matches the previous behaviour - we're only splitting this out so it's clear why
	if err = client.DeleteThenPoll(ctx, *id, virtualmachinescalesets.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting Orchestrated %s: %w", id, err)
	}

	return nil
}
