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
					Schema: subscription.SubscriptionSchema(),
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

		s["subscription_id"] = *val.SubscriptionID
		s["display_name"] = *val.DisplayName
		s["state"] = val.State
		if policies := val.SubscriptionPolicies; policies != nil {
			s["location_placement_id"] = *policies.LocationPlacementID
			s["quota_id"] = *policies.QuotaID
			s["spending_limit"] = policies.SpendingLimit
		}
		
		subscriptions = append(subscriptions, s)
	}

	d.SetId("subscriptions-" + armClient.tenantId)
	if err := d.Set("subscriptions", subscriptions); err != nil {
		return fmt.Errorf("Error flattening `subscriptions`: %+v", err)
	}

	return nil
}
