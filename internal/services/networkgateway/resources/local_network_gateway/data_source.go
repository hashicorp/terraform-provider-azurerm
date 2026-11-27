// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package local_network_gateway

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Local Network Gateway.",
		Read:        dataSourceLocalNetworkGatewayRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the Local Network Gateway.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"gateway_address": {
				Description: "The gateway IP address the Local Network Gateway uses.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"gateway_fqdn": {
				Description: "The gateway FQDN the Local Network Gateway uses.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"address_space": {
				Description: "The list of string CIDRs representing the address spaces the gateway exposes.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"bgp_settings": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"asn": {
							Description: "The BGP speaker's ASN.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"bgp_peering_address": {
							Description: "The BGP peering address and BGP identifier of this BGP speaker.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"peer_weight": {
							Description: "The weight added to routes learned from this BGP speaker.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceLocalNetworkGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.LocalNetworkGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := localnetworkgateways.NewLocalNetworkGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.LocalNetworkGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties
		d.Set("gateway_address", props.GatewayIPAddress)
		d.Set("gateway_fqdn", props.Fqdn)

		if lnas := props.LocalNetworkAddressSpace; lnas != nil {
			d.Set("address_space", lnas.AddressPrefixes)
		}
		if err := d.Set("bgp_settings", flattenLocalNetworkGatewayDataSourceBGPSettings(props.BgpSettings)); err != nil {
			return err
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func flattenLocalNetworkGatewayDataSourceBGPSettings(input *localnetworkgateways.BgpSettings) []any {
	output := make(map[string]any)

	if input == nil {
		return []any{}
	}

	output["asn"] = int(*input.Asn)
	output["bgp_peering_address"] = *input.BgpPeeringAddress
	output["peer_weight"] = int(*input.PeerWeight)

	return []any{output}
}
