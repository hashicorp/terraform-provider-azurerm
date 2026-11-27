// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVpnGatewayConnectionResourceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseVPNConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.VpnConnectionsGet(ctx, *id)
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

	gatewayId := virtualwans.NewVpnGatewayID(id.SubscriptionId, id.ResourceGroupName, id.GatewayName)

	locks.ByName(gatewayId.VpnGatewayName, VPNGatewayResourceName)
	defer locks.UnlockByName(gatewayId.VpnGatewayName, VPNGatewayResourceName)

	if d.HasChange("internet_security_enabled") {
		payload.Properties.EnableInternetSecurity = pointer.To(d.Get("internet_security_enabled").(bool))
	}

	if d.HasChange("routing") {
		payload.Properties.RoutingConfiguration = expandVpnGatewayConnectionRoutingConfiguration(d.Get("routing").([]any))
	}

	if d.HasChange("vpn_link") {
		payload.Properties.VpnLinkConnections = expandVpnGatewayConnectionVpnSiteLinkConnections(d.Get("vpn_link").([]any))
	}

	if d.HasChange("traffic_selector_policy") {
		payload.Properties.TrafficSelectorPolicies = expandVpnGatewayConnectionTrafficSelectorPolicy(d.Get("traffic_selector_policy").(*pluginsdk.Set).List())
	}

	if err = client.VpnConnectionsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVpnGatewayConnectionResourceRead(d, meta)
}
