// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceExtendedLocations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceExtendedLocationsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.LocationWithoutForceNew(),

			"extended_locations": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceExtendedLocationsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.SubscriptionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSubscriptionID(subscriptionId)
	opts := subscriptions.DefaultListLocationsOperationOptions()
	opts.IncludeExtendedLocations = pointer.To(true)
	resp, err := client.ListLocations(ctx, id, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	normalizedLocation := location.Normalize(d.Get("location").(string))
	d.SetId(fmt.Sprintf("%s/locations/%s", id.ID(), normalizedLocation))

	extendedLocations := getExtendedLocations(resp.Model.Value, normalizedLocation)
	if len(extendedLocations) == 0 {
		return fmt.Errorf("no extended locations were found for the location %q", normalizedLocation)
	}
	if err := d.Set("extended_locations", extendedLocations); err != nil {
		return fmt.Errorf("setting `extended_locations`: %s", err)
	}

	return nil
}

func getExtendedLocations(input *[]subscriptions.Location, normalizedLocation string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.Type == nil || item.Metadata == nil || item.Metadata.HomeLocation == nil || item.Name == nil {
			continue
		}

		if *item.Type != subscriptions.LocationTypeEdgeZone {
			continue
		}

		if location.Normalize(*item.Metadata.HomeLocation) != normalizedLocation {
			continue
		}

		results = append(results, location.Normalize(*item.Name))
	}

	return results
}
