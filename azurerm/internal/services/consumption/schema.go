package consumption

import (
	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-01-01/consumption"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/validate"
)

func SchemaAzureConsumptionBudgetResourceGroupResource() map[string]*schema.Schema {
	resourceGroupNameSchema := map[string]*schema.Schema{
		"resource_group_name": azure.SchemaResourceGroupName(),
	}

	return azure.MergeSchema(SchemaAzureConsumptionBudgetSubscriptionResource(), resourceGroupNameSchema)
}

func SchemaAzureConsumptionBudgetSubscriptionResource() map[string]*schema.Schema {
	subscriptionIDSchema := map[string]*schema.Schema{
		"subscription_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}

	return azure.MergeSchema(SchemaAzureConsumptionBudgetCommonResource(), subscriptionIDSchema)
}

func SchemaAzureConsumptionBudgetFilterTagElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaAzureConsumptionBudgetNotificationElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"threshold": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(consumption.EqualTo),
					string(consumption.GreaterThan),
					string(consumption.GreaterThanOrEqualTo),
				}, false),
			},

			"contact_emails": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"contact_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"contact_roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaAzureConsumptionBudgetCommonResource() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConsumptionBudgetName(),
		},

		"amount": {
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(1.0),
		},

		"category": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(consumption.Cost),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(consumption.Cost),
				string(consumption.Usage),
			}, false),
		},

		"filter": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"resource_groups": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"resources": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
					"meters": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},
					"tag": {
						Type:     schema.TypeSet,
						Optional: true,
						Set:      schema.HashResource(SchemaAzureConsumptionBudgetFilterTagElement()),
						Elem:     SchemaAzureConsumptionBudgetFilterTagElement(),
					},
				},
			},
		},

		"notification": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 5,
			Set:      schema.HashResource(SchemaAzureConsumptionBudgetNotificationElement()),
			Elem:     SchemaAzureConsumptionBudgetNotificationElement(),
		},

		"time_grain": {
			Type:     schema.TypeString,
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
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_date": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validate.ConsumptionBudgetTimePeriodStartDate,
						ForceNew:     true,
					},
					"end_date": {
						Type:             schema.TypeString,
						Optional:         true,
						Computed:         true,
						ValidateFunc:     validation.IsRFC3339Time,
						DiffSuppressFunc: DiffSuppressFuncConsumptionBudgetTimePeriodEndDate,
					},
				},
			},
		},
	}
}
