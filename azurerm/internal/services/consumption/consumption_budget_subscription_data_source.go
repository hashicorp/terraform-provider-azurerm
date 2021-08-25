package consumption

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	subscriptionParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetSubscriptionDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: resourceArmConsumptionBudgetSubscriptionDataSourceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: SchemaConsumptionBudgetSubscriptionDataSource(),
	}
}

func resourceArmConsumptionBudgetSubscriptionDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	subscriptionID := subscriptionParse.NewSubscriptionId(d.Get("subscription_id").(string))

	err := resourceArmConsumptionBudgetRead(d, meta, subscriptionID.ID(), name)

	if err != nil {
		return fmt.Errorf("error making read request on Azure Consumption Budget %q for scope %q: %+v", d.Get("name").(string), subscriptionID.ID(), err)
	}

	d.SetId(parse.NewConsumptionBudgetSubscriptionID(subscriptionID.SubscriptionID, d.Get("name").(string)).ID())

	// The scope of a Subscription budget resource is the Subscription budget ID
	d.Set("subscription_id", d.Get("subscription_id").(string))

	return nil
}
