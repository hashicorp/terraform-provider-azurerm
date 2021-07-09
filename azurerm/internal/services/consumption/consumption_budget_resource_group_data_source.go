package consumption

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	resourceParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetResourceGroupDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: resourceArmConsumptionBudgetResourceGroupDataSourceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: SchemaConsumptionBudgetResourceGroupResource(),
	}
}

func resourceArmConsumptionBudgetResourceGroupDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	resourceGroupId := parse.NewConsumptionBudgetResourceGroupID(d.Get("subscription_id").(string), d.Get("resource_group_id").(string), d.Get("name").(string))
	resourceGroupID := resourceParse.NewResourceGroupID(d.Get("subscription_id").(string), d.Get("resource_group_id").(string)).ID()

	err := resourceArmConsumptionBudgetRead(d, meta, resourceGroupID, d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("error making read request on Azure Consumption Budget %q for scope %q: %+v", d.Get("name").(string), resourceGroupID, err)
	}
	
	d.SetId(parse.NewConsumptionBudgetResourceGroupID(resourceGroupId.SubscriptionId, resourceGroupId.ResourceGroup, d.Get("name").(string)).ID())

	// The scope of a Resource Group consumption budget is the Resource Group ID
	d.Set("resource_group_id", d.Get("resource_group_id").(string))
	d.Set("subscription_id", d.Get("subscription_id").(string))
	
	return nil
}