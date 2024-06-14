// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"tags": commonschema.TagsDataSource(),

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
	client := meta.(*clients.Client).Network.VirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("guid", props.ResourceGuid)

			if as := props.AddressSpace; as != nil {
				if err := d.Set("address_space", utils.FlattenStringSlice(as.AddressPrefixes)); err != nil {
					return fmt.Errorf("setting `address_space`: %v", err)
				}
			}

			if options := props.DhcpOptions; options != nil {
				if err := d.Set("dns_servers", utils.FlattenStringSlice(options.DnsServers)); err != nil {
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

		}
		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func flattenVnetSubnetsNames(input *[]virtualnetworks.Subnet) []interface{} {
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

func flattenVnetPeerings(input *[]virtualnetworks.VirtualNetworkPeering) map[string]interface{} {
	output := make(map[string]interface{})

	if peerings := input; peerings != nil {
		for _, vnetPeering := range *peerings {
			if vnetPeering.Name == nil {
				continue
			}

			value := ""
			if props := vnetPeering.Properties; props != nil {
				if props.RemoteVirtualNetwork == nil || props.RemoteVirtualNetwork.Id == nil {
					continue
				}
				value = *props.RemoteVirtualNetwork.Id

			}
			key := *vnetPeering.Name
			output[key] = value
		}
	}

	return output
}

func flattenVnetPeeringsdAddressList(input *[]virtualnetworks.VirtualNetworkPeering) []string {
	var output []string
	if peerings := input; peerings != nil {
		for _, vnetPeering := range *peerings {
			if props := vnetPeering.Properties; props != nil {
				for _, addresses := range *props.RemoteVirtualNetworkAddressSpace.AddressPrefixes {
					if addresses != "" {
						output = append(output, addresses)
					}
				}
			}
		}
	}
	return output
}
