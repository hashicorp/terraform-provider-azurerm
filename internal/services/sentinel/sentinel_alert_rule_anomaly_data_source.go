// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AlertRuleAnomalyDataSourceModel struct {
	Name                         string                                  `tfschema:"name"`
	DisplayName                  string                                  `tfschema:"display_name"`
	WorkspaceId                  string                                  `tfschema:"log_analytics_workspace_id"`
	AnomalyVersion               string                                  `tfschema:"anomaly_version"`
	AnomalySettingsVersion       int64                                   `tfschema:"anomaly_settings_version"`
	Description                  string                                  `tfschema:"description"`
	Enabled                      bool                                    `tfschema:"enabled"`
	Frequency                    string                                  `tfschema:"frequency"`
	RequiredDataConnectors       []AnomalyRuleRequiredDataConnectorModel `tfschema:"required_data_connector"`
	SettingsDefinitionId         string                                  `tfschema:"settings_definition_id"`
	Mode                         string                                  `tfschema:"mode"`
	Tactics                      []string                                `tfschema:"tactics"`
	Techniques                   []string                                `tfschema:"techniques"`
	ThresholdObservation         []AnomalyRuleThresholdModel             `tfschema:"threshold_observation"`
	MultiSelectObservation       []AnomalyRuleMultiSelectModel           `tfschema:"multi_select_observation"`
	SingleSelectObservation      []AnomalyRuleSingleSelectModel          `tfschema:"single_select_observation"`
	PrioritizeExcludeObservation []AnomalyRulePriorityModel              `tfschema:"prioritized_exclude_observation"`
}

type AlertRuleAnomalyDataSource struct{}

var _ sdk.DataSource = AlertRuleAnomalyDataSource{}

func (a AlertRuleAnomalyDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},
	}
}

func (a AlertRuleAnomalyDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"anomaly_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"anomaly_settings_version": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"frequency": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"required_data_connector": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"connector_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"data_types": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"settings_definition_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tactics": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"techniques": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"multi_select_observation": AnomalyRuleMultiSelectSchema(),

		"single_select_observation": AnomalyRuleSingleSelectSchema(),

		"prioritized_exclude_observation": AnomalyRulePrioritySchema(),

		"threshold_observation": AnomalyRuleThresholdSchema(),
	}
}

func (a AlertRuleAnomalyDataSource) ModelObject() interface{} {
	return &AlertRuleAnomalyDataSourceModel{}
}

func (a AlertRuleAnomalyDataSource) ResourceType() string {
	return "azurerm_sentinel_alert_rule_anomaly"
}

func (a AlertRuleAnomalyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleAnomalyDataSourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			workspaceId, err := workspaces.ParseWorkspaceID(config.WorkspaceId)
			if err != nil {
				return err
			}

			items, err := client.ListComplete(ctx, securitymlanalyticssettings.WorkspaceId(*workspaceId))
			if err != nil {
				return fmt.Errorf("listing alerts rules on %s: %+v", workspaceId, err)
			}

			var rule *securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings
			var ruleName string
			for _, item := range items.Items {
				v, ok := item.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings)
				if !ok {
					continue
				}

				switch {
				case config.Name != "":
					if strings.EqualFold(pointer.From(v.Name), config.Name) {
						ruleName = config.Name
						rule = &v
						break
					}
				case config.DisplayName != "":
					if v.Properties != nil && strings.EqualFold(v.Properties.DisplayName, config.DisplayName) {
						ruleName = config.DisplayName
						rule = &v
						break
					}
				}
			}

			if rule == nil {
				return fmt.Errorf("retrieving anomaly rule (%s): not found", ruleName)
			}

			id := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, pointer.From(rule.Name))
			metadata.SetID(id)

			state := AlertRuleAnomalyDataSourceModel{
				WorkspaceId: workspaceId.ID(),
			}

			state.Name = pointer.From(rule.Name)
			if props := rule.Properties; props != nil {
				state.Mode = string(props.SettingsStatus)
				state.DisplayName = props.DisplayName
				state.AnomalyVersion = props.AnomalyVersion
				state.AnomalySettingsVersion = pointer.From(props.AnomalySettingsVersion)
				state.Description = pointer.From(props.Description)
				state.Enabled = props.Enabled
				state.Frequency = props.Frequency
				state.RequiredDataConnectors = flattenSentinelAlertRuleAnomalyRequiredDataConnectors(props.RequiredDataConnectors)
				state.SettingsDefinitionId = pointer.From(props.SettingsDefinitionId)
				state.Tactics = pointer.FromEnumSlice(props.Tactics)
				state.Techniques = pointer.From(props.Techniques)

				if co := props.CustomizableObservations; co != nil {
					state.MultiSelectObservation = flattenSentinelAlertRuleAnomalyMultiSelect(co.MultiSelectObservations)
					state.SingleSelectObservation = flattenSentinelAlertRuleAnomalySingleSelect(co.SingleSelectObservations)
					state.PrioritizeExcludeObservation = flattenSentinelAlertRuleAnomalyPriority(co.PrioritizeExcludeObservations)
					state.ThresholdObservation = flattenSentinelAlertRuleAnomalyThreshold(co.ThresholdObservations)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
