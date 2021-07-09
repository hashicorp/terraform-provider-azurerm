package consumption

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/parse"
	resourceParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	// subscriptionParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
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
			"resource_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: resourceValidate.ResourceGroupID,
			},
			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"amount": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},

			"filter": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Set:      pluginsdk.HashResource(SchemaConsumptionBudgetFilterDimensionElement()),
							Elem:     SchemaConsumptionBudgetFilterDimensionElement(),
						},
						"tag": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Set:      pluginsdk.HashResource(SchemaConsumptionBudgetFilterTagElement()),
							Elem:     SchemaConsumptionBudgetFilterTagElement(),
						},
						"not": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"dimension": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem:     SchemaConsumptionBudgetFilterDimensionElement(),
									},
									"tag": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem:     SchemaConsumptionBudgetFilterTagElement(),
									},
								},
							},
						},
					},
				},
			},

			"notification": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Set:      pluginsdk.HashResource(SchemaConsumptionBudgetNotificationElement()),
				Elem:     SchemaConsumptionBudgetNotificationElement(),
			},

			"time_grain": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"time_period": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
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