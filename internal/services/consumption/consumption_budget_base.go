package consumption

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/shopspring/decimal"
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
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
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
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
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
					"not": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dimension": {
									Type:         pluginsdk.TypeList,
									MaxItems:     1,
									Optional:     true,
									ExactlyOneOf: []string{"filter.0.not.0.tag"},
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
									Type:         pluginsdk.TypeList,
									MaxItems:     1,
									Optional:     true,
									ExactlyOneOf: []string{"filter.0.not.0.dimension"},
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
					// Issue: https://github.com/Azure/azure-rest-api-specs/issues/16240
					// Toggling between these two values doesn't work at the moment and also doesn't throw an error
					// but it seems unlikely that a user would switch the threshold_type of their budgets frequently
					"threshold_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(consumption.ThresholdTypeActual),
						ForceNew: true, // TODO: remove this when the above issue is fixed
						ValidateFunc: validation.StringInSlice([]string{
							string(consumption.ThresholdTypeActual),
							"Forecasted",
						}, false),
					},
					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(consumption.OperatorTypeEqualTo),
							string(consumption.OperatorTypeGreaterThan),
							string(consumption.OperatorTypeGreaterThanOrEqualTo),
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

// CLEANUP remove in 3.0
func parseScope(scope string) (string, error) {
	// The validation behaviour of subscription_id isn't correct, it should only accept
	// the resource ID format for the subscription. This ensures backward compatibility
	// for the time being but should be removed in 3.0
	if _, err := uuid.ParseUUID(scope); err == nil {
		return "/subscriptions/" + scope, nil
	}
	if _, err := commonids.ParseSubscriptionID(scope); err == nil {
		return scope, nil
	}
	return "", fmt.Errorf("could not parse %s, was not a valid subscription ID or subscription resource ID", scope)
}

func (br consumptionBudgetBaseResource) createFunc(resourceName, scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			var err error
			scope := metadata.ResourceData.Get(scopeFieldName).(string)
			if scopeFieldName == "subscription_id" && !features.ThreePointOhBeta() {
				scope, err = parseScope(metadata.ResourceData.Get(scopeFieldName).(string))
				if err != nil {
					return err
				}
			}

			id := parse.NewConsumptionBudgetId(scope, metadata.ResourceData.Get("name").(string))

			existing, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !utils.ResponseWasNotFound(existing.Response) {
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
			id, err := parse.ConsumptionBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s, %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.Name)
			//lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			amount := 0.0
			if v := resp.Amount; v != nil {
				amount, _ = v.Float64()
			}
			metadata.ResourceData.Set("amount", amount)

			eTag := ""
			if v := resp.ETag; v != nil {
				eTag = *v
			}
			metadata.ResourceData.Set("etag", eTag)
			metadata.ResourceData.Set("time_grain", string(resp.TimeGrain))
			metadata.ResourceData.Set("time_period", flattenConsumptionBudgetTimePeriod(resp.TimePeriod))
			metadata.ResourceData.Set("notification", flattenConsumptionBudgetNotifications(resp.Notifications, scopeFieldName))
			metadata.ResourceData.Set("filter", flattenConsumptionBudgetFilter(resp.Filter))

			return nil
		},
	}
}

func (br consumptionBudgetBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient
			id, err := parse.ConsumptionBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, id.Scope, id.Name); err != nil {
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

			id, err := parse.ConsumptionBudgetID(metadata.ResourceData.Id())
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

func (br consumptionBudgetBaseResource) importerFunc(expectScope string) sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		var err error
		id, err := parse.ConsumptionBudgetID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		switch expectScope {
		case "subscription":
			_, err = parse.ConsumptionBudgetSubscriptionID(metadata.ResourceData.Id())
		case "resource_group":
			_, err = parse.ConsumptionBudgetResourceGroupID(metadata.ResourceData.Id())
		case "management_group":
			_, err = parse.ConsumptionBudgetManagementGroupID(metadata.ResourceData.Id())
		}

		if err != nil {
			return fmt.Errorf("budget has mismatched scope, expected a budget with %s scope, got %s", expectScope, id.Scope)
		}

		return nil
	}
}

