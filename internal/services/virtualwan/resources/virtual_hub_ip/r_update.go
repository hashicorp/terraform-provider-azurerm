// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_ip

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubIPUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtualHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubId.VirtualHubName, virtualHubResourceName)

	id, err := commonids.ParseVirtualHubIPConfigurationID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.VirtualHubIPConfigurationGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("private_ip_address") {
		payload.Properties.PrivateIPAddress = pointer.To(d.Get("private_ip_address").(string))
	}

	if d.HasChange("private_ip_allocation_method") {
		payload.Properties.PrivateIPAllocationMethod = pointer.ToEnum[virtualwans.IPAllocationMethod](d.Get("private_ip_allocation_method").(string))
	}

	if err := client.VirtualHubIPConfigurationCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubIPRead(d, meta)
}
