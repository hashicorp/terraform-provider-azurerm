// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
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

func flattenSentinelAlertRuleAnomalyMultiSelect(input *[]azuresdkhacks.AnomalySecurityMLAnalyticsMultiSelectObservations) []AnomalyRuleMultiSelectModel {
	if input == nil {
		return []AnomalyRuleMultiSelectModel{}
	}

	output := make([]AnomalyRuleMultiSelectModel, 0)
	for _, item := range *input {
		o := AnomalyRuleMultiSelectModel{}
		if item.Values != nil {
			values := make([]string, 0)
			values = append(values, *item.Values...)
			o.Values = values
		}
		if item.SupportValues != nil {
			supportValues := make([]string, 0)
			supportValues = append(supportValues, *item.SupportValues...)
			o.SupportValues = supportValues
		}
		if item.Name != nil {
			o.Name = *item.Name
		}
		if item.Description != nil {
			o.Description = *item.Description
		}
		output = append(output, o)
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

func flattenSentinelAlertRuleAnomalySingleSelect(input *[]azuresdkhacks.AnomalySecurityMLAnalyticsSingleSelectObservations) []AnomalyRuleSingleSelectModel {
	if input == nil {
		return []AnomalyRuleSingleSelectModel{}
	}

	output := make([]AnomalyRuleSingleSelectModel, 0)
	for _, item := range *input {
		o := AnomalyRuleSingleSelectModel{}
		if item.Value != nil {
			o.Value = *item.Value
		}
		if item.SupportValues != nil {
			supportValues := make([]string, 0)
			supportValues = append(supportValues, *item.SupportValues...)
			o.SupportValues = supportValues
		}
		if item.Name != nil {
			o.Name = *item.Name
		}
		if item.Description != nil {
			o.Description = *item.Description
		}
		output = append(output, o)
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

func flattenSentinelAlertRuleAnomalyPriority(input *[]azuresdkhacks.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations) []AnomalyRulePriorityModel {
	if input == nil {
		return []AnomalyRulePriorityModel{}
	}

	output := make([]AnomalyRulePriorityModel, 0)
	for _, item := range *input {
		o := AnomalyRulePriorityModel{}
		if item.Prioritize != nil {
			o.Prioritize = *item.Prioritize
		}
		if item.Exclude != nil {
			o.Exclude = *item.Exclude
		}
		if item.Name != nil {
			o.Name = *item.Name
		}
		if item.Description != nil {
			o.Description = *item.Description
		}
		output = append(output, o)
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

func flattenSentinelAlertRuleAnomalyThreshold(input *[]azuresdkhacks.AnomalySecurityMLAnalyticsThresholdObservations) []AnomalyRuleThresholdModel {
	if input == nil {
		return []AnomalyRuleThresholdModel{}
	}

	output := make([]AnomalyRuleThresholdModel, 0)
	for _, item := range *input {
		o := AnomalyRuleThresholdModel{}
		if item.Max != nil {
			o.Max = *item.Max
		}
		if item.Min != nil {
			o.Min = *item.Min
		}
		if item.Value != nil {
			o.Value = *item.Value
		}
		if item.Name != nil {
			o.Name = *item.Name
		}
		if item.Description != nil {
			o.Description = *item.Description
		}
		output = append(output, o)
	}
	return output
}

// the service always return a fixed one no matter what id we pass, tracked on https://github.com/Azure/azure-rest-api-specs/issues/22485
func AlertRuleAnomalyReadWithPredicate(ctx context.Context, baseClient securityinsight.BaseClient, workspaceId workspaces.WorkspaceId, predicateFunc func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool) (*azuresdkhacks.AnomalySecurityMLAnalyticsSettings, error) {
	client := azuresdkhacks.SecurityMLAnalyticsSettingsClient{BaseClient: baseClient}
	resp, err := client.ListComplete(ctx, workspaceId.ResourceGroupName, workspaceId.WorkspaceName)
	if err != nil {
		return nil, fmt.Errorf("retrieving: %+v", err)
	}

	for resp.NotDone() {
		item := resp.Value()
		if v, ok := item.AsAnomalySecurityMLAnalyticsSettings(); ok {
			if predicateFunc(v) {
				return v, nil
			}

		}
		if err := resp.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("listing next: %+v", err)
		}
	}
	return nil, nil
}

// when the id of workspace is too long, the service return without workspace name:
// "/subscriptions/{sub_id}/resourceGroups/{rg_name}/providers/Microsoft.OperationalInsights/workspaces//providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings/5020e404-9768-4364-98f6-679940c21362",
// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22500
func AlertRuleAnomalyIdFromWorkspaceId(workspaceId workspaces.WorkspaceId, name string) string {
	return parse.NewMLAnalyticsSettingsID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name).ID()
}
