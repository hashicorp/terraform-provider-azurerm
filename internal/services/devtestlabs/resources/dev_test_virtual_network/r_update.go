// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_virtual_network

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestVirtualNetworkUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworks.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("subnet") {
		payload.Properties.SubnetOverrides = expandDevTestVirtualNetworkSubnets(d.Get("subnet").([]any), subscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestVirtualNetworkRead(d, meta)
}
