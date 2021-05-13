package consumption

import (
	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/consumption/validate"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
)

func SchemaConsumptionBudgetResourceGroupResource() map[string]*schema.Schema {
	resourceGroupNameSchema := map[string]*schema.Schema{
		"resource_group_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
	}

	return azure.MergeSchema(SchemaConsumptionBudgetCommonResource(), resourceGroupNameSchema)
}

func SchemaConsumptionBudgetSubscriptionResource() map[string]*schema.Schema {
	subscriptionIDSchema := map[string]*schema.Schema{
		"subscription_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}

	return azure.MergeSchema(SchemaConsumptionBudgetCommonResource(), subscriptionIDSchema)
}

func SchemaConsumptionBudgetFilterDimensionElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ChargeType",
					"Frequency",
					"InvoiceId",
					"Meter",
					"MeterCategory",
					"MeterSubCategory",
					"PartNumber",
					"PricingModel",
					"Product",
					"ProductOrderId",
					"ProductOrderName",
					"PublisherType",
					"ReservationId",
					"ReservationName",
					"ResourceGroupName",
					"ResourceGuid",
					"ResourceId",
					"ResourceLocation",
					"ResourceType",
					"ServiceFamily",
					"ServiceName",
					"UnitOfMeasure",
				}, false),
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "In",
				ValidateFunc: validation.StringInSlice([]string{
					"In",
				}, false),
			},
			"values": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaConsumptionBudgetFilterTagElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "In",
				ValidateFunc: validation.StringInSlice([]string{
					"In",
				}, false),
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

func SchemaConsumptionBudgetNotificationElement() *schema.Resource {
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

func SchemaConsumptionBudgetCommonResource() map[string]*schema.Schema {
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

		"filter": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"dimension": {
						Type:     schema.TypeSet,
						Optional: true,
						Set:      schema.HashResource(SchemaConsumptionBudgetFilterDimensionElement()),
						Elem:     SchemaConsumptionBudgetFilterDimensionElement(),
					},
					"tag": {
						Type:     schema.TypeSet,
						Optional: true,
						Set:      schema.HashResource(SchemaConsumptionBudgetFilterTagElement()),
						Elem:     SchemaConsumptionBudgetFilterTagElement(),
					},
					"not": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dimension": {
									Type:          schema.TypeList,
									MaxItems:      1,
									Optional:      true,
									ConflictsWith: []string{"filter.0.not.0.tag"},
									Elem:          SchemaConsumptionBudgetFilterDimensionElement(),
								},
								"tag": {
									Type:          schema.TypeList,
									MaxItems:      1,
									Optional:      true,
									ConflictsWith: []string{"filter.0.not.0.dimension"},
									Elem:          SchemaConsumptionBudgetFilterTagElement(),
								},
							},
						},
					},
				},
			},
		},

		"notification": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 5,
			Set:      schema.HashResource(SchemaConsumptionBudgetNotificationElement()),
			Elem:     SchemaConsumptionBudgetNotificationElement(),
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
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
	}
}
