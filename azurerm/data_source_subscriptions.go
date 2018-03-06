package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/subscription"
)

func dataSourceArmSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubscriptionsRead,

		Schema: map[string]*schema.Schema{
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: subscription.SubscriptionSchema(false),
				},
			},
		},
	}
}

func dataSourceArmSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	subClient := armClient.subscriptionsClient
	ctx := armClient.StopContext

	//ListComplete returns an iterator struct
	results, err := subClient.ListComplete(ctx)
	if err != nil {
		return fmt.Errorf("Error listing subscriptions: %+v", err)
	}

	//iterate across each subscriptions and append them to slice
	subscriptions := make([]map[string]interface{}, 0)
	for err = nil; results.NotDone(); err = results.Next() {
		val := results.Value()

		s := make(map[string]interface{})

		if v := val.SubscriptionID; v != nil {
			s["subscription_id"] = *v
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

		subscriptions = append(subscriptions, s)
	}

	d.SetId("subscriptions-" + armClient.tenantId)
	if err := d.Set("subscriptions", subscriptions); err != nil {
		return fmt.Errorf("Error flattening `subscriptions`: %+v", err)
	}

	return nil
}
