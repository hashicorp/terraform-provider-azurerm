package consumption

import (
	"time"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	resourceParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read:   resourceArmConsumptionBudgetDataSourceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: schemaConsumptionBudgetDataSource(),
	}
}

func schemaConsumptionBudgetDataSource() map[string]*pluginsdk.Schema {
	resourceGroupSubscriptionSchema := map[string]*pluginsdk.Schema{
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}

	return azure.MergeSchema(SchemaConsumptionBudgetCommonResource(), resourceGroupSubscriptionSchema)
}

func resourceArmConsumptionBudgetDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("subscription_id", consumptionBudgetId.SubscriptionId)

	return nil
}