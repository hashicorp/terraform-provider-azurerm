// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
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

	var resp *privatezones.PrivateZone
	if id.ResourceGroupName != "" {
		zone, err := client.Get(ctx, id)
		if err != nil || zone.Model == nil {
			if response.WasNotFound(zone.HttpResponse) {
				return fmt.Errorf("%s was not found", id)
			}
			return fmt.Errorf("reading %s: %+v", id, err)
		}
		resp = zone.Model
	} else {
		resourcesClient := meta.(*clients.Client).Resource.ResourcesClient

		zone, err := findPrivateZone(ctx, client, resourcesClient, id.PrivateDnsZoneName)
		if err != nil {
			return err
		}

		if zone == nil {
			return fmt.Errorf("%s was not found", id)
		}

		resp = &zone.zone
		id.ResourceGroupName = zone.resourceGroup
	}

	d.SetId(id.ID())

	d.Set("name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if props := resp.Properties; props != nil {
		d.Set("number_of_record_sets", props.NumberOfRecordSets)
		d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)
		d.Set("max_number_of_virtual_network_links", props.MaxNumberOfVirtualNetworkLinks)
		d.Set("max_number_of_virtual_network_links_with_registration", props.MaxNumberOfVirtualNetworkLinksWithRegistration)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

type privateDnsZone struct {
	zone          privatezones.PrivateZone
	resourceGroup string
}

func findPrivateZone(ctx context.Context, client *privatezones.PrivateZonesClient, resourcesClient *resources.Client, name string) (*privateDnsZone, error) {
	filter := fmt.Sprintf("resourceType eq 'Microsoft.Network/privateDnsZones' and name eq '%s'", name)
	privateZones, err := resourcesClient.List(ctx, filter, "", nil)
	if err != nil {
		return nil, fmt.Errorf("listing Private DNS Zones: %+v", err)
	}

	if len(privateZones.Values()) > 1 {
		return nil, fmt.Errorf("More than one Private DNS Zone found with name: %q", name)
	}

	for _, z := range privateZones.Values() {
		if z.ID == nil {
			continue
		}

		id, err := privatezones.ParsePrivateDnsZoneID(*z.ID)
		if err != nil {
			continue
		}

		zone, err := client.Get(ctx, *id)
		if err != nil || zone.Model == nil {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return &privateDnsZone{
			zone:          *zone.Model,
			resourceGroup: id.ResourceGroupName,
		}, nil
	}

	return nil, fmt.Errorf("No Private DNS Zones found with name: %q", name)
}
