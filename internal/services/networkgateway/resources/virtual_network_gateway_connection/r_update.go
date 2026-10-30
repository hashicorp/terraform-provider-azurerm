// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.VirtualNetworkGatewayConnections
	vnetGatewayClient := meta.(*clients.Client).NetworkGateway.VirtualNetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgatewayconnections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	var virtualNetworkGateway virtualnetworkgateways.VirtualNetworkGateway
	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)

		gwid, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(virtualNetworkGatewayId)
		if err != nil {
			return err
		}

		resp, err := vnetGatewayClient.Get(ctx, *gwid)
		if err != nil {
			return err
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: %+v", gwid, err)
		}

		virtualNetworkGateway = *resp.Model
	}

	if d.HasChange("authorization_key") {
		payload.Properties.AuthorizationKey = pointer.To(d.Get("authorization_key").(string))
	}

	if d.HasChange("egress_nat_rule_ids") {
		payload.Properties.EgressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(d.Get("egress_nat_rule_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("ingress_nat_rule_ids") {
		payload.Properties.IngressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(d.Get("ingress_nat_rule_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("local_network_gateway_id") {
		localNetworkGatewayId := d.Get("local_network_gateway_id").(string)
		name, err := localNetworkGatewayFromId(localNetworkGatewayId)
		if err != nil {
			return fmt.Errorf("getting LocalNetworkGateway Name and Group: %+v", err)
		}

		payload.Properties.LocalNetworkGateway2 = &virtualnetworkgatewayconnections.LocalNetworkGateway{
			Id:   &localNetworkGatewayId,
			Name: &name,
			Properties: virtualnetworkgatewayconnections.LocalNetworkGatewayPropertiesFormat{
				LocalNetworkAddressSpace: &virtualnetworkgatewayconnections.AddressSpace{},
			},
		}
	}

	if d.HasChange("bgp_enabled") {
		payload.Properties.EnableBgp = pointer.To(d.Get("bgp_enabled").(bool))
	}

	if d.HasChange("use_policy_based_traffic_selectors") {
		payload.Properties.UsePolicyBasedTrafficSelectors = pointer.To(d.Get("use_policy_based_traffic_selectors").(bool))
	}

	if d.HasChange("routing_weight") {
		payload.Properties.RoutingWeight = pointer.To(int64(d.Get("routing_weight").(int)))
	}

	if d.HasChange("express_route_gateway_bypass") {
		payload.Properties.ExpressRouteGatewayBypass = pointer.To(d.Get("express_route_gateway_bypass").(bool))
	}

	if d.HasChange("private_link_fast_path_enabled") {
		payload.Properties.EnablePrivateLinkFastPath = pointer.To(d.Get("private_link_fast_path_enabled").(bool))
	}

	if d.HasChange("traffic_selector_policy") {
		payload.Properties.TrafficSelectorPolicies = expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(d.Get("traffic_selector_policy").([]any))
	}

	if d.HasChange("custom_bgp_addresses") {
		if virtualNetworkGateway.Properties.BgpSettings == nil || virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses == nil {
			return fmt.Errorf("retrieving BGP peering address from `virtual_network_gateway` %s (%s) failed: get nil", *virtualNetworkGateway.Name, *virtualNetworkGateway.Id)
		}

		gatewayCustomBgpIPAddresses, err := expandGatewayCustomBgpIPAddresses(d, virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses)
		if err != nil {
			return err
		}

		payload.Properties.GatewayCustomBgpIPAddresses = gatewayCustomBgpIPAddresses
	}

	if d.HasChange("ipsec_policy") {
		payload.Properties.IPsecPolicies = expandVirtualNetworkGatewayConnectionIpsecPolicies(d.Get("ipsec_policy").([]any))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	// Only send the update if there are changes other than the shared key
	if d.HasChangesExcept("shared_key") {
		if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	if d.HasChange("shared_key") {
		if err := client.SetSharedKeyThenPoll(ctx, *id, virtualnetworkgatewayconnections.ConnectionSharedKey{
			Value: d.Get("shared_key").(string),
		}); err != nil {
			return fmt.Errorf("updating Shared Key for %s: %+v", id, err)
		}

		// Once this issue https://github.com/Azure/azure-rest-api-specs/issues/26660 is fixed, below this part will be removed
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{string(virtualnetworkgatewayconnections.ProvisioningStateUpdating)},
			Target:     []string{string(virtualnetworkgatewayconnections.ProvisioningStateSucceeded)},
			Refresh:    virtualNetworkGatewayConnectionStateRefreshFunc(ctx, client, *id),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayConnectionRead(d, meta)
}
