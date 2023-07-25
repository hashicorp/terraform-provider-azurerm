// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func dataSourceVirtualNetwork() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVnetRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"tags": tags.SchemaDataSource(),

			"address_space": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"dns_servers": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subnets": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"vnet_peerings": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"vnet_peerings_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceVnetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroupName, id.VirtualNetworkName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: %s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.VirtualNetworkPropertiesFormat; props != nil {
		d.Set("guid", props.ResourceGUID)

		if as := props.AddressSpace; as != nil {
			if err := d.Set("address_space", utils.FlattenStringSlice(as.AddressPrefixes)); err != nil {
				return fmt.Errorf("setting `address_space`: %v", err)
			}
		}

		if options := props.DhcpOptions; options != nil {
			if err := d.Set("dns_servers", utils.FlattenStringSlice(options.DNSServers)); err != nil {
				return fmt.Errorf("setting `dns_servers`: %v", err)
			}
		}

		if err := d.Set("subnets", flattenVnetSubnetsNames(props.Subnets)); err != nil {
			return fmt.Errorf("setting `subnets`: %v", err)
		}

		if err := d.Set("vnet_peerings", flattenVnetPeerings(props.VirtualNetworkPeerings)); err != nil {
			return fmt.Errorf("setting `vnet_peerings`: %v", err)
		}

		if err := d.Set("vnet_peerings_addresses", flattenVnetPeeringsdAddressList(props.VirtualNetworkPeerings)); err != nil {
			return fmt.Errorf("setting `vnet_peerings_addresses`: %v", err)
		}

		return tags.FlattenAndSet(d, resp.Tags)
	}

	return nil
}

func flattenVnetSubnetsNames(input *[]network.Subnet) []interface{} {
	subnets := make([]interface{}, 0)

	if mysubnets := input; mysubnets != nil {
		for _, subnet := range *mysubnets {
			if v := subnet.Name; v != nil {
				subnets = append(subnets, *v)
			}
		}
	}
	return subnets
}

func flattenVnetPeerings(input *[]network.VirtualNetworkPeering) map[string]interface{} {
	output := make(map[string]interface{})

	if peerings := input; peerings != nil {
		for _, vnetpeering := range *peerings {
			if vnetpeering.Name == nil || vnetpeering.RemoteVirtualNetwork == nil || vnetpeering.RemoteVirtualNetwork.ID == nil {
				continue
			}

			key := *vnetpeering.Name
			value := *vnetpeering.RemoteVirtualNetwork.ID

			output[key] = value
		}
	}

	return output
}

func flattenVnetPeeringsdAddressList(input *[]network.VirtualNetworkPeering) []string {
	var output []string
	if peerings := input; peerings != nil {
		for _, vnetpeering := range *peerings {
			for _, addresses := range *vnetpeering.RemoteVirtualNetworkAddressSpace.AddressPrefixes {
				if addresses != "" {
					output = append(output, addresses)
				}
			}
		}
	}
	return output
}