func createOrUpdateConsumptionBudget(ctx context.Context, client *consumption.BudgetsClient, metadata sdk.ResourceMetaData, id parse.ConsumptionBudgetId) error {
	amount := decimal.NewFromFloat(metadata.ResourceData.Get("amount").(float64))
	timePeriod, err := expandConsumptionBudgetTimePeriod(metadata.ResourceData.Get("time_period").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `time_period`: %+v", err)
	}

	// The Consumption Budget API requires the category type field to be set in a budget's properties.
	// 'Cost' is the only valid Budget type today according to the API spec.
	category := "Cost"
	parameters := consumption.Budget{
		Name: utils.String(id.Name),
		BudgetProperties: &consumption.BudgetProperties{
			Amount:        &amount,
			Category:      &category,
			Filter:        expandConsumptionBudgetFilter(metadata.ResourceData.Get("filter").([]interface{})),
			Notifications: expandConsumptionBudgetNotifications(metadata.ResourceData.Get("notification").(*pluginsdk.Set).List()),
			TimeGrain:     consumption.TimeGrainType(metadata.ResourceData.Get("time_grain").(string)),
			TimePeriod:    timePeriod,
		},
	}

	if v, ok := metadata.ResourceData.GetOk("etag"); ok {
		parameters.ETag = utils.String(v.(string))
	}

	_, err = client.CreateOrUpdate(ctx, id.Scope, id.Name, parameters)
	if err != nil {
		return err
	}

	return nil
}

func expandConsumptionBudgetTimePeriod(i []interface{}) (*consumption.BudgetTimePeriod, error) {
	if len(i) == 0 || i[0] == nil {
		return nil, nil
	}

	input := i[0].(map[string]interface{})
	timePeriod := consumption.BudgetTimePeriod{}

	if startDateInput, ok := input["start_date"].(string); ok {
		startDate, err := date.ParseTime(time.RFC3339, startDateInput)
		if err != nil {
			return nil, fmt.Errorf("start_date '%s' was not in the correct format: %+v", startDateInput, err)
		}

		timePeriod.StartDate = &date.Time{
			Time: startDate,
		}
	}

	if endDateInput, ok := input["end_date"].(string); ok {
		if endDateInput != "" {
			endDate, err := date.ParseTime(time.RFC3339, endDateInput)
			if err != nil {
				return nil, fmt.Errorf("end_date '%s' was not in the correct format: %+v", endDateInput, err)
			}

			timePeriod.EndDate = &date.Time{
				Time: endDate,
			}
		}
	}

	return &timePeriod, nil
}

func flattenConsumptionBudgetTimePeriod(input *consumption.BudgetTimePeriod) []interface{} {
	timePeriod := make([]interface{}, 0)

	if input == nil {
		return timePeriod
	}

	startDate := ""
	if v := input.StartDate; v != nil {
		startDate = v.String()
	}

	endDate := ""
	if v := input.EndDate; v != nil {
		endDate = v.String()
	}

	return append(timePeriod, map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	})
}

func expandConsumptionBudgetNotifications(input []interface{}) map[string]*consumption.Notification {
	if len(input) == 0 {
		return nil
	}

	notifications := make(map[string]*consumption.Notification)

	for _, v := range input {
		if v != nil {
			notificationRaw := v.(map[string]interface{})
			notification := consumption.Notification{}

			notification.Enabled = utils.Bool(notificationRaw["enabled"].(bool))
			notification.Operator = consumption.OperatorType(notificationRaw["operator"].(string))

			thresholdDecimal := decimal.NewFromInt(int64(notificationRaw["threshold"].(int)))
			notification.Threshold = &thresholdDecimal

			notification.ThresholdType = consumption.ThresholdType(notificationRaw["threshold_type"].(string))

			notification.ContactEmails = utils.ExpandStringSlice(notificationRaw["contact_emails"].([]interface{}))

			// contact_roles cannot be set on consumption budgets for management groups
			if _, ok := notificationRaw["contact_roles"]; ok {
				notification.ContactRoles = utils.ExpandStringSlice(notificationRaw["contact_roles"].([]interface{}))
			}

			// contact_groups cannot be set on consumption budgets for management groups
			if _, ok := notificationRaw["contact_groups"]; ok {
				notification.ContactGroups = utils.ExpandStringSlice(notificationRaw["contact_groups"].([]interface{}))
			}

			notificationKey := fmt.Sprintf("actual_%s_%s_Percent", string(notification.Operator), notification.Threshold.StringFixed(0))
			notifications[notificationKey] = &notification
		}
	}

	return notifications
}

