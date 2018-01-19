package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmSubscription() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubscriptionRead,
		Schema: map[string]*schema.Schema{

			"subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func dataSourceArmSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	groupClient := client.subscriptionsClient
	ctx := client.StopContext

	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = client.subscriptionId
	}

	resp, err := groupClient.Get(ctx, subscriptionId)
	if err != nil {
		return fmt.Errorf("Error reading subscription: %+v", err)
	}

	d.SetId(*resp.ID)
	d.Set("subscription_id", resp.SubscriptionID)
	d.Set("display_name", resp.DisplayName)
	d.Set("state", resp.State)
	if resp.SubscriptionPolicies != nil {
		d.Set("location_placement_id", resp.SubscriptionPolicies.LocationPlacementID)
		d.Set("quota_id", resp.SubscriptionPolicies.QuotaID)
		d.Set("spending_limit", resp.SubscriptionPolicies.SpendingLimit)
	}

	return nil
}
