// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
					string(network.ExpressRoutePeeringTypeAzurePrivatePeering),
					string(network.ExpressRoutePeeringTypeAzurePublicPeering),
					string(network.ExpressRoutePeeringTypeMicrosoftPeering),
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
	client := meta.(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("peering_type", id.PeeringName)
	d.Set("express_route_circuit_name", id.ExpressRouteCircuitName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.ExpressRouteCircuitPeeringPropertiesFormat; props != nil {
		d.Set("azure_asn", props.AzureASN)
		d.Set("peer_asn", props.PeerASN)
		d.Set("primary_azure_port", props.PrimaryAzurePort)
		d.Set("secondary_azure_port", props.SecondaryAzurePort)
		d.Set("primary_peer_address_prefix", props.PrimaryPeerAddressPrefix)
		d.Set("secondary_peer_address_prefix", props.SecondaryPeerAddressPrefix)
		d.Set("vlan_id", props.VlanID)
		d.Set("gateway_manager_etag", props.GatewayManagerEtag)
		d.Set("ipv4_enabled", props.State == network.ExpressRoutePeeringStateEnabled)

		routeFilterId := ""
		if props.RouteFilter != nil && props.RouteFilter.ID != nil {
			routeFilterId = *props.RouteFilter.ID
		}
		d.Set("route_filter_id", routeFilterId)
	}

	return nil
}