func flattenConsumptionBudgetNotifications(input map[string]*consumption.Notification, scope string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	notifications := make([]interface{}, 0)
	for _, n := range input {
		if n != nil {
			block := make(map[string]interface{})

			enabled := true
			if v := n.Enabled; v != nil && !*v {
				enabled = false
			}
			block["enabled"] = enabled

			operator := ""
			if v := n.Operator; v != "" {
				operator = string(v)
			}
			block["operator"] = operator

			threshold := 0
			if v := n.Threshold; v != nil {
				t, _ := v.Float64()
				threshold = int(t)
			}
			block["threshold"] = threshold

			thresholdType := string(consumption.ThresholdTypeActual)
			if v := n.ThresholdType; v != consumption.ThresholdTypeActual {
				t := v
				thresholdType = string(t)
			}
			block["threshold_type"] = thresholdType

			var emails []interface{}
			if v := n.ContactEmails; v != nil {
				emails = utils.FlattenStringSlice(v)
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
	}

	return notifications
}

func expandConsumptionBudgetComparisonExpression(input interface{}) *consumption.BudgetComparisonExpression {
	if input == nil {
		return nil
	}

	v := input.(map[string]interface{})

	return &consumption.BudgetComparisonExpression{
		Name:     utils.String(v["name"].(string)),
		Operator: utils.String(v["operator"].(string)),
		Values:   utils.ExpandStringSlice(v["values"].([]interface{})),
	}
}

func flattenConsumptionBudgetComparisonExpression(input *consumption.BudgetComparisonExpression) *map[string]interface{} {
	consumptionBudgetComparisonExpression := make(map[string]interface{})

	consumptionBudgetComparisonExpression["name"] = input.Name
	consumptionBudgetComparisonExpression["operator"] = input.Operator
	consumptionBudgetComparisonExpression["values"] = utils.FlattenStringSlice(input.Values)

	return &consumptionBudgetComparisonExpression
}

func expandConsumptionBudgetFilterDimensions(input []interface{}) []consumption.BudgetFilterProperties {
	if len(input) == 0 {
		return nil
	}

	dimensions := make([]consumption.BudgetFilterProperties, 0)

	for _, v := range input {
		dimension := consumption.BudgetFilterProperties{
			Dimensions: expandConsumptionBudgetComparisonExpression(v),
		}
		dimensions = append(dimensions, dimension)
	}

	return dimensions
}

func expandConsumptionBudgetFilterTag(input []interface{}) []consumption.BudgetFilterProperties {
	if len(input) == 0 {
		return nil
	}

	tags := make([]consumption.BudgetFilterProperties, 0)

	for _, v := range input {
		tag := consumption.BudgetFilterProperties{
			Tags: expandConsumptionBudgetComparisonExpression(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

func expandConsumptionBudgetFilter(i []interface{}) *consumption.BudgetFilter {
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	filter := consumption.BudgetFilter{}

	notBlock := input["not"].([]interface{})
	if len(notBlock) != 0 && notBlock[0] != nil {
		not := notBlock[0].(map[string]interface{})

		tags := expandConsumptionBudgetFilterTag(not["tag"].([]interface{}))
		dimensions := expandConsumptionBudgetFilterDimensions(not["dimension"].([]interface{}))

		if len(dimensions) != 0 {
			filter.Not = &dimensions[0]
		} else if len(tags) != 0 {
			filter.Not = &tags[0]
		}
	}

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

func flattenConsumptionBudgetFilter(input *consumption.BudgetFilter) []interface{} {
	filter := make([]interface{}, 0)

	if input == nil {
		return filter
	}

	dimensions := make([]interface{}, 0)
	tags := make([]interface{}, 0)

	filterBlock := make(map[string]interface{})

	notBlock := make(map[string]interface{})

	if input.Not != nil {
		if input.Not.Dimensions != nil {
			notBlock["dimension"] = []interface{}{flattenConsumptionBudgetComparisonExpression(input.Not.Dimensions)}
		}

		if input.Not.Tags != nil {
			notBlock["tag"] = []interface{}{flattenConsumptionBudgetComparisonExpression(input.Not.Tags)}
		}

		if len(notBlock) != 0 {
			filterBlock["not"] = []interface{}{notBlock}
		}
	}

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
