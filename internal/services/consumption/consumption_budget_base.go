// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type consumptionBudgetBaseResource struct{}

// Shared nested model structs for consumption budget resources

type ConsumptionBudgetFilterDimensionModel struct {
	Name     string   `tfschema:"name"`
	Operator string   `tfschema:"operator"`
	Values   []string `tfschema:"values"`
}

type ConsumptionBudgetFilterTagModel struct {
	Name     string   `tfschema:"name"`
	Operator string   `tfschema:"operator"`
	Values   []string `tfschema:"values"`
}

type ConsumptionBudgetFilterModel struct {
	Dimension []ConsumptionBudgetFilterDimensionModel `tfschema:"dimension"`
	Tag       []ConsumptionBudgetFilterTagModel       `tfschema:"tag"`
}

type ConsumptionBudgetTimePeriodModel struct {
	StartDate string `tfschema:"start_date"`
	EndDate   string `tfschema:"end_date"`
}

// Full notification model (subscription + resource group variants)
type ConsumptionBudgetNotificationModel struct {
	Enabled       bool     `tfschema:"enabled"`
	Threshold     int64    `tfschema:"threshold"`
	ThresholdType string   `tfschema:"threshold_type"`
	Operator      string   `tfschema:"operator"`
	ContactEmails []string `tfschema:"contact_emails"`
	ContactGroups []string `tfschema:"contact_groups"`
	ContactRoles  []string `tfschema:"contact_roles"`
}

func getDimensionNames() []string {
	return []string{
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
		"SubscriptionID",
		"SubscriptionName",
		"UnitOfMeasure",
	}
}

func (br consumptionBudgetBaseResource) arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},

		"amount": {
			Type:         pluginsdk.TypeFloat,
			Required:     true,
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
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(getDimensionNames(), false),
								},
								"operator": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "In",
									ValidateFunc: validation.StringInSlice([]string{
										"In",
									}, false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									MinItems: 1,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"tag": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"operator": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "In",
									ValidateFunc: validation.StringInSlice([]string{
										"In",
									}, false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},

		"notification": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"threshold": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"threshold_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(budgets.ThresholdTypeActual),
						ValidateFunc: validation.StringInSlice([]string{
							string(budgets.ThresholdTypeActual),
							"Forecasted",
						}, false),
					},
					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(budgets.OperatorTypeEqualTo),
							string(budgets.OperatorTypeGreaterThan),
							string(budgets.OperatorTypeGreaterThanOrEqualTo),
						}, false),
					},

					"contact_emails": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"contact_groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"contact_roles": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"time_grain": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(budgets.TimeGrainTypeMonthly),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(budgets.TimeGrainTypeBillingAnnual),
				string(budgets.TimeGrainTypeBillingMonth),
				string(budgets.TimeGrainTypeBillingQuarter),
				string(budgets.TimeGrainTypeAnnually),
				string(budgets.TimeGrainTypeMonthly),
				string(budgets.TimeGrainTypeQuarterly),
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
	}

	// Consumption Budgets for Management Groups have a different notification schema,
	// here we override the notification schema in the base resource
	for k, v := range fields {
		output[k] = v
	}

	return output
}

func (br consumptionBudgetBaseResource) attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br consumptionBudgetBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient
			id, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (br consumptionBudgetBaseResource) importerFunc() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		if _, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id()); err != nil {
			return err
		}

		return nil
	}
}

func expandConsumptionBudgetTimePeriodFromModel(input []ConsumptionBudgetTimePeriodModel) (*budgets.BudgetTimePeriod, error) {
	if len(input) == 0 {
		return nil, nil
	}

	tp := input[0]
	timePeriod := budgets.BudgetTimePeriod{}

	if _, err := date.ParseTime(time.RFC3339, tp.StartDate); err != nil {
		return nil, fmt.Errorf("start_date '%s' was not in the correct format: %+v", tp.StartDate, err)
	}
	timePeriod.StartDate = tp.StartDate

	if tp.EndDate != "" {
		if _, err := date.ParseTime(time.RFC3339, tp.EndDate); err != nil {
			return nil, fmt.Errorf("end_date '%s' was not in the correct format: %+v", tp.EndDate, err)
		}
		timePeriod.EndDate = pointer.To(tp.EndDate)
	}

	return &timePeriod, nil
}

func flattenConsumptionBudgetTimePeriodToModel(input *budgets.BudgetTimePeriod) []ConsumptionBudgetTimePeriodModel {
	if input == nil {
		return []ConsumptionBudgetTimePeriodModel{}
	}

	return []ConsumptionBudgetTimePeriodModel{
		{
			StartDate: input.StartDate,
			EndDate:   pointer.From(input.EndDate),
		},
	}
}

