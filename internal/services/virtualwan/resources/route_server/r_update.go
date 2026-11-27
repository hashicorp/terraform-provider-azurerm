// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_server

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

func resourceRouteServerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, "azurerm_route_server")
	defer locks.UnlockByName(id.VirtualHubName, "azurerm_route_server")

	existing, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
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

	if d.HasChange("hub_routing_preference") {
		payload.Properties.HubRoutingPreference = pointer.ToEnum[virtualwans.HubRoutingPreference](d.Get("hub_routing_preference").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.VirtualHubsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning", "Updating"},
		Target:                    []string{"Succeeded", "Provisioned"},
		Refresh:                   routeServerCreateRefreshFunc(ctx, client, *id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceRouteServerRead(d, meta)
}
