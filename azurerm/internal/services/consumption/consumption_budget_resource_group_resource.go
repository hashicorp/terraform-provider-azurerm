package consumption

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceArmConsumptionBudgetResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmConsumptionBudgetResourceGroupCreateUpdate,
		Read:   resourceArmConsumptionBudgetResourceGroupRead,
		Update: resourceArmConsumptionBudgetResourceGroupCreateUpdate,
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

		Schema: SchemaAzureConsumptionBudgetResourceGroupResource(),
	}
}

func resourceArmConsumptionBudgetResourceGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := d.Get("subscription_id").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, resourceGroup)

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

	// Parse the scope to get the subscription ID
	resourceId, err := azure.ParseAzureResourceID(consumptionBudgetId.Scope)
	if err != nil {
		return err
	}

	d.Set("subscription_id", resourceId.SubscriptionID)
	d.Set("resource_group_name", resourceId.ResourceGroup)

	return nil
}
