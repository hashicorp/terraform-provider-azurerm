package subscription

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2021-01-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceSubscriptionLocations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSubscriptionLocationsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"locations": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"extended_location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSubscriptionLocationsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := d.Get("subscription_id").(string)

	id := parse.NewLocationID(subscriptionId)

	resp, err := client.ListLocations(ctx, id.SubscriptionId, utils.Bool(true))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := d.Set("locations", flattenSubscriptionLocations(resp.Value)); err != nil {
		return fmt.Errorf("setting `locations`: %s", err)
	}

	return nil
}

func flattenSubscriptionLocations(input *[]subscriptions.Location) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var extendedLocation string
		var location string
		if item.Type == subscriptions.LocationTypeEdgeZone {
			if v := item.Name; v != nil {
				extendedLocation = *v
			}

			if item.Metadata != nil && item.Metadata.HomeLocation != nil {
				location = *item.Metadata.HomeLocation
			}
		} else if v := item.Name; v != nil {
			location = *v
		}

		results = append(results, map[string]interface{}{
			"extended_location": extendedLocation,
			"location":          location,
			"type":              string(item.Type),
		})
	}

	return results
}
