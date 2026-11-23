// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		return fmt.Errorf("deadline is not properly set for Virtual Hub")
	}

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	existing, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("branch_to_branch_traffic_enabled") {
		payload.Properties.AllowBranchToBranchTraffic = pointer.To(d.Get("branch_to_branch_traffic_enabled").(bool))
	}

	if d.HasChange("route") {
		payload.Properties.RouteTable = expandVirtualHubRoute(d.Get("route").(*pluginsdk.Set).List())
	}

	if d.HasChange("hub_routing_preference") {
		payload.Properties.HubRoutingPreference = pointer.ToEnum[virtualwans.HubRoutingPreference](d.Get("hub_routing_preference").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("virtual_router_auto_scale_min_capacity") {
		if v, ok := d.GetOk("virtual_router_auto_scale_min_capacity"); ok {
			payload.Properties.VirtualRouterAutoScaleConfiguration = &virtualwans.VirtualRouterAutoScaleConfiguration{
				MinCapacity: pointer.To(int64(v.(int))),
			}
		}
	}

	if err := client.VirtualHubsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Hub returns provisioned while the routing state is still "provisioning". This might cause issues with following hubvnet connection operations.
	// https://github.com/Azure/azure-rest-api-specs/issues/10391
	// As a workaround, we will poll the routing state and ensure it is "Provisioned".

	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning"},
		Target:                    []string{"Provisioned", "Failed", "None"},
		Refresh:                   virtualHubCreateRefreshFunc(ctx, client, *id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s provisioning route: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRead(d, meta)
}
