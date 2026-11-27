// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// @ArcturusZhang (mimicking from virtual_machine_pluginsdk.go): sending `nil` here omits this value from being sent
	// which matches the previous behaviour - we're only splitting this out so it's clear why
	opts := virtualmachinescalesets.DefaultDeleteOperationOptions()
	opts.ForceDeletion = nil
	if err := client.DeleteThenPoll(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
