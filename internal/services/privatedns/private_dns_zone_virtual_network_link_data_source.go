// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePrivateDnsZoneVirtualNetworkLink() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePrivateDnsZoneVirtualNetworkLinkRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			// TODO: in 4.0 switch this to `private_dns_zone_id`
			"private_dns_zone_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"registration_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourcePrivateDnsZoneVirtualNetworkLinkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := virtualnetworklinks.NewVirtualNetworkLinkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("private_dns_zone_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VirtualNetworkLinkName)
	d.Set("private_dns_zone_name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("registration_enabled", props.RegistrationEnabled)

			if network := props.VirtualNetwork; network != nil {
				d.Set("virtual_network_id", network.Id)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
