// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AnomalyRuleRequiredDataConnectorModel struct {
	ConnectorId string   `tfschema:"connector_id"`
	DataTypes   []string `tfschema:"data_types"`
}

type AnomalyRuleMultiSelectModel struct {
	SupportValues []string `tfschema:"supported_values"`
	Values        []string `tfschema:"values"`
	Name          string   `tfschema:"name"`
	Description   string   `tfschema:"description"`
}

func AnomalyRuleMultiSelectSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"description": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"supported_values": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"values": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func flattenSentinelAlertRuleAnomalyMultiSelect(input *[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsMultiSelectObservations) []AnomalyRuleMultiSelectModel {
	output := make([]AnomalyRuleMultiSelectModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		output = append(output, AnomalyRuleMultiSelectModel{
			Name:          pointer.From(item.Name),
			Description:   pointer.From(item.Description),
			SupportValues: pointer.From(item.SupportedValues),
			Values:        pointer.From(item.Values),
		})
	}

	return output
}

type AnomalyRuleSingleSelectModel struct {
	Name          string   `tfschema:"name"`
	Description   string   `tfschema:"description"`
	SupportValues []string `tfschema:"supported_values"`
	Value         string   `tfschema:"value"`
}

func AnomalyRuleSingleSelectSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"description": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"supported_values": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"value": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenSentinelAlertRuleAnomalySingleSelect(input *[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsSingleSelectObservations) []AnomalyRuleSingleSelectModel {
	output := make([]AnomalyRuleSingleSelectModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		output = append(output, AnomalyRuleSingleSelectModel{
			Name:          pointer.From(item.Name),
			Description:   pointer.From(item.Description),
			SupportValues: pointer.From(item.SupportedValues),
			Value:         pointer.From(item.Value),
		})
	}
	return output
}

type AnomalyRulePriorityModel struct {
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	Prioritize  string `tfschema:"prioritize"`
	Exclude     string `tfschema:"exclude"`
}

func AnomalyRulePrioritySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"description": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"prioritize": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"exclude": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenSentinelAlertRuleAnomalyPriority(input *[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations) []AnomalyRulePriorityModel {
	output := make([]AnomalyRulePriorityModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		output = append(output, AnomalyRulePriorityModel{
			Name:        pointer.From(item.Name),
			Description: pointer.From(item.Description),
			Prioritize:  pointer.From(item.Prioritize),
			Exclude:     pointer.From(item.Exclude),
		})
	}
	return output
}

type AnomalyRuleThresholdModel struct {
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	Max         string `tfschema:"max"`
	Min         string `tfschema:"min"`
	Value       string `tfschema:"value"`
}

func AnomalyRuleThresholdSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"description": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"max": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"min": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"value": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenSentinelAlertRuleAnomalyThreshold(input *[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsThresholdObservations) []AnomalyRuleThresholdModel {
	output := make([]AnomalyRuleThresholdModel, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		output = append(output, AnomalyRuleThresholdModel{
			Name:        pointer.From(item.Name),
			Description: pointer.From(item.Description),
			Max:         pointer.From(item.Maximum),
			Min:         pointer.From(item.Minimum),
			Value:       pointer.From(item.Value),
		})
	}
	return output
}

// when the id of workspace is too long, the service return without workspace name:
// "/subscriptions/{sub_id}/resourceGroups/{rg_name}/providers/Microsoft.OperationalInsights/workspaces//providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings/5020e404-9768-4364-98f6-679940c21362",
// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22500
func AlertRuleAnomalyIdFromWorkspaceId(workspaceId workspaces.WorkspaceId, name string) string {
	return securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name).ID()
}

func flattenSentinelAlertRuleAnomalyRequiredDataConnectors(input *[]securitymlanalyticssettings.SecurityMLAnalyticsSettingsDataSource) []AnomalyRuleRequiredDataConnectorModel {
	output := make([]AnomalyRuleRequiredDataConnectorModel, 0)
	if input == nil {
		return output
	}

	for _, v := range *input {
		if v.ConnectorId == nil || v.DataTypes == nil {
			continue
		}

		output = append(output, AnomalyRuleRequiredDataConnectorModel{
			ConnectorId: *v.ConnectorId,
			DataTypes:   *v.DataTypes,
		})
	}

	return output
}
