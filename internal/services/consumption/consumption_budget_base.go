// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type consumptionBudgetBaseResource struct{}

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
			MaxItems: 5,
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

func (br consumptionBudgetBaseResource) createFunc(resourceName, scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			var err error
			scope := metadata.ResourceData.Get(scopeFieldName).(string)

			id := budgets.NewScopedBudgetID(scope, metadata.ResourceData.Get("name").(string))

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			if err = createOrUpdateConsumptionBudget(ctx, client, metadata, id); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (br consumptionBudgetBaseResource) readFunc(scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient
			id, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s, %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.BudgetName)
			// lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			if model := resp.Model; model != nil {
				eTag := ""
				if v := model.ETag; v != nil {
					eTag = *v
				}
				metadata.ResourceData.Set("etag", eTag)

				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("amount", props.Amount)
					metadata.ResourceData.Set("time_grain", string(props.TimeGrain))
					metadata.ResourceData.Set("time_period", flattenConsumptionBudgetTimePeriod(&props.TimePeriod))
					metadata.ResourceData.Set("notification", flattenConsumptionBudgetNotifications(props.Notifications, scopeFieldName))
					metadata.ResourceData.Set("filter", flattenConsumptionBudgetFilter(props.Filter))
				}
			}

			return nil
		},
	}
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

func (br consumptionBudgetBaseResource) updateFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			id, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = createOrUpdateConsumptionBudget(ctx, client, metadata, *id); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (br consumptionBudgetBaseResource) importerFunc() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		_, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		return nil
	}
}

func createOrUpdateConsumptionBudget(ctx context.Context, client *budgets.BudgetsClient, metadata sdk.ResourceMetaData, id budgets.ScopedBudgetId) error {
	timePeriod, err := expandConsumptionBudgetTimePeriod(metadata.ResourceData.Get("time_period").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `time_period`: %+v", err)
	}

	// The Consumption Budget API requires the category type field to be set in a budget's properties.
	// 'Cost' is the only valid Budget type today according to the API spec.

	parameters := budgets.Budget{
		Name: utils.String(id.BudgetName),
		Properties: &budgets.BudgetProperties{
			Amount:        metadata.ResourceData.Get("amount").(float64),
			Category:      budgets.CategoryTypeCost,
			Filter:        expandConsumptionBudgetFilter(metadata.ResourceData.Get("filter").([]interface{})),
			Notifications: expandConsumptionBudgetNotifications(metadata.ResourceData.Get("notification").(*pluginsdk.Set).List()),
			TimeGrain:     budgets.TimeGrainType(metadata.ResourceData.Get("time_grain").(string)),
			TimePeriod:    *timePeriod,
		},
	}

	if v, ok := metadata.ResourceData.GetOk("etag"); ok {
		parameters.ETag = utils.String(v.(string))
	}

	_, err = client.CreateOrUpdate(ctx, id, parameters)
	if err != nil {
		return err
	}

	return nil
}

func expandConsumptionBudgetTimePeriod(i []interface{}) (*budgets.BudgetTimePeriod, error) {
	if len(i) == 0 || i[0] == nil {
		return nil, nil
	}

	input := i[0].(map[string]interface{})
	timePeriod := budgets.BudgetTimePeriod{}

	if startDateInput, ok := input["start_date"].(string); ok {
		_, err := date.ParseTime(time.RFC3339, startDateInput)
		if err != nil {
			return nil, fmt.Errorf("start_date '%s' was not in the correct format: %+v", startDateInput, err)
		}
		timePeriod.StartDate = input["start_date"].(string)
	}

	if endDateInput, ok := input["end_date"].(string); ok {
		if endDateInput != "" {
			_, err := date.ParseTime(time.RFC3339, endDateInput)
			if err != nil {
				return nil, fmt.Errorf("end_date '%s' was not in the correct format: %+v", endDateInput, err)
			}

			timePeriod.EndDate = utils.String(input["end_date"].(string))
		}
	}

	return &timePeriod, nil
}

func flattenConsumptionBudgetTimePeriod(input *budgets.BudgetTimePeriod) []interface{} {
	timePeriod := make([]interface{}, 0)

	if input == nil {
		return timePeriod
	}

	startDate := input.StartDate

	endDate := ""
	if v := input.EndDate; v != nil {
		endDate = *v
	}

	return append(timePeriod, map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	})
}

func expandConsumptionBudgetNotifications(input []interface{}) *map[string]budgets.Notification {
	if len(input) == 0 {
		return nil
	}

	notifications := make(map[string]budgets.Notification)

	for _, v := range input {
		if v != nil {
			notificationRaw := v.(map[string]interface{})
			notification := budgets.Notification{}

			notification.Enabled = notificationRaw["enabled"].(bool)
			notification.Operator = budgets.OperatorType(notificationRaw["operator"].(string))

			notification.Threshold = float64(notificationRaw["threshold"].(int))

			thresholdType := budgets.ThresholdType(notificationRaw["threshold_type"].(string))
			notification.ThresholdType = &thresholdType

			contactEmails := utils.ExpandStringSlice(notificationRaw["contact_emails"].([]interface{}))
			notification.ContactEmails = *contactEmails

			// contact_roles cannot be set on consumption budgets for management groups
			if _, ok := notificationRaw["contact_roles"]; ok {
				notification.ContactRoles = utils.ExpandStringSlice(notificationRaw["contact_roles"].([]interface{}))
			}

			// contact_groups cannot be set on consumption budgets for management groups
			if _, ok := notificationRaw["contact_groups"]; ok {
				notification.ContactGroups = utils.ExpandStringSlice(notificationRaw["contact_groups"].([]interface{}))
			}

			notificationKey := fmt.Sprintf("%s_%s_%f_Percent", string(thresholdType), string(notification.Operator), notification.Threshold)
			notifications[notificationKey] = notification
		}
	}

	return &notifications
}

