package consumption

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
				Type:         pluginsdk.TypeFloat,
				Computed:     true,
				ValidateFunc: validation.FloatAtLeast(1.0),
			},

			"filter": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:         pluginsdk.TypeSet,
							Optional:     true,
							Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterDimensionElement()),
							Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
							AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
						},
						"tag": {
							Type:         pluginsdk.TypeSet,
							Optional:     true,
							Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterTagElement()),
							Elem:         SchemaConsumptionBudgetFilterTagElement(),
							AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
						},
						"not": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"dimension": {
										Type:         pluginsdk.TypeList,
										MaxItems:     1,
										Optional:     true,
										ExactlyOneOf: []string{"filter.0.not.0.tag"},
										Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
									},
									"tag": {
										Type:         pluginsdk.TypeList,
										MaxItems:     1,
										Optional:     true,
										ExactlyOneOf: []string{"filter.0.not.0.dimension"},
										Elem:         SchemaConsumptionBudgetFilterTagElement(),
									},
								},
							},
							AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
						},
					},
				},
			},

			"notification": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Set:      pluginsdk.HashResource(SchemaConsumptionBudgetNotificationElement()),
				Elem:     SchemaConsumptionBudgetNotificationElement(),
			},

			"time_grain": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(consumption.TimeGrainTypeMonthly),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(consumption.TimeGrainTypeBillingAnnual),
					string(consumption.TimeGrainTypeBillingMonth),
					string(consumption.TimeGrainTypeBillingQuarter),
					string(consumption.TimeGrainTypeAnnually),
					string(consumption.TimeGrainTypeMonthly),
					string(consumption.TimeGrainTypeQuarterly),
				}, false),
			},

			"time_period": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ConsumptionBudgetTimePeriodStartDate,
							ForceNew:     true,
						},
						"end_date": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsRFC3339Time,
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