func expandConsumptionBudgetNotificationsFromModel(input []ConsumptionBudgetNotificationModel) *map[string]budgets.Notification {
	if len(input) == 0 {
		return nil
	}

	notifications := make(map[string]budgets.Notification)
	for _, n := range input {
		notification := budgets.Notification{
			Enabled:  n.Enabled,
			Operator: budgets.OperatorType(n.Operator),
			// nolint: gosec
			Threshold: float64(n.Threshold),
		}

		thresholdType := budgets.ThresholdType(n.ThresholdType)
		notification.ThresholdType = &thresholdType

		notification.ContactEmails = n.ContactEmails
		if len(n.ContactRoles) > 0 {
			notification.ContactRoles = &n.ContactRoles
		}
		if len(n.ContactGroups) > 0 {
			notification.ContactGroups = &n.ContactGroups
		}

		notificationKey := fmt.Sprintf("%s_%s_%f_Percent", string(thresholdType), string(notification.Operator), notification.Threshold)
		notifications[notificationKey] = notification
	}

	return &notifications
}

func flattenConsumptionBudgetNotificationsToModel(input *map[string]budgets.Notification) []ConsumptionBudgetNotificationModel {
	if input == nil {
		return []ConsumptionBudgetNotificationModel{}
	}

	result := make([]ConsumptionBudgetNotificationModel, 0)
	for _, n := range *input {
		thresholdType := string(budgets.ThresholdTypeActual)
		if v := n.ThresholdType; v != nil {
			thresholdType = string(*v)
		}

		model := ConsumptionBudgetNotificationModel{
			Enabled:       n.Enabled,
			Operator:      string(n.Operator),
			Threshold:     int64(n.Threshold),
			ThresholdType: thresholdType,
			ContactEmails: n.ContactEmails,
		}

		if v := n.ContactRoles; v != nil {
			model.ContactRoles = *v
		}
		if v := n.ContactGroups; v != nil {
			model.ContactGroups = *v
		}

		result = append(result, model)
	}

	return result
}

func expandConsumptionBudgetFilterFromModel(input []ConsumptionBudgetFilterModel) *budgets.BudgetFilter {
	if len(input) == 0 {
		return nil
	}

	f := input[0]
	filter := budgets.BudgetFilter{}

	dimensions := make([]budgets.BudgetFilterProperties, 0)
	for _, d := range f.Dimension {
		dimensions = append(dimensions, budgets.BudgetFilterProperties{
			Dimensions: &budgets.BudgetComparisonExpression{
				Name:     d.Name,
				Operator: budgets.BudgetOperatorType(d.Operator),
				Values:   d.Values,
			},
		})
	}

	tags := make([]budgets.BudgetFilterProperties, 0)
	for _, t := range f.Tag {
		tags = append(tags, budgets.BudgetFilterProperties{
			Tags: &budgets.BudgetComparisonExpression{
				Name:     t.Name,
				Operator: budgets.BudgetOperatorType(t.Operator),
				Values:   t.Values,
			},
		})
	}

	dimensionsSet := len(dimensions) > 0
	tagsSet := len(tags) > 0

	if dimensionsSet && tagsSet {
		dimensions = append(dimensions, tags...)
		filter.And = &dimensions
	} else {
		if dimensionsSet {
			if len(dimensions) > 1 {
				filter.And = &dimensions
			} else {
				filter.Dimensions = dimensions[0].Dimensions
			}
		} else if tagsSet {
			if len(tags) > 1 {
				filter.And = &tags
			} else {
				filter.Tags = tags[0].Tags
			}
		}
	}

	return &filter
}

func flattenConsumptionBudgetFilterToModel(input *budgets.BudgetFilter) []ConsumptionBudgetFilterModel {
	if input == nil {
		return []ConsumptionBudgetFilterModel{}
	}

	dimensions := make([]ConsumptionBudgetFilterDimensionModel, 0)
	tags := make([]ConsumptionBudgetFilterTagModel, 0)

	if input.And != nil {
		for _, v := range *input.And {
			if v.Dimensions != nil {
				dimensions = append(dimensions, ConsumptionBudgetFilterDimensionModel{
					Name:     v.Dimensions.Name,
					Operator: string(v.Dimensions.Operator),
					Values:   v.Dimensions.Values,
				})
			}
			if v.Tags != nil {
				tags = append(tags, ConsumptionBudgetFilterTagModel{
					Name:     v.Tags.Name,
					Operator: string(v.Tags.Operator),
					Values:   v.Tags.Values,
				})
			}
		}
	} else {
		if input.Dimensions != nil {
			dimensions = append(dimensions, ConsumptionBudgetFilterDimensionModel{
				Name:     input.Dimensions.Name,
				Operator: string(input.Dimensions.Operator),
				Values:   input.Dimensions.Values,
			})
		}
		if input.Tags != nil {
			tags = append(tags, ConsumptionBudgetFilterTagModel{
				Name:     input.Tags.Name,
				Operator: string(input.Tags.Operator),
				Values:   input.Tags.Values,
			})
		}
	}

	if len(dimensions) == 0 && len(tags) == 0 {
		return []ConsumptionBudgetFilterModel{}
	}

	return []ConsumptionBudgetFilterModel{
		{
			Dimension: dimensions,
			Tag:       tags,
		},
	}
}
