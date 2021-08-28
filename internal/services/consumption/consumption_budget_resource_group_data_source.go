package consumption

import (
	"time"

	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetResourceGroupDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: resourceArmConsumptionBudgetResourceGroupDataSourceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ConsumptionBudgetName(),
			},

			"resource_group_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmConsumptionBudgetResourceGroupDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := resourceParse.NewConsumptionBudgetResourceGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	d.SetId(id.ID())

	err := resourceArmConsumptionBudgetRead(d, meta, id.ID(), d.Get("name").(string))

	if err != nil {
		return err
	}

	return nil
}
