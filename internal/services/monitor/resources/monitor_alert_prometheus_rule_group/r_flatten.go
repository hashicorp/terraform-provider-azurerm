// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_prometheus_rule_group

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroupresources"
)

func flattenPrometheusRuleModel(inputList *[]prometheusrulegroupresources.PrometheusRule) []PrometheusRuleModel {
	outputList := make([]PrometheusRuleModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := PrometheusRuleModel{
			Expression: input.Expression,
		}

		output.Action = flattenPrometheusRuleGroupActionModel(input.Actions)
		output.Alert = pointer.From(input.Alert)
		output.Annotations = pointer.From(input.Annotations)
		output.Enabled = pointer.From(input.Enabled)
		output.For = pointer.From(input.For)
		output.Labels = pointer.From(input.Labels)
		output.Record = pointer.From(input.Record)
		output.AlertResolution = flattenPrometheusRuleAlertResolutionModel(input.ResolveConfiguration)
		output.Severity = pointer.From(input.Severity)
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenPrometheusRuleGroupActionModel(inputList *[]prometheusrulegroupresources.PrometheusRuleGroupAction) []PrometheusRuleGroupActionModel {
	outputList := make([]PrometheusRuleGroupActionModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := PrometheusRuleGroupActionModel{}
		output.ActionGroupId = pointer.From(input.ActionGroupId)
		output.ActionProperties = pointer.From(input.ActionProperties)
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenPrometheusRuleAlertResolutionModel(input *prometheusrulegroupresources.PrometheusRuleResolveConfiguration) []PrometheusRuleAlertResolutionModel {
	outputList := make([]PrometheusRuleAlertResolutionModel, 0)
	if input == nil {
		return outputList
	}

	output := PrometheusRuleAlertResolutionModel{}
	output.AutoResolved = pointer.From(input.AutoResolved)
	output.TimeToResolve = pointer.From(input.TimeToResolve)

	return append(outputList, output)
}
