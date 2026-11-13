// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_log

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
)

func flattenAzureRmScheduledQueryRulesLogCriteria(input []scheduledqueryrules.Criteria) []any {
	result := make([]any, 0)
	for _, criteria := range input {
		v := make(map[string]any)

		v["dimension"] = flattenAzureRmScheduledQueryRulesLogDimension(criteria.Dimensions)
		v["metric_name"] = criteria.MetricName

		result = append(result, v)
	}
	return result
}

func flattenAzureRmScheduledQueryRulesLogDimension(input *[]scheduledqueryrules.Dimension) []any {
	result := make([]any, 0)

	if input != nil {
		for _, dimension := range *input {
			v := make(map[string]any)

			v["name"] = dimension.Name
			v["operator"] = dimension.Operator
			v["values"] = dimension.Values
			result = append(result, v)
		}
	}
	return result
}
