package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceLocalNetworkGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceLocalNetworkGatewayRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": location.SchemaComputed(),

			"gateway_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"gateway_fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"address_space": {
				Type:     pluginsdk.TypeList,
				Computed: true,
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
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"bgp_peering_address": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"peer_weight": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceLocalNetworkGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LocalNetworkGatewaysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLocalNetworkGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading the state of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.LocalNetworkGatewayPropertiesFormat; props != nil {
		d.Set("gateway_address", props.GatewayIPAddress)
		d.Set("gateway_fqdn", props.Fqdn)

		if lnas := props.LocalNetworkAddressSpace; lnas != nil {
			d.Set("address_space", lnas.AddressPrefixes)
		}
		flattenedSettings := flattenLocalNetworkGatewayDataSourceBGPSettings(props.BgpSettings)
		if err := d.Set("bgp_settings", flattenedSettings); err != nil {
			return err
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenLocalNetworkGatewayDataSourceBGPSettings(input *network.BgpSettings) []interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	output["asn"] = int(*input.Asn)
	output["bgp_peering_address"] = *input.BgpPeeringAddress
	output["peer_weight"] = int(*input.PeerWeight)

	return []interface{}{output}
}
