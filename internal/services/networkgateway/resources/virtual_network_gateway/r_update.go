// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.VirtualNetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return err
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("bgp_enabled") {
		payload.Properties.EnableBgp = pointer.To(d.Get("bgp_enabled").(bool))
	}

	if d.HasChange("active_active") {
		payload.Properties.ActiveActive = pointer.To(d.Get("active_active").(bool))
	}

	if d.HasChange("sku") {
		payload.Properties.Sku = expandVirtualNetworkGatewaySku(d)
	}

	if d.HasChange("ip_configuration") {
		payload.Properties.IPConfigurations = expandVirtualNetworkGatewayIPConfigurations(d)
	}

	if d.HasChange("policy_group") {
		payload.Properties.VirtualNetworkGatewayPolicyGroups = expandVirtualNetworkGatewayPolicyGroups(d.Get("policy_group").([]any))
	}

	if d.HasChange("vpn_client_configuration") {
		payload.Properties.VpnClientConfiguration = expandVirtualNetworkGatewayVpnClientConfig(d, *id)
	}

	if d.HasChange("bgp_settings") {
		bgpSettings, err := expandVirtualNetworkGatewayBgpSettings(*id, d)
		if err != nil {
			return err
		}
		payload.Properties.BgpSettings = bgpSettings
	}

	if d.HasChange("custom_route") {
		payload.Properties.CustomRoutes = expandVirtualNetworkGatewayAddressSpace(d.Get("custom_route").([]any))
	}

	if d.HasChange("default_local_network_gateway_id") {
		payload.Properties.GatewayDefaultSite = nil
		if gatewayDefaultSiteID := d.Get("default_local_network_gateway_id").(string); gatewayDefaultSiteID != "" {
			payload.Properties.GatewayDefaultSite = &virtualnetworkgateways.SubResource{
				Id: &gatewayDefaultSiteID,
			}
		}
	}

	if d.HasChange("bgp_route_translation_for_nat_enabled") {
		payload.Properties.EnableBgpRouteTranslationForNat = pointer.To(d.Get("bgp_route_translation_for_nat_enabled").(bool))
	}

	if d.HasChange("dns_forwarding_enabled") {
		payload.Properties.EnableDnsForwarding = pointer.To(d.Get("dns_forwarding_enabled").(bool))
	}

	if d.HasChange("ip_sec_replay_protection_enabled") {
		payload.Properties.DisableIPSecReplayProtection = pointer.To(!d.Get("ip_sec_replay_protection_enabled").(bool))
	}

	if d.HasChange("remote_vnet_traffic_enabled") {
		payload.Properties.AllowRemoteVnetTraffic = pointer.To(d.Get("remote_vnet_traffic_enabled").(bool))
	}

	if d.HasChange("virtual_wan_traffic_enabled") {
		payload.Properties.AllowVirtualWanTraffic = pointer.To(d.Get("virtual_wan_traffic_enabled").(bool))
	}

	// SKU changes that cross the availability-zone boundary are ForceNew (see
	// resourceVirtualNetworkGatewayCustomizeDiff), so here we only need to push
	// autoscale changes while the gateway stays on the ErGwScale SKU.
	if d.HasChanges("minimum_scale_unit", "maximum_scale_unit") {
		payload.Properties.AutoScaleConfiguration = expandVirtualNetworkGatewayAutoScaleConfiguration(d)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayRead(d, meta)
}
