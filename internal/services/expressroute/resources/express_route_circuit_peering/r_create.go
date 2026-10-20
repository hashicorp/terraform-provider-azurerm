// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_peering

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitpeerings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routefilters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitPeeringCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := commonids.NewExpressRouteCircuitPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_circuit_name").(string), d.Get("peering_type").(string))

	locks.ByName(id.CircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.CircuitName, expressRouteCircuitResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_express_route_circuit_peering", id.ID())
		}
	}

	primaryPeerAddressPrefix := d.Get("primary_peer_address_prefix").(string)
	secondaryPeerAddressPrefix := d.Get("secondary_peer_address_prefix").(string)

	routeFilterId := d.Get("route_filter_id").(string)

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

	parameters := expressroutecircuitpeerings.ExpressRouteCircuitPeering{
		Properties: &expressroutecircuitpeerings.ExpressRouteCircuitPeeringPropertiesFormat{
			PeeringType:        pointer.ToEnum[expressroutecircuitpeerings.ExpressRoutePeeringType](id.PeeringName),
			SharedKey:          pointer.To(d.Get("shared_key").(string)),
			AzureASN:           pointer.To(int64(d.Get("azure_asn").(int))),
			PeerASN:            pointer.To(int64(d.Get("peer_asn").(int))),
			VlanId:             pointer.To(int64(d.Get("vlan_id").(int))),
			GatewayManagerEtag: pointer.To(d.Get("gateway_manager_etag").(string)),
			Connections:        connection,
		},
	}

	ipv4Enabled := d.Get("ipv4_enabled").(bool)
	if ipv4Enabled {
		parameters.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateEnabled)
	} else {
		parameters.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateDisabled)
	}

	if !strings.EqualFold(primaryPeerAddressPrefix, "") {
		parameters.Properties.PrimaryPeerAddressPrefix = pointer.To(primaryPeerAddressPrefix)
	}

	if !strings.EqualFold(secondaryPeerAddressPrefix, "") {
		parameters.Properties.SecondaryPeerAddressPrefix = pointer.To(secondaryPeerAddressPrefix)
	}

	if strings.EqualFold(id.PeeringName, string(expressroutecircuitpeerings.ExpressRoutePeeringTypeMicrosoftPeering)) {
		peerings := d.Get("microsoft_peering_config").([]any)
		if len(peerings) == 0 && primaryPeerAddressPrefix != "" {
			return fmt.Errorf("`microsoft_peering_config` must be specified when config for Ipv4 and `peering_type` is set to `MicrosoftPeering`")
		}

		if len(peerings) != 0 && (primaryPeerAddressPrefix == "" || secondaryPeerAddressPrefix == "") {
			return fmt.Errorf("`primary_peer_address_prefix, secondary_peer_address_prefix` must be specified when config for Ipv4")
		}

		parameters.Properties.MicrosoftPeeringConfig = expandExpressRouteCircuitPeeringMicrosoftConfig(peerings)

		if routeFilterId != "" {
			parameters.Properties.RouteFilter = &expressroutecircuitpeerings.SubResource{
				Id: pointer.To(routeFilterId),
			}
		}
	} else if routeFilterId != "" {
		return fmt.Errorf("`routeFilterId` may only be specified when `peering_type` is set to `MicrosoftPeering`")
	}

	ipv6Peering := d.Get("ipv6").([]any)
	if len(ipv6Peering) != 0 && id.PeeringName == string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePublicPeering) {
		return fmt.Errorf("`ipv6` may only be specified when `peering_type` is `MicrosoftPeering` or `AzurePrivatePeering`")
	}

	ipv6PeeringConfig, err := expandExpressRouteCircuitIpv6PeeringConfig(ipv6Peering)
	if err != nil {
		return err
	}
	parameters.Properties.IPv6PeeringConfig = ipv6PeeringConfig

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitPeeringRead(d, meta)
}

func expandExpressRouteCircuitPeeringMicrosoftConfig(input []any) *expressroutecircuitpeerings.ExpressRouteCircuitPeeringConfig {
	if len(input) == 0 {
		return nil
	}
	peering := input[0].(map[string]any)

	prefixes := make([]string, 0)
	for _, v := range peering["advertised_public_prefixes"].([]any) {
		prefixes = append(prefixes, v.(string))
	}

	advertisedCommunities := make([]string, 0)
	for _, v := range peering["advertised_communities"].([]any) {
		advertisedCommunities = append(advertisedCommunities, v.(string))
	}

	return &expressroutecircuitpeerings.ExpressRouteCircuitPeeringConfig{
		AdvertisedPublicPrefixes: pointer.To(prefixes),
		CustomerASN:              pointer.To(int64(peering["customer_asn"].(int))),
		RoutingRegistryName:      pointer.To(peering["routing_registry_name"].(string)),
		AdvertisedCommunities:    pointer.To(advertisedCommunities),
	}
}

func expandExpressRouteCircuitIpv6PeeringConfig(input []any) (*expressroutecircuitpeerings.IPv6ExpressRouteCircuitPeeringConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0].(map[string]any)
	peeringConfig := expressroutecircuitpeerings.IPv6ExpressRouteCircuitPeeringConfig{
		MicrosoftPeeringConfig: expandExpressRouteCircuitPeeringMicrosoftConfig(v["microsoft_peering"].([]any)),
		State:                  pointer.To(expressroutecircuitpeerings.ExpressRouteCircuitPeeringStateEnabled),
	}

	primaryPeerAddressPrefix := v["primary_peer_address_prefix"].(string)
	secondaryPeerAddressPrefix := v["secondary_peer_address_prefix"].(string)
	if !strings.EqualFold(primaryPeerAddressPrefix, "") {
		peeringConfig.PrimaryPeerAddressPrefix = pointer.To(primaryPeerAddressPrefix)
	}
	if !strings.EqualFold(secondaryPeerAddressPrefix, "") {
		peeringConfig.SecondaryPeerAddressPrefix = pointer.To(secondaryPeerAddressPrefix)
	}

	ipv6Enabled := v["enabled"].(bool)
	if !ipv6Enabled {
		peeringConfig.State = pointer.To(expressroutecircuitpeerings.ExpressRouteCircuitPeeringStateDisabled)
	}

	routeFilterId := v["route_filter_id"].(string)
	if routeFilterId != "" {
		if _, err := routefilters.ParseRouteFilterID(routeFilterId); err != nil {
			return nil, err
		}
		peeringConfig.RouteFilter = &expressroutecircuitpeerings.SubResource{
			Id: pointer.To(routeFilterId),
		}
	}
	return &peeringConfig, nil
}
