// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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

func flattenConsumptionBudgetComparisonExpression(input *budgets.BudgetComparisonExpression) *map[string]interface{} {
	consumptionBudgetComparisonExpression := make(map[string]interface{})

	consumptionBudgetComparisonExpression["name"] = input.Name
	consumptionBudgetComparisonExpression["operator"] = input.Operator
	consumptionBudgetComparisonExpression["values"] = utils.FlattenStringSlice(&input.Values)

	return &consumptionBudgetComparisonExpression
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
