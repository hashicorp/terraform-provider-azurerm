// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVpnGatewayConnectionResourceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	gatewayId, err := virtualwans.ParseVpnGatewayID(d.Get("vpn_gateway_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewVPNConnectionID(gatewayId.SubscriptionId, gatewayId.ResourceGroupName, gatewayId.VpnGatewayName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.VpnConnectionsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_gateway_connection", id.ID())
		}
	}

	locks.ByName(gatewayId.VpnGatewayName, VPNGatewayResourceName)
	defer locks.UnlockByName(gatewayId.VpnGatewayName, VPNGatewayResourceName)

	payload := virtualwans.VpnConnection{
		Properties: &virtualwans.VpnConnectionProperties{
			EnableInternetSecurity: pointer.To(d.Get("internet_security_enabled").(bool)),
			RemoteVpnSite: &virtualwans.SubResource{
				Id: pointer.To(d.Get("remote_vpn_site_id").(string)),
			},
			VpnLinkConnections:   expandVpnGatewayConnectionVpnSiteLinkConnections(d.Get("vpn_link").([]any)),
			RoutingConfiguration: expandVpnGatewayConnectionRoutingConfiguration(d.Get("routing").([]any)),
		},
	}

	if v, ok := d.GetOk("traffic_selector_policy"); ok {
		payload.Properties.TrafficSelectorPolicies = expandVpnGatewayConnectionTrafficSelectorPolicy(v.(*pluginsdk.Set).List())
	}

	if err := client.VpnConnectionsCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVpnGatewayConnectionResourceRead(d, meta)
}
