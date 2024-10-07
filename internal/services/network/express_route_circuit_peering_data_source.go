// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceExpressRouteCircuitPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceExpressRouteCircuitPeeringRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"peering_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePrivatePeering),
					string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePublicPeering),
					string(expressroutecircuitpeerings.ExpressRoutePeeringTypeMicrosoftPeering),
				}, false),
			},

			"express_route_circuit_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"primary_peer_address_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_peer_address_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ipv4_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"vlan_id": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"shared_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"peer_asn": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"azure_asn": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"primary_azure_port": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_azure_port": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"route_filter_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"gateway_manager_etag": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceExpressRouteCircuitPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := commonids.NewExpressRouteCircuitPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_circuit_name").(string), d.Get("peering_type").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
		}
	}

	return nil
}
