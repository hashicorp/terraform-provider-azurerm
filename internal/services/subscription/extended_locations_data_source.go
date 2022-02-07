package subscription

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2021-01-01/subscriptions"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceExtendedLocations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceExtendedLocationsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": location.SchemaWithoutForceNew(),

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
	client := meta.(*clients.Client).Subscription.Client
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSubscriptionID(subscriptionId)

	location := d.Get("location").(string)
	includeExtendedLocations := utils.Bool(true)

	resp, err := client.ListLocations(ctx, id.SubscriptionId, includeExtendedLocations)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := d.Set("extended_locations", getExtendedLocations(resp.Value, location)); err != nil {
		return fmt.Errorf("setting `extended_locations`: %s", err)
	}

	return nil
}

func getExtendedLocations(input *[]subscriptions.Location, location string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.Type == subscriptions.LocationTypeEdgeZone && item.Metadata != nil && item.Metadata.HomeLocation != nil && azure.NormalizeLocation(*item.Metadata.HomeLocation) == azure.NormalizeLocation(location) && item.Name != nil {
			extendedLocation := *item.Name
			results = append(results, extendedLocation)
		}
	}

	return results
}
