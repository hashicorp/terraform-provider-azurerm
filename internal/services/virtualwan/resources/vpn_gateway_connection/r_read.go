// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_connection

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVpnGatewayConnectionResourceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVPNConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("vpn_gateway_id", virtualwans.NewVpnGatewayID(id.SubscriptionId, id.ResourceGroupName, id.GatewayName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			vpnSiteId := ""
			if props.RemoteVpnSite != nil && props.RemoteVpnSite.Id != nil {
				theVpnSiteId, err := virtualwans.ParseVpnSiteIDInsensitively(*props.RemoteVpnSite.Id)
				if err != nil {
					return err
				}
				vpnSiteId = theVpnSiteId.ID()
			}
			d.Set("remote_vpn_site_id", vpnSiteId)

			d.Set("internet_security_enabled", pointer.From(props.EnableInternetSecurity))

			if err := d.Set("routing", flattenVpnGatewayConnectionRoutingConfiguration(props.RoutingConfiguration)); err != nil {
				return fmt.Errorf(`setting "routing": %v`, err)
			}

			if err := d.Set("vpn_link", flattenVpnGatewayConnectionVpnSiteLinkConnections(props.VpnLinkConnections)); err != nil {
				return fmt.Errorf(`setting "vpn_link": %v`, err)
			}

			if err := d.Set("traffic_selector_policy", flattenVpnGatewayConnectionTrafficSelectorPolicy(props.TrafficSelectorPolicies)); err != nil {
				return fmt.Errorf("setting `traffic_selector_policy`: %+v", err)
			}
		}
	}

	return nil
}
