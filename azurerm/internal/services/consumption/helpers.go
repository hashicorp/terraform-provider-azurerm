package consumption

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/shopspring/decimal"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// expand and flatten
func ExpandConsumptionBudgetTimePeriod(i []interface{}) (*consumption.BudgetTimePeriod, error) {
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

func FlattenConsumptionBudgetTimePeriod(input *consumption.BudgetTimePeriod) []interface{} {
	timePeriod := make([]interface{}, 0)

	if input == nil {
		return timePeriod
	}

	timePeriodBlock := make(map[string]interface{})

	timePeriodBlock["start_date"] = input.StartDate.String()
	timePeriodBlock["end_date"] = input.EndDate.String()

	return append(timePeriod, timePeriodBlock)
}

func ExpandConsumptionBudgetNotifications(input []interface{}) map[string]*consumption.Notification {
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

			notification.ContactEmails = utils.ExpandStringSlice(notificationRaw["contact_emails"].([]interface{}))
			notification.ContactRoles = utils.ExpandStringSlice(notificationRaw["contact_roles"].([]interface{}))
			notification.ContactGroups = utils.ExpandStringSlice(notificationRaw["contact_groups"].([]interface{}))

			notificationKey := fmt.Sprintf("actual_%s_%s_Percent", string(notification.Operator), notification.Threshold.StringFixed(0))
			notifications[notificationKey] = &notification
		}
	}

	return notifications
}

func FlattenConsumptionBudgetNotifications(input map[string]*consumption.Notification) []interface{} {
	notifications := make([]interface{}, 0)

	if input == nil {
		return notifications
	}

	for _, v := range input {
		if v != nil {
			notificationBlock := make(map[string]interface{})

			notificationBlock["enabled"] = *v.Enabled
			notificationBlock["operator"] = string(v.Operator)
			threshold, _ := v.Threshold.Float64()
			notificationBlock["threshold"] = int(threshold)
			notificationBlock["contact_emails"] = utils.FlattenStringSlice(v.ContactEmails)
			notificationBlock["contact_roles"] = utils.FlattenStringSlice(v.ContactRoles)
			notificationBlock["contact_groups"] = utils.FlattenStringSlice(v.ContactGroups)

			notifications = append(notifications, notificationBlock)
		}
	}

	return notifications
}

func ExpandConsumptionBudgetComparisonExpression(input interface{}) *consumption.BudgetComparisonExpression {
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

func FlattenConsumptionBudgetComparisonExpression(input *consumption.BudgetComparisonExpression) *map[string]interface{} {
	consumptionBudgetComparisonExpression := make(map[string]interface{})

	consumptionBudgetComparisonExpression["name"] = input.Name
	consumptionBudgetComparisonExpression["operator"] = input.Operator
	consumptionBudgetComparisonExpression["values"] = utils.FlattenStringSlice(input.Values)

	return &consumptionBudgetComparisonExpression
}

func ExpandConsumptionBudgetFilterDimensions(input []interface{}) []consumption.BudgetFilterProperties {
	if len(input) == 0 {
		return nil
	}

	dimensions := make([]consumption.BudgetFilterProperties, 0)

	for _, v := range input {
		dimension := consumption.BudgetFilterProperties{
			Dimensions: ExpandConsumptionBudgetComparisonExpression(v),
		}
		dimensions = append(dimensions, dimension)
	}

	return dimensions
}

func ExpandConsumptionBudgetFilterTag(input []interface{}) []consumption.BudgetFilterProperties {
	if len(input) == 0 {
		return nil
	}

	tags := make([]consumption.BudgetFilterProperties, 0)

	for _, v := range input {
		tag := consumption.BudgetFilterProperties{
			Tags: ExpandConsumptionBudgetComparisonExpression(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

func ExpandConsumptionBudgetFilter(i []interface{}) *consumption.BudgetFilter {
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	filter := consumption.BudgetFilter{}

	notBlock := input["not"].([]interface{})
	if len(notBlock) != 0 && notBlock[0] != nil {
		not := notBlock[0].(map[string]interface{})

		tags := ExpandConsumptionBudgetFilterTag(not["tag"].([]interface{}))
		dimensions := ExpandConsumptionBudgetFilterDimensions(not["dimension"].([]interface{}))

		if len(dimensions) != 0 {
			filter.Not = &dimensions[0]
		} else if len(tags) != 0 {
			filter.Not = &tags[0]
		}
	}

	tags := ExpandConsumptionBudgetFilterTag(input["tag"].(*pluginsdk.Set).List())
	dimensions := ExpandConsumptionBudgetFilterDimensions(input["dimension"].(*pluginsdk.Set).List())

	tagsSet := len(tags) > 0
	dimensionsSet := len(dimensions) > 0

	if dimensionsSet && tagsSet {
		and := append(dimensions, tags...)
		filter.And = &and
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

func FlattenConsumptionBudgetFilter(input *consumption.BudgetFilter) []interface{} {
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
			notBlock["dimension"] = []interface{}{FlattenConsumptionBudgetComparisonExpression(input.Not.Dimensions)}
		}

		if input.Not.Tags != nil {
			notBlock["tag"] = []interface{}{FlattenConsumptionBudgetComparisonExpression(input.Not.Tags)}
		}

		if len(notBlock) != 0 {
			filterBlock["not"] = []interface{}{notBlock}
		}
	}

	if input.And != nil {
		for _, v := range *input.And {
			if v.Dimensions != nil {
				dimensions = append(dimensions, FlattenConsumptionBudgetComparisonExpression(v.Dimensions))
			} else {
				tags = append(tags, FlattenConsumptionBudgetComparisonExpression(v.Tags))
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
			filterBlock["tag"] = append(tags, FlattenConsumptionBudgetComparisonExpression(input.Tags))
		}

		if input.Dimensions != nil {
			filterBlock["dimension"] = append(dimensions, FlattenConsumptionBudgetComparisonExpression(input.Dimensions))
		}
	}

	if len(filterBlock) != 0 {
		filter = append(filter, filterBlock)
	}

	return filter
}
