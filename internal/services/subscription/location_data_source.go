// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	resourcesSubscription "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = LocationDataSource{}

type LocationDataSource struct{}

type LocationDataSourceModel struct {
	Location     string                `tfschema:"location"`
	DisplayName  string                `tfschema:"display_name"`
	ZoneMappings []LocationZoneMapping `tfschema:"zone_mappings"`
}

type LocationZoneMapping struct {
	LogicalZone  string `tfschema:"logical_zone"`
	PhysicalZone string `tfschema:"physical_zone"`
}

func (r LocationDataSource) ResourceType() string {
	return "azurerm_location"
}

func (r LocationDataSource) ModelObject() interface{} {
	return &LocationDataSourceModel{}
}

func (r LocationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationWithoutForceNew(),
	}
}

func (r LocationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"zone_mappings": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"logical_zone": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"physical_zone": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r LocationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Subscription.SubscriptionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := commonids.NewSubscriptionID(subscriptionId)
			resp, err := client.ListLocations(ctx, id, resourcesSubscription.DefaultListLocationsOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if resp.Model.Value == nil {
				return fmt.Errorf("retrieving %s: model value was nil", id)
			}

			var model LocationDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			normalizedLocation := location.Normalize(model.Location)

			locationValue, err := getLocation(normalizedLocation, resp.Model.Value)
			if err != nil {
				return err
			}

			var state LocationDataSourceModel
			state.ZoneMappings = flattenZonesMapping(locationValue)
			state.DisplayName = pointer.From(locationValue.DisplayName)
			state.Location = normalizedLocation

			metadata.ResourceData.SetId(fmt.Sprintf("%s/locations/%s", id.ID(), normalizedLocation))

			return metadata.Encode(&state)
		},
	}
}

func getLocation(location string, input *[]resourcesSubscription.Location) (*resourcesSubscription.Location, error) {
	for _, item := range *input {
		if pointer.From(item.Name) == location && strings.EqualFold(string(pointer.From(item.Metadata.RegionType)), "Physical") {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("no location was found for %q", location)
}

func flattenZonesMapping(location *resourcesSubscription.Location) (zoneMappings []LocationZoneMapping) {
	zoneMappings = make([]LocationZoneMapping, 0)

	if location == nil || location.AvailabilityZoneMappings == nil {
		return zoneMappings
	}

	for _, zoneMapping := range *location.AvailabilityZoneMappings {
		locationZoneMapping := LocationZoneMapping{
			LogicalZone:  pointer.From(zoneMapping.LogicalZone),
			PhysicalZone: pointer.From(zoneMapping.PhysicalZone),
		}

		zoneMappings = append(zoneMappings, locationZoneMapping)
	}

	return zoneMappings
}
