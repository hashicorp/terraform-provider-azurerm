package consumption

import (
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	subscriptionParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmConsumptionBudgetSubscriptionCreateUpdate,
		Read:   resourceArmConsumptionBudgetSubscriptionRead,
		Update: resourceArmConsumptionBudgetSubscriptionCreateUpdate,
		Delete: resourceArmConsumptionBudgetSubscriptionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConsumptionBudgetSubscriptionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: SchemaConsumptionBudgetSubscriptionResource(),
	}
}

func resourceArmConsumptionBudgetSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	subscriptionId := subscriptionParse.NewSubscriptionId(d.Get("subscription_id").(string))

	err := resourceArmConsumptionBudgetCreateUpdate(d, meta, consumptionBudgetSubscriptionName, subscriptionId.ID())
	if err != nil {
		return err
	}

	d.SetId(parse.NewConsumptionBudgetSubscriptionID(subscriptionId.SubscriptionID, name).ID())

	return resourceArmConsumptionBudgetSubscriptionRead(d, meta)
}

func resourceArmConsumptionBudgetSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionId := subscriptionParse.NewSubscriptionId(consumptionBudgetId.SubscriptionId)

	err = resourceArmConsumptionBudgetRead(d, meta, subscriptionId.ID(), consumptionBudgetId.BudgetName)
	if err != nil {
		return err
	}

	d.Set("subscription_id", consumptionBudgetId.SubscriptionId)

	return nil
}

func resourceArmConsumptionBudgetSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionId := subscriptionParse.NewSubscriptionId(consumptionBudgetId.SubscriptionId)

	return resourceArmConsumptionBudgetDelete(d, meta, subscriptionId.ID())
}
