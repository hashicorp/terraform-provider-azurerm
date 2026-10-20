// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitPeeringRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("peering_type", id.PeeringName)
	d.Set("express_route_circuit_name", id.CircuitName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("azure_asn", props.AzureASN)
			d.Set("peer_asn", props.PeerASN)
			d.Set("primary_azure_port", props.PrimaryAzurePort)
			d.Set("secondary_azure_port", props.SecondaryAzurePort)
			d.Set("primary_peer_address_prefix", props.PrimaryPeerAddressPrefix)
			d.Set("secondary_peer_address_prefix", props.SecondaryPeerAddressPrefix)
			d.Set("vlan_id", props.VlanId)
			d.Set("gateway_manager_etag", props.GatewayManagerEtag)
			d.Set("ipv4_enabled", pointer.From(props.State) == expressroutecircuitpeerings.ExpressRoutePeeringStateEnabled)

			routeFilterId := ""
			if props.RouteFilter != nil && props.RouteFilter.Id != nil {
				routeFilterId = *props.RouteFilter.Id
			}
			d.Set("route_filter_id", routeFilterId)

			if err := d.Set("microsoft_peering_config", flattenExpressRouteCircuitPeeringMicrosoftConfig(props.MicrosoftPeeringConfig)); err != nil {
				return fmt.Errorf("setting `microsoft_peering_config`: %+v", err)
			}
			if err := d.Set("ipv6", flattenExpressRouteCircuitIpv6PeeringConfig(props.IPv6PeeringConfig)); err != nil {
				return fmt.Errorf("setting `ipv6`: %+v", err)
			}
		}
	}

	return nil
}

func flattenExpressRouteCircuitPeeringMicrosoftConfig(input *expressroutecircuitpeerings.ExpressRouteCircuitPeeringConfig) any {
	if input == nil {
		return []any{}
	}

	config := make(map[string]any)
	prefixes := make([]string, 0)
	if customerASN := input.CustomerASN; customerASN != nil {
		config["customer_asn"] = *customerASN
	}
	if routingRegistryName := input.RoutingRegistryName; routingRegistryName != nil {
		config["routing_registry_name"] = *routingRegistryName
	}
	if ps := input.AdvertisedPublicPrefixes; ps != nil {
		prefixes = *ps
	}

	if community := input.AdvertisedCommunities; community != nil {
		config["advertised_communities"] = *community
	}
	config["advertised_public_prefixes"] = prefixes

	return []any{config}
}

func flattenExpressRouteCircuitIpv6PeeringConfig(input *expressroutecircuitpeerings.IPv6ExpressRouteCircuitPeeringConfig) []any {
	if input == nil {
		return make([]any, 0)
	}

	routeFilterId := ""
	if input.RouteFilter != nil && input.RouteFilter.Id != nil {
		routeFilterId = *input.RouteFilter.Id
	}
	return []any{
		map[string]any{
			"microsoft_peering":             flattenExpressRouteCircuitPeeringMicrosoftConfig(input.MicrosoftPeeringConfig),
			"primary_peer_address_prefix":   pointer.From(input.PrimaryPeerAddressPrefix),
			"secondary_peer_address_prefix": pointer.From(input.SecondaryPeerAddressPrefix),
			"route_filter_id":               routeFilterId,
			"enabled":                       pointer.From(input.State) == expressroutecircuitpeerings.ExpressRouteCircuitPeeringStateEnabled,
		},
	}
}

func convertCircuitConnectionsToPeeringConnections(input *[]expressroutecircuitconnections.ExpressRouteCircuitConnection) *[]expressroutecircuitpeerings.ExpressRouteCircuitConnection {
	output := make([]expressroutecircuitpeerings.ExpressRouteCircuitConnection, 0)

	if input == nil || len(*input) == 0 {
		return &output
	}

	for _, i := range *input {
		o := expressroutecircuitpeerings.ExpressRouteCircuitConnection{
			Etag: i.Etag,
			Id:   i.Id,
			Name: i.Name,
			Type: i.Type,
		}

		if props := i.Properties; props != nil {
			o.Properties = &expressroutecircuitpeerings.ExpressRouteCircuitConnectionPropertiesFormat{
				AddressPrefix:           props.AddressPrefix,
				AuthorizationKey:        props.AuthorizationKey,
				CircuitConnectionStatus: (*expressroutecircuitpeerings.CircuitConnectionStatus)(props.CircuitConnectionStatus),
				ExpressRouteCircuitPeering: &expressroutecircuitpeerings.SubResource{
					Id: props.ExpressRouteCircuitPeering.Id,
				},
				IPv6CircuitConnectionConfig: &expressroutecircuitpeerings.IPv6CircuitConnectionConfig{
					AddressPrefix:           props.IPv6CircuitConnectionConfig.AddressPrefix,
					CircuitConnectionStatus: (*expressroutecircuitpeerings.CircuitConnectionStatus)(props.IPv6CircuitConnectionConfig.CircuitConnectionStatus),
				},
				PeerExpressRouteCircuitPeering: &expressroutecircuitpeerings.SubResource{
					Id: props.PeerExpressRouteCircuitPeering.Id,
				},
				ProvisioningState: (*expressroutecircuitpeerings.ProvisioningState)(props.ProvisioningState),
			}
		}
		output = append(output, o)
	}
	return &output
}
