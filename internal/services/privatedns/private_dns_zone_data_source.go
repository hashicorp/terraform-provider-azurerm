// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePrivateDnsZone() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePrivateDnsZoneRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"number_of_record_sets": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links_with_registration": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourcePrivateDnsZoneRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privatezones.NewPrivateDnsZoneID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if id.ResourceGroupName == "" {
		// we need to discover the Private DNS Zone's resource group
		subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
		zoneId, err := findPrivateDnsZoneId(ctx, client, subscriptionResourceId, id.PrivateDnsZoneName)
		if err != nil {
			return err
		}

		if zoneId == nil {
			return fmt.Errorf("unable to determine the Resource Group for Private DNS Zone %q in Subscription %q", id.PrivateDnsZoneName, id.SubscriptionId)
		}

		id = *zoneId
	}

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("number_of_record_sets", props.NumberOfRecordSets)
			d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)
			d.Set("max_number_of_virtual_network_links", props.MaxNumberOfVirtualNetworkLinks)
			d.Set("max_number_of_virtual_network_links_with_registration", props.MaxNumberOfVirtualNetworkLinksWithRegistration)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func findPrivateDnsZoneId(ctx context.Context, client *privatezones.PrivateZonesClient, subscriptionId commonids.SubscriptionId, name string) (*privatezones.PrivateDnsZoneId, error) {
	opts := privatezones.DefaultListOperationOptions()
	results, err := client.ListComplete(ctx, subscriptionId, opts)
	if err != nil {
		return nil, fmt.Errorf("listing the Private DNS Zones within %s: %+v", subscriptionId, err)
	}

	for _, item := range results.Items {
		if item.Id == nil {
			continue
		}

		itemId := *item.Id
		parsed, err := privatezones.ParsePrivateDnsZoneIDInsensitively(itemId)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a Private DNS Zone ID: %+v", itemId, err)
		}
		if parsed.PrivateDnsZoneName == name {
			return parsed, nil
		}
	}

	return nil, fmt.Errorf("No Private DNS Zones found with name: %q", name)
}
