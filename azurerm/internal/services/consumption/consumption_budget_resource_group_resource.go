package consumption

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	resourceParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmConsumptionBudgetResourceGroupCreateUpdate,
		Read:   resourceArmConsumptionBudgetResourceGroupRead,
		Update: resourceArmConsumptionBudgetResourceGroupCreateUpdate,
		Delete: resourceArmConsumptionBudgetResourceGroupDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConsumptionBudgetResourceGroupID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: SchemaConsumptionBudgetResourceGroupResource(),
	}
}

func resourceArmConsumptionBudgetResourceGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	resourceGroupId, err := resourceParse.ResourceGroupID(d.Get("resource_group_id").(string))
	if err != nil {
		return err
	}

	err = resourceArmConsumptionBudgetCreateUpdate(d, meta, consumptionBudgetResourceGroupName, resourceGroupId.ID())
	if err != nil {
		return err
	}

	d.SetId(parse.NewConsumptionBudgetResourceGroupID(resourceGroupId.SubscriptionId, resourceGroupId.ResourceGroup, name).ID())

	return resourceArmConsumptionBudgetResourceGroupRead(d, meta)
}

func resourceArmConsumptionBudgetResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	resourceGroupId := resourceParse.NewResourceGroupID(consumptionBudgetId.SubscriptionId, consumptionBudgetId.ResourceGroup)

	err = resourceArmConsumptionBudgetRead(d, meta, resourceGroupId.ID(), consumptionBudgetId.BudgetName)
	if err != nil {
		return err
	}

	// The scope of a Resource Group consumption budget is the Resource Group ID
	d.Set("resource_group_id", resourceGroupId.ID())

	return nil
}

func resourceArmConsumptionBudgetResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	resourceGroupId := resourceParse.NewResourceGroupID(consumptionBudgetId.SubscriptionId, consumptionBudgetId.ResourceGroup)

	return resourceArmConsumptionBudgetDelete(d, meta, resourceGroupId.ID())
}
