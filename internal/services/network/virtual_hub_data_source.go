// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVirtualHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualHubRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.VirtualHubName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"address_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"virtual_wan_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"virtual_router_asn": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"virtual_router_ips": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.TagsDataSource(),

			"default_route_table_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVirtualHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVirtualHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.VirtualHubsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VirtualHubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	defaultRouteTable := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, "defaultRouteTable")
	d.Set("default_route_table_id", defaultRouteTable.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("address_prefix", props.AddressPrefix)

			var virtualWanId *string
			if props.VirtualWAN != nil {
				virtualWanId = props.VirtualWAN.Id
			}
			d.Set("virtual_wan_id", virtualWanId)

			var virtualRouterAsn *int64
			if props.VirtualRouterAsn != nil {
				virtualRouterAsn = props.VirtualRouterAsn
			}
			d.Set("virtual_router_asn", virtualRouterAsn)

			var virtualRouterIps *[]string
			if props.VirtualRouterIPs != nil {
				virtualRouterIps = props.VirtualRouterIPs
			}
			d.Set("virtual_router_ips", virtualRouterIps)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
