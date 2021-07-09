package consumption

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetSubscriptionDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: resourceArmConsumptionBudgetSubscriptionDataSourceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: SchemaConsumptionBudgetSubscriptionResource(),
	}
}

func resourceArmConsumptionBudgetSubscriptionDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := parse.NewConsumptionBudgetSubscriptionID(d.Get("subscription_id").(string), d.Get("name").(string))

 	err := resourceArmConsumptionBudgetRead(d, meta, d.Get("subscription_id").(string), d.Get("name").(string))

 	if err != nil {
 		return fmt.Errorf("error making read request on Azure Consumption Budget %q for scope %q: %+v", d.Get("name").(string), subscriptionId.ID(), err)
 	}

 	d.SetId(parse.NewConsumptionBudgetSubscriptionID(subscriptionId.SubscriptionId, d.Get("name").(string)).ID())

 	// The scope of a Resource Group consumption budget is the Resource Group ID
 	d.Set("resource_group_id", d.Get("resource_group_id").(string))
 	d.Set("subscription_id", d.Get("subscription_id").(string))

 	return nil
}
