package consumption

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceArmConsumptionBudgetSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmConsumptionBudgetSubscriptionCreateUpdate,
		Read:   resourceArmConsumptionBudgetSubscriptionRead,
		Update: resourceArmConsumptionBudgetSubscriptionCreateUpdate,
		Delete: resourceArmConsumptionBudgetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: SchemaAzureConsumptionBudgetSubscriptionResource(),
	}
}

func resourceArmConsumptionBudgetSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := d.Get("subscription_id").(string)
	scope := fmt.Sprintf("/subscriptions/%s", subscriptionId)

	err := resourceArmConsumptionBudgetCreateUpdate(d, meta, consumptionBudgetSubscriptionName, scope)
	if err != nil {
		return err
	}

	return resourceArmConsumptionBudgetRead(d, meta)
}

func resourceArmConsumptionBudgetSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	err := resourceArmConsumptionBudgetRead(d, meta)
	if err != nil {
		return err
	}

	// Parse the consumption budget to get the scope
	consumptionBudgetId, err := parse.ConsumptionBudgetID(d.Id())
	if err != nil {
		return err
	}

	// Parse the scope to get the subscription ID
	resourceId, err := azure.ParseAzureResourceID(consumptionBudgetId.Scope)
	if err != nil {
		return err
	}

	d.Set("subscription_id", resourceId.SubscriptionID)

	return nil
}
