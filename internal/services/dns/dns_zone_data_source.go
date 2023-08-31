// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDnsZone() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDnsZoneRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": {
				// TODO: we need a CommonSchema type for this which doesn't have ForceNew
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

			"name_servers": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDnsZoneRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.Zones
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := zones.NewDnsZoneID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	var zone *zones.Zone
	if id.ResourceGroupName != "" {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%s was not found", id)
			}
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		zone = resp.Model
	} else {
		result, resourceGroupName, err := findZone(ctx, client, id.SubscriptionId, id.DnsZoneName)
		if err != nil {
			return err
		}

		if resourceGroupName == nil {
			return fmt.Errorf("unable to locate the Resource Group for DNS Zone %q in Subscription %q", id.DnsZoneName, subscriptionId)
		}

		zone = result
		id.ResourceGroupName = *resourceGroupName
	}

	if zone == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	d.SetId(id.ID())

	d.Set("name", id.DnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if props := zone.Properties; props != nil {
		d.Set("number_of_record_sets", props.NumberOfRecordSets)
		d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)

		nameServers := make([]string, 0)
		if ns := props.NameServers; ns != nil {
			nameServers = *ns
		}
		if err := d.Set("name_servers", nameServers); err != nil {
			return err
		}
	}

	if err := tags.FlattenAndSet(d, zone.Tags); err != nil {
		return err
	}

	return nil
}

func findZone(ctx context.Context, client *zones.ZonesClient, subscriptionId, name string) (*zones.Zone, *string, error) {
	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	zonesIterator, err := client.ListComplete(ctx, subscriptionResourceId, zones.DefaultListOperationOptions())
	if err != nil {
		return nil, nil, fmt.Errorf("listing DNS Zones: %+v", err)
	}

	var found zones.Zone
	for _, zone := range zonesIterator.Items {
		if zone.Name != nil && *zone.Name == name {
			if found.Id != nil {
				return nil, nil, fmt.Errorf("found multiple DNS zones with name %q, please specify the resource group", name)
			}
			found = zone
		}
	}

	if found.Id == nil {
		return nil, nil, fmt.Errorf("could not find DNS zone with name: %q", name)
	}

	id, err := zones.ParseDnsZoneIDInsensitively(*found.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing %q as a DNS Zone ID: %+v", *found.Id, err)
	}
	return &found, &id.ResourceGroupName, nil
}
