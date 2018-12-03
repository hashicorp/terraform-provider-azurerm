package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func dataSourceArmSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubscriptionsRead,

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
					Schema: azure.SchemaSubscription(false),
				},
			},
		},
	}
}

func dataSourceArmSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	subClient := armClient.subscriptionsClient
	ctx := armClient.StopContext

	displayNamePrefix := strings.ToLower(d.Get("display_name_prefix").(string))
	displayNameContains := strings.ToLower(d.Get("display_name_contains").(string))

	//ListComplete returns an iterator struct
	results, err := subClient.ListComplete(ctx)
	if err != nil {
		return fmt.Errorf("Error listing subscriptions: %+v", err)
	}

	//iterate across each subscriptions and append them to slice
	subscriptions := make([]map[string]interface{}, 0)
	for results.NotDone() {
		val := results.Value()

		//check if the display name prefix matches the given input
		if displayNamePrefix != "" {
			if !strings.HasPrefix(strings.ToLower(*val.DisplayName), displayNamePrefix) {
				//the display name does not match the given prefix
				log.Printf("[DEBUG][data_azurerm_subscriptions] %q does not match the prefix check %q", *val.DisplayName, displayNamePrefix)
				if err = results.Next(); err != nil {
					return fmt.Errorf("Error going to next subscriptions value: %+v", err)
				}
				continue
			}
		}

		//check if the display name matches the 'contains' comparison
		if displayNameContains != "" {
			if !strings.Contains(strings.ToLower(*val.DisplayName), displayNameContains) {
				//the display name does not match the contains check
				log.Printf("[DEBUG][data_azurerm_subscriptions] %q does not match the contains check %q", *val.DisplayName, displayNameContains)
				if err = results.Next(); err != nil {
					return fmt.Errorf("Error going to next subscriptions value: %+v", err)
				}
				continue
			}
		}

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

		if err = results.Next(); err != nil {
			return fmt.Errorf("Error going to next subscriptions value: %+v", err)
		}
	}

	d.SetId("subscriptions-" + armClient.tenantId)
	if err = d.Set("subscriptions", subscriptions); err != nil {
		return fmt.Errorf("Error setting `subscriptions`: %+v", err)
	}

	return nil
}
