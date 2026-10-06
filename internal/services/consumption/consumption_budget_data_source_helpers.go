// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenConsumptionBudgetTimePeriod(input *budgets.BudgetTimePeriod) []any {
	timePeriod := make([]any, 0)

	if input == nil {
		return timePeriod
	}

	startDate := input.StartDate

	return append(timePeriod, map[string]any{
		"start_date": startDate,
		"end_date":   pointer.From(input.EndDate),
	})
}

func flattenConsumptionBudgetNotifications(input *map[string]budgets.Notification, scope string) []any {
	if input == nil {
		return []any{}
	}

	notifications := make([]any, 0)
	for _, n := range *input {
		block := make(map[string]any)

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

		var emails []any
		if v := n.ContactEmails; v != nil {
			emails = pluginsdk.FlattenSlice(&v)
		}
		block["contact_emails"] = emails

		if scope != "management_group_id" {
			var roles []any
			if v := n.ContactRoles; v != nil {
				roles = pluginsdk.FlattenSlice(v)
			}
			block["contact_roles"] = roles

			var groups []any
			if v := n.ContactGroups; v != nil {
				groups = pluginsdk.FlattenSlice(v)
			}
			block["contact_groups"] = groups
		}

		notifications = append(notifications, block)
	}

	return notifications
}

func flattenConsumptionBudgetComparisonExpression(input *budgets.BudgetComparisonExpression) *map[string]any {
	consumptionBudgetComparisonExpression := make(map[string]any)

	consumptionBudgetComparisonExpression["name"] = input.Name
	consumptionBudgetComparisonExpression["operator"] = input.Operator
	consumptionBudgetComparisonExpression["values"] = pluginsdk.FlattenSlice(&input.Values)

	return &consumptionBudgetComparisonExpression
}

func flattenConsumptionBudgetFilter(input *budgets.BudgetFilter) []any {
	filter := make([]any, 0)

	if input == nil {
		return filter
	}

	dimensions := make([]any, 0)
	tags := make([]any, 0)

	filterBlock := make(map[string]any)

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
