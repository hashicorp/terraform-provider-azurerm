// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRouteTableUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	virtHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.VirtualHubName, virtualHubResourceName)

	id, err := virtualwans.ParseHubRouteTableID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.HubRouteTablesGet(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
	}

	payload := existing.Model

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("labels") {
		payload.Properties.Labels = pluginsdk.ExpandStringSlice(d.Get("labels").(*pluginsdk.Set).List())
	}

	if d.HasChange("route") {
		payload.Properties.Routes = expandVirtualHubRouteTableHubRoutes(d.Get("route").(*pluginsdk.Set).List())
	}

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRead(d, meta)
}
