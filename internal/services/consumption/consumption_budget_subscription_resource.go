package consumption

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
	subscriptionId := commonids.NewSubscriptionID(d.Get("subscription_id").(string))
	id := parse.NewConsumptionBudgetSubscriptionID(subscriptionId.SubscriptionId, d.Get("name").(string))

	err := resourceArmConsumptionBudgetCreateUpdate(d, meta, consumptionBudgetSubscriptionName, subscriptionId.ID())
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceArmConsumptionBudgetSubscriptionRead(d, meta)
}

func resourceArmConsumptionBudgetSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionId := commonids.NewSubscriptionID(consumptionBudgetId.SubscriptionId)

	err = resourceArmConsumptionBudgetRead(d, meta, subscriptionId.ID(), consumptionBudgetId.BudgetName, SchemaConsumptionBudgetNotificationElement, FlattenConsumptionBudgetNotifications)
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

	subscriptionId := commonids.NewSubscriptionID(consumptionBudgetId.SubscriptionId)

	return resourceArmConsumptionBudgetDelete(d, meta, subscriptionId.ID())
}
