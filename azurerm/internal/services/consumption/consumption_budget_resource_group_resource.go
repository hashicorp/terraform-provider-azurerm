package consumption

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmConsumptionBudgetResourceGroupCreateUpdate,
		Read:   resourceArmConsumptionBudgetResourceGroupRead,
		Update: resourceArmConsumptionBudgetResourceGroupCreateUpdate,
		Delete: resourceArmConsumptionBudgetDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConsumptionBudgetID(id)
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
	scope := d.Get("resource_group_id").(string)

	err := resourceArmConsumptionBudgetCreateUpdate(d, meta, consumptionBudgetResourceGroupName, scope)
	if err != nil {
		return err
	}

	return resourceArmConsumptionBudgetRead(d, meta)
}

func resourceArmConsumptionBudgetResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	err := resourceArmConsumptionBudgetRead(d, meta)
	if err != nil {
		return err
	}

	// Parse the consumption budget to get the scope
	consumptionBudgetId, err := parse.ConsumptionBudgetID(d.Id())
	if err != nil {
		return err
	}

	// The scope of a Resource Group consumption budget is the Resource Group ID
	d.Set("resource_group_id", consumptionBudgetId.Scope)

	return nil
}
