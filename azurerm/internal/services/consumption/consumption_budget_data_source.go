package consumption

import (
	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/validate"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"time"
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
				ForceNew:     true,
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
				Type:         pluginsdk.TypeFloat,
				Computed:     true,
				ValidateFunc: validation.FloatAtLeast(1.0),
			},

			"filter": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:         pluginsdk.TypeSet,
							Computed:     true,
							Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterDimensionElement()),
							Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
							AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
						},
						"tag": {
							Type:         pluginsdk.TypeSet,
							Computed:     true,
							Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterTagElement()),
							Elem:         SchemaConsumptionBudgetFilterTagElement(),
							AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
						},
						"not": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"dimension": {
										Type:         pluginsdk.TypeList,
										MaxItems:     1,
										Computed:     true,
										ExactlyOneOf: []string{"filter.0.not.0.tag"},
										Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
									},
									"tag": {
										Type:         pluginsdk.TypeList,
										MaxItems:     1,
										Computed:     true,
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
				Computed: true,
				MinItems: 1,
				MaxItems: 5,
				Set:      pluginsdk.HashResource(SchemaConsumptionBudgetNotificationElement()),
				Elem:     SchemaConsumptionBudgetNotificationElement(),
			},

			"time_grain": {
				Type:     pluginsdk.TypeString,
				Computed: true,
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
				Computed: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date": {
							Type:         pluginsdk.TypeString,
							Computed:     true,
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

func resourceArmConsumptionBudgetDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	resourceGroupId := d.Get("resource_group_id").(string)
	subscriptionId := d.Get("subscription_id").(string)
	name := d.Get("name").(string)

	if resourceGroupId == "" {
		err := resourceArmConsumptionBudgetRead(d, meta, subscriptionId, name)
		if err != nil {
			return err
		}
	} else {
		err := resourceArmConsumptionBudgetRead(d, meta, resourceGroupId, name)
		if err != nil {
			return err
		}
	}

	// The scope of a Resource Group consumption budget is the Resource Group ID
	d.Set("resource_group_id", resourceGroupId)
	d.Set("subscription_id", subscriptionId)

	return nil
}
