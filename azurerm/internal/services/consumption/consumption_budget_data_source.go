package consumption

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/validate"
	resourceParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	subscriptionParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmConsumptionBudgetDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: resourceArmConsumptionBudgetDataSourceRead,
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

func resourceArmConsumptionBudgetDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	resourceGroupID := resourceParse.NewResourceGroupID(d.Get("subscription_id").(string), d.Get("resource_group_id").(string)).ID()
	subscriptionID := subscriptionParse.NewSubscriptionId(d.Get("subscription_id").(string)).ID()
	name := d.Get("name").(string)

	scope := subscriptionID
	if d.Get("resource_group_id").(string) != "" {
		scope = resourceGroupID
	}

	err := resourceArmConsumptionBudgetHelperRead(d, meta, scope, name)
	if err != nil {
		return err
	}

	// The scope of a Resource Group consumption budget is the Resource Group ID
	d.Set("resource_group_id", d.Get("resource_group_id").(string))
	d.Set("subscription_id", d.Get("subscription_id").(string))

	return nil
}

func resourceArmConsumptionBudgetHelperRead(d *pluginsdk.ResourceData, meta interface{}, scope, name string) error {
	client := meta.(*clients.Client).Consumption.BudgetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, scope, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error making read request on Azure Consumption Budget %q for scope %q: %+v", name, scope, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	if resp.Amount != nil {
		amount, _ := resp.Amount.Float64()
		d.Set("amount", amount)
	}
	d.Set("time_grain", string(resp.TimeGrain))
	d.Set("time_period", FlattenConsumptionBudgetTimePeriod(resp.TimePeriod))
	d.Set("notification", pluginsdk.NewSet(pluginsdk.HashResource(SchemaConsumptionBudgetNotificationElement()), FlattenConsumptionBudgetNotifications(resp.Notifications)))
	d.Set("filter", FlattenConsumptionBudgetFilter(resp.Filter))

	return err
}
