// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePrivateDnsZoneVirtualNetworkLink() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourcePrivateDnsZoneVirtualNetworkLinkRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"private_dns_zone_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"registration_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"resolution_policy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
	if !features.FivePointOh() {
		resource.Schema["private_dns_zone_id"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"private_dns_zone_name"},
		}

		resource.Schema["resource_group_name"] = commonschema.ResourceGroupNameOptional()

		resource.Schema["private_dns_zone_name"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			Deprecated:    "The `private_dns_zone_name` field is deprecated in favor of `private_dns_zone_id`. This will be removed in version 5.0.",
			ConflictsWith: []string{"private_dns_zone_id"},
			RequiredWith:  []string{"resource_group_name"},
		}
	}

	return resource
}

func dataSourcePrivateDnsZoneVirtualNetworkLinkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rawDnsZoneId := d.Get("private_dns_zone_id").(string)
	if !features.FivePointOh() && rawDnsZoneId == "" {
		dnsZoneId := &virtualnetworklinks.PrivateDnsZoneId{
			ResourceGroupName:  d.Get("resource_group_name").(string),
			PrivateDnsZoneName: d.Get("private_dns_zone_name").(string),
			SubscriptionId:     meta.(*clients.Client).Account.SubscriptionId,
		}
		rawDnsZoneId = dnsZoneId.ID()
	}
	dnsZoneId, err := virtualnetworklinks.ParsePrivateDnsZoneID(rawDnsZoneId)
	if err != nil {
		return fmt.Errorf("parsing private DNS zone ID: %+v", err)
	}
	id := virtualnetworklinks.NewVirtualNetworkLinkID(meta.(*clients.Client).Account.SubscriptionId, dnsZoneId.ResourceGroupName, dnsZoneId.PrivateDnsZoneName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VirtualNetworkLinkName)
	if !features.FivePointOh() {
		d.Set("private_dns_zone_name", id.PrivateDnsZoneName)
		d.Set("resource_group_name", id.ResourceGroupName)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("registration_enabled", props.RegistrationEnabled)
			d.Set("resolution_policy", pointer.From(props.ResolutionPolicy))

			if network := props.VirtualNetwork; network != nil {
				d.Set("virtual_network_id", network.Id)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