func flattenConsumptionBudgetNotifications(input *map[string]budgets.Notification, scope string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	notifications := make([]interface{}, 0)
	for _, n := range *input {
		block := make(map[string]interface{})

		block["enabled"] = n.Enabled

		operator := ""
		if v := n.Operator; v != "" {
			operator = string(v)
		}
		block["operator"] = operator

		block["threshold"] = n.Threshold

		thresholdType := string(budgets.ThresholdTypeActual)
		if v := n.ThresholdType; v != nil {
			thresholdType = string(*v)
		}
		block["threshold_type"] = thresholdType

		var emails []interface{}
		if v := n.ContactEmails; v != nil {
			emails = utils.FlattenStringSlice(&v)
		}
		block["contact_emails"] = emails

		if scope != "management_group_id" {
			var roles []interface{}
			if v := n.ContactRoles; v != nil {
				roles = utils.FlattenStringSlice(v)
			}
			block["contact_roles"] = roles

			var groups []interface{}
			if v := n.ContactGroups; v != nil {
				groups = utils.FlattenStringSlice(v)
			}
			block["contact_groups"] = groups
		}

		notifications = append(notifications, block)
	}

	return notifications
}

func expandConsumptionBudgetComparisonExpression(input interface{}) *budgets.BudgetComparisonExpression {
	if input == nil {
		return nil
	}

	v := input.(map[string]interface{})

	return &budgets.BudgetComparisonExpression{
		Name:     v["name"].(string),
		Operator: budgets.BudgetOperatorType(v["operator"].(string)),
		Values:   *utils.ExpandStringSlice(v["values"].([]interface{})),
	}
}

func flattenConsumptionBudgetComparisonExpression(input *budgets.BudgetComparisonExpression) *map[string]interface{} {
	consumptionBudgetComparisonExpression := make(map[string]interface{})

	consumptionBudgetComparisonExpression["name"] = input.Name
	consumptionBudgetComparisonExpression["operator"] = input.Operator
	consumptionBudgetComparisonExpression["values"] = utils.FlattenStringSlice(&input.Values)

	return &consumptionBudgetComparisonExpression
}

func expandConsumptionBudgetFilterDimensions(input []interface{}) []budgets.BudgetFilterProperties {
	if len(input) == 0 {
		return nil
	}

	dimensions := make([]budgets.BudgetFilterProperties, 0)

	for _, v := range input {
		dimension := budgets.BudgetFilterProperties{
			Dimensions: expandConsumptionBudgetComparisonExpression(v),
		}
		dimensions = append(dimensions, dimension)
	}

	return dimensions
}

func expandConsumptionBudgetFilterTag(input []interface{}) []budgets.BudgetFilterProperties {
	if len(input) == 0 {
		return nil
	}

	tags := make([]budgets.BudgetFilterProperties, 0)

	for _, v := range input {
		tag := budgets.BudgetFilterProperties{
			Tags: expandConsumptionBudgetComparisonExpression(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

func expandConsumptionBudgetFilter(i []interface{}) *budgets.BudgetFilter {
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	filter := budgets.BudgetFilter{}

	tags := expandConsumptionBudgetFilterTag(input["tag"].(*pluginsdk.Set).List())
	dimensions := expandConsumptionBudgetFilterDimensions(input["dimension"].(*pluginsdk.Set).List())

	tagsSet := len(tags) > 0
	dimensionsSet := len(dimensions) > 0

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

func flattenConsumptionBudgetFilter(input *budgets.BudgetFilter) []interface{} {
	filter := make([]interface{}, 0)

	if input == nil {
		return filter
	}

	dimensions := make([]interface{}, 0)
	tags := make([]interface{}, 0)

	filterBlock := make(map[string]interface{})

	if input.And != nil {
		for _, v := range *input.And {
			if v.Dimensions != nil {
				dimensions = append(dimensions, flattenConsumptionBudgetComparisonExpression(v.Dimensions))
			} else {
				tags = append(tags, flattenConsumptionBudgetComparisonExpression(v.Tags))
			}
		}

		if len(dimensions) != 0 {
			filterBlock["dimension"] = dimensions
		}

		if len(tags) != 0 {
			filterBlock["tag"] = tags
		}
	} else {
		if input.Tags != nil {
			filterBlock["tag"] = append(tags, flattenConsumptionBudgetComparisonExpression(input.Tags))
		}

		if input.Dimensions != nil {
			filterBlock["dimension"] = append(dimensions, flattenConsumptionBudgetComparisonExpression(input.Dimensions))
		}
	}

	if len(filterBlock) != 0 {
		filter = append(filter, filterBlock)
	}

	return filter
}
