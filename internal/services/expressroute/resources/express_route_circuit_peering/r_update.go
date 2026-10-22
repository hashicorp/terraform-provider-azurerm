// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_peering

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitPeeringUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.CircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.CircuitName, expressRouteCircuitResourceName)

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	payload := existing.Model

	circuitConnClient := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuitConnections
	circuitConnectionId := commonids.NewExpressRouteCircuitPeeringID(id.SubscriptionId, id.ResourceGroupName, id.CircuitName, id.PeeringName)
	connResp, err := circuitConnClient.List(ctx, circuitConnectionId)
	if err != nil && !response.WasNotFound(connResp.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", circuitConnectionId, err)
	}

	var connection *[]expressroutecircuitpeerings.ExpressRouteCircuitConnection
	if model := connResp.Model; model != nil {
		connection = convertCircuitConnectionsToPeeringConnections(model)
	}

	payload.Properties.Connections = connection

	if d.HasChange("shared_key") {
		payload.Properties.SharedKey = pointer.To(d.Get("shared_key").(string))
	}

	if d.HasChange("azure_asn") {
		payload.Properties.AzureASN = pointer.To(int64(d.Get("azure_asn").(int)))
	}

	if d.HasChange("peer_asn") {
		payload.Properties.PeerASN = pointer.To(int64(d.Get("peer_asn").(int)))
	}

	if d.HasChange("vlan_id") {
		payload.Properties.VlanId = pointer.To(int64(d.Get("vlan_id").(int)))
	}

	if d.HasChange("gateway_manager_etag") {
		payload.Properties.GatewayManagerEtag = pointer.To(d.Get("gateway_manager_etag").(string))
	}

	if d.HasChange("ipv4_enabled") {
		ipv4Enabled := d.Get("ipv4_enabled").(bool)
		if ipv4Enabled {
			payload.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateEnabled)
		} else {
			payload.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateDisabled)
		}
	}

	primaryPeerAddressPrefix := d.Get("primary_peer_address_prefix").(string)
	secondaryPeerAddressPrefix := d.Get("secondary_peer_address_prefix").(string)

	routeFilterId := d.Get("route_filter_id").(string)

	if d.HasChange("primary_peer_address_prefix") && !strings.EqualFold(primaryPeerAddressPrefix, "") {
		payload.Properties.PrimaryPeerAddressPrefix = pointer.To(primaryPeerAddressPrefix)
	}

	if d.HasChange("secondary_peer_address_prefix") && !strings.EqualFold(secondaryPeerAddressPrefix, "") {
		payload.Properties.SecondaryPeerAddressPrefix = pointer.To(secondaryPeerAddressPrefix)
	}

	if d.HasChange("microsoft_peering_config") {
		if strings.EqualFold(id.PeeringName, string(expressroutecircuitpeerings.ExpressRoutePeeringTypeMicrosoftPeering)) {
			peerings := d.Get("microsoft_peering_config").([]any)
			if len(peerings) == 0 && primaryPeerAddressPrefix != "" {
				return fmt.Errorf("`microsoft_peering_config` must be specified when config for Ipv4 and `peering_type` is set to `MicrosoftPeering`")
			}

			if len(peerings) != 0 && (primaryPeerAddressPrefix == "" || secondaryPeerAddressPrefix == "") {
				return fmt.Errorf("`primary_peer_address_prefix, secondary_peer_address_prefix` must be specified when config for Ipv4")
			}

			payload.Properties.MicrosoftPeeringConfig = expandExpressRouteCircuitPeeringMicrosoftConfig(peerings)

			if d.HasChange("route_filter_id") && routeFilterId != "" {
				payload.Properties.RouteFilter = &expressroutecircuitpeerings.SubResource{
					Id: pointer.To(routeFilterId),
				}
			}
		} else if routeFilterId != "" {
			return fmt.Errorf("`routeFilterId` may only be specified when `peering_type` is set to `MicrosoftPeering`")
		}
	}

	if d.HasChange("ipv6") {
		ipv6Peering := d.Get("ipv6").([]any)
		if len(ipv6Peering) != 0 && id.PeeringName == string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePublicPeering) {
			return fmt.Errorf("`ipv6` may only be specified when `peering_type` is `MicrosoftPeering` or `AzurePrivatePeering`")
		}

		ipv6PeeringConfig, err := expandExpressRouteCircuitIpv6PeeringConfig(ipv6Peering)
		if err != nil {
			return err
		}
		payload.Properties.IPv6PeeringConfig = ipv6PeeringConfig
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitPeeringRead(d, meta)
}
