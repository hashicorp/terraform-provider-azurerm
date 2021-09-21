package consumption

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"amount": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},

			"filter": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
							Set:      pluginsdk.HashResource(SchemaConsumptionBudgetFilterDimensionElement()),
							Elem:     SchemaConsumptionBudgetFilterDimensionElement(),
						},
						"tag": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Set:      pluginsdk.HashResource(SchemaConsumptionBudgetFilterTagElement()),
							Elem:     SchemaConsumptionBudgetFilterTagElement(),
						},
						"not": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"dimension": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Computed: true,
										Elem:     SchemaConsumptionBudgetFilterDimensionElement(),
									},
									"tag": {
										Type:     pluginsdk.TypeList,
										Optional: true,
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
				Optional: true,
				Computed: true,
			},

			"time_period": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date": {
							Type:     pluginsdk.TypeString,
							Required: true,
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
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	id := resourceParse.NewConsumptionBudgetResourceGroupID(subscriptionID, d.Get("resource_group_name").(string), d.Get("name").(string))
	d.SetId(id.ID())

	err := resourceArmConsumptionBudgetDataSourceRead(d, meta, id.ID(), d.Get("name").(string))

	if err != nil {
		return err
	}

	return nil
}

func resourceArmConsumptionBudgetDataSourceRead(d *pluginsdk.ResourceData, meta interface{}, scope, name string) error {
	client := meta.(*clients.Client).Consumption.BudgetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making read request on Azure Consumption Budget %q for scope %q: %+v", name, scope, err)
	}

	d.Set("name", resp.Name)
	if resp.Amount != nil {
		amount, _ := resp.Amount.Float64()
		d.Set("amount", amount)
	}
	d.Set("time_grain", string(resp.TimeGrain))
	d.Set("time_period", FlattenConsumptionBudgetTimePeriod(resp.TimePeriod))
	d.Set("notification", pluginsdk.NewSet(pluginsdk.HashResource(SchemaConsumptionBudgetNotificationElement()), FlattenConsumptionBudgetNotifications(resp.Notifications)))
	d.Set("filter", FlattenConsumptionBudgetFilter(resp.Filter))

	return nil
}
