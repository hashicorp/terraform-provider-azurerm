package subscription

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubscriptionsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"display_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name_contains": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"location_placement_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"quota_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"spending_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceArmSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*clients.Client)
	subClient := armClient.Subscription.Client
	tagsClient := armClient.Subscription.TagsClient
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

		if subscriptionId := val.SubscriptionID; subscriptionId != nil {
			subscriptionId := *subscriptionId
			tagsResp, err := tagsClient.GetAtScope(ctx, "subscriptions/"+subscriptionId)
			if err != nil {
				if utils.ResponseWasNotFound(tagsResp.Response) {
					return fmt.Errorf("Error: default tags for Subscription %q was not found", subscriptionId)
				}
				return fmt.Errorf("Error reading default tags for Subscription: %+v", err)
			}
			if tagsResp.Properties == nil {
				return fmt.Errorf("nil tags properties of Subscription %q", subscriptionId)
			}
			s["tags"] = tags.Flatten(tagsResp.Properties.Tags)
		}

		subscriptions = append(subscriptions, s)
	}

	d.SetId("subscriptions-" + armClient.Account.TenantId)
	if err = d.Set("subscriptions", subscriptions); err != nil {
		return fmt.Errorf("Error setting `subscriptions`: %+v", err)
	}

	return nil
}
