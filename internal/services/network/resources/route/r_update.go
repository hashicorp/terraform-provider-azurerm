// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Routes
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
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

	rtId := routetables.NewRouteTableID(id.SubscriptionId, id.ResourceGroupName, id.RouteTableName)
	locks.ByID(rtId.ID())
	defer locks.UnlockByID(rtId.ID())

	if d.HasChange("address_prefix") {
		payload.Properties.AddressPrefix = pointer.To(d.Get("address_prefix").(string))
	}

	if d.HasChange("next_hop_type") {
		payload.Properties.NextHopType = routes.RouteNextHopType(d.Get("next_hop_type").(string))
	}

	if d.HasChange("next_hop_in_ip_address") {
		payload.Properties.NextHopIPAddress = pointer.To(d.Get("next_hop_in_ip_address").(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceRouteRead(d, meta)
}
