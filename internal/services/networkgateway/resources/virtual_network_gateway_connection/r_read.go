// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.VirtualNetworkGatewayConnections
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgatewayconnections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)

	respKey, err := client.GetSharedKey(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Shared Key for %s: %+v", id, err)
	}

	if model := respKey.Model; model != nil {
		if model.Value != "" {
			d.Set("shared_key", model.Value)
		}
	}

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties

		if string(props.ConnectionType) != "" {
			d.Set("type", string(props.ConnectionType))
		}

		d.Set("virtual_network_gateway_id", props.VirtualNetworkGateway1.Id)

		// not returned from api - getting from state
		d.Set("authorization_key", d.Get("authorization_key").(string))

		if props.DpdTimeoutSeconds != nil {
			d.Set("dpd_timeout_seconds", props.DpdTimeoutSeconds)
		}

		if props.Peer != nil {
			d.Set("express_route_circuit_id", props.Peer.Id)
		}

		if props.VirtualNetworkGateway2 != nil {
			d.Set("peer_virtual_network_gateway_id", props.VirtualNetworkGateway2.Id)
		}

		if props.UseLocalAzureIPAddress != nil {
			d.Set("local_azure_ip_address_enabled", props.UseLocalAzureIPAddress)
		}

		if props.LocalNetworkGateway2 != nil {
			d.Set("local_network_gateway_id", props.LocalNetworkGateway2.Id)
		}

		d.Set("bgp_enabled", pointer.From(props.EnableBgp))

		if props.UsePolicyBasedTrafficSelectors != nil {
			d.Set("use_policy_based_traffic_selectors", props.UsePolicyBasedTrafficSelectors)
		}

		if props.RoutingWeight != nil {
			d.Set("routing_weight", props.RoutingWeight)
		}

		if err := d.Set("custom_bgp_addresses", flattenGatewayCustomBgpIPAddresses(props.GatewayCustomBgpIPAddresses)); err != nil {
			return fmt.Errorf("setting `custom_bgp_addresses`: %+v", err)
		}

		d.Set("connection_protocol", pointer.FromEnum(props.ConnectionProtocol))

		d.Set("connection_mode", pointer.FromEnum(props.ConnectionMode))

		if props.ExpressRouteGatewayBypass != nil {
			d.Set("express_route_gateway_bypass", props.ExpressRouteGatewayBypass)
		}

		if props.EnablePrivateLinkFastPath != nil {
			d.Set("private_link_fast_path_enabled", props.EnablePrivateLinkFastPath)
		}

		if props.IPsecPolicies != nil {
			if err := d.Set("ipsec_policy", flattenVirtualNetworkGatewayConnectionIpsecPolicies(props.IPsecPolicies)); err != nil {
				return fmt.Errorf("setting `ipsec_policy`: %+v", err)
			}
		}

		if err := d.Set("traffic_selector_policy", flattenVirtualNetworkGatewayConnectionTrafficSelectorPolicies(props.TrafficSelectorPolicies)); err != nil {
			return fmt.Errorf("setting `traffic_selector_policy`: %+v", err)
		}

		if err := d.Set("egress_nat_rule_ids", flattenVirtualNetworkGatewayConnectionNatRuleIds(props.EgressNatRules)); err != nil {
			return fmt.Errorf("setting `egress_nat_rule_ids`: %+v", err)
		}

		if err := d.Set("ingress_nat_rule_ids", flattenVirtualNetworkGatewayConnectionNatRuleIds(props.IngressNatRules)); err != nil {
			return fmt.Errorf("setting `ingress_nat_rule_ids`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func virtualNetworkGatewayConnectionStateRefreshFunc(ctx context.Context, client *virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionsClient, id virtualnetworkgatewayconnections.ConnectionId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties.ProvisioningState != nil {
			return res, pointer.FromEnum(res.Model.Properties.ProvisioningState), nil
		}

		return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
	}
}
