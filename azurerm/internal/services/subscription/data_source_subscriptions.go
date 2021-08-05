package subscription

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceSubscriptions() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSubscriptionsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"display_name_prefix": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"display_name_contains": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"subscriptions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subscription_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"display_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"state": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"location_placement_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"quota_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"spending_limit": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceSubscriptionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	armClient := meta.(*clients.Client)
	subClient := armClient.Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	displayNamePrefix := strings.ToLower(d.Get("display_name_prefix").(string))
	displayNameContains := strings.ToLower(d.Get("display_name_contains").(string))

	// ListComplete returns an iterator struct
	results, err := subClient.ListComplete(ctx)
	if err != nil {
		return fmt.Errorf("Error listing subscriptions: %+v", err)
	}

	// iterate across each subscriptions and append them to slice
	subscriptions := make([]map[string]interface{}, 0)
	for results.NotDone() {
		val := results.Value()

		s := make(map[string]interface{})

		if v := val.ID; v != nil {
			s["id"] = *v
		}
		if v := val.SubscriptionID; v != nil {
			s["subscription_id"] = *v
		}
		if v := val.TenantID; v != nil {
			s["tenant_id"] = *v
		}
		if v := val.DisplayName; v != nil {
			s["display_name"] = *v
		}
		s["state"] = string(val.State)

		if policies := val.SubscriptionPolicies; policies != nil {
			if v := policies.LocationPlacementID; v != nil {
				s["location_placement_id"] = *v
			}
			if v := policies.QuotaID; v != nil {
				s["quota_id"] = *v
			}
			s["spending_limit"] = string(policies.SpendingLimit)
		}

		if err = results.Next(); err != nil {
			return fmt.Errorf("Error going to next subscriptions value: %+v", err)
		}

		// check if the display name prefix matches the given input
		if displayNamePrefix != "" {
			if !strings.HasPrefix(strings.ToLower(s["display_name"].(string)), displayNamePrefix) {
				// the display name does not match the given prefix
				continue
			}
		}

		// check if the display name matches the 'contains' comparison
		if displayNameContains != "" {
			if !strings.Contains(strings.ToLower(s["display_name"].(string)), displayNameContains) {
				// the display name does not match the contains check
				continue
			}
		}

		s["tags"] = tags.Flatten(val.Tags)

		subscriptions = append(subscriptions, s)
	}

	d.SetId("subscriptions-" + armClient.Account.TenantId)
	if err = d.Set("subscriptions", subscriptions); err != nil {
		return fmt.Errorf("Error setting `subscriptions`: %+v", err)
	}

	return nil
}
