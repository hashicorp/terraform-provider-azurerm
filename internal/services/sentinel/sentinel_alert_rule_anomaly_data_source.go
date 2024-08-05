// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
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
			var metaModel AlertRuleAnomalyDataSourceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient
			workspaceId, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			setting, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, *workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.Name != nil && strings.EqualFold(*v.Name, metaModel.Name) {
					return true
				}

				if v.DisplayName != nil && strings.EqualFold(*v.DisplayName, metaModel.DisplayName) {
					return true
				}

				return false
			})

			if err != nil {
				return fmt.Errorf("retrieving: %+v", err)
			}
			if setting == nil {
				if metaModel.DisplayName != "" {
					return fmt.Errorf("reading Sentinel Anomaly Rule (Display Name %q) was not found", metaModel.DisplayName)
				}
				return fmt.Errorf("reading Sentinel Anomaly Rule (Name %q) was not found", metaModel.Name)
			}

			id, err := parse.MLAnalyticsSettingsID(AlertRuleAnomalyIdFromWorkspaceId(*workspaceId, *setting.Name))
			if err != nil {
				return fmt.Errorf("parsing: %+v", err)
			}

			state := AlertRuleAnomalyDataSourceModel{
				WorkspaceId: workspaceId.ID(),
				Mode:        string(setting.SettingsStatus),
			}

			if setting.Name != nil {
				state.Name = *setting.Name
			}
			if setting.DisplayName != nil {
				state.DisplayName = *setting.DisplayName
			}
			if setting.AnomalyVersion != nil {
				state.AnomalyVersion = *setting.AnomalyVersion
			}
			if setting.AnomalySettingsVersion != nil {
				state.AnomalySettingsVersion = int64(*setting.AnomalySettingsVersion)
			}
			if setting.Description != nil {
				state.Description = *setting.Description
			}
			if setting.Enabled != nil {
				state.Enabled = *setting.Enabled
			}
			if setting.Frequency != nil {
				state.Frequency = *setting.Frequency
			}
			state.RequiredDataConnectors = flattenSentinelAlertRuleAnomalyRequiredDataConnectors(setting.RequiredDataConnectors)
			if setting.SettingsDefinitionID != nil {
				state.SettingsDefinitionId = setting.SettingsDefinitionID.String()
			}
			state.Tactics = flattenSentinelAlertRuleAnomalyTactics(setting.Tactics)
			if setting.Techniques != nil {
				state.Techniques = *setting.Techniques
			}

			if setting.CustomizableObservations != nil {
				state.MultiSelectObservation = flattenSentinelAlertRuleAnomalyMultiSelect(setting.CustomizableObservations.MultiSelectObservations)
				state.SingleSelectObservation = flattenSentinelAlertRuleAnomalySingleSelect(setting.CustomizableObservations.SingleSelectObservations)
				state.PrioritizeExcludeObservation = flattenSentinelAlertRuleAnomalyPriority(setting.CustomizableObservations.PrioritizeExcludeObservations)
				state.ThresholdObservation = flattenSentinelAlertRuleAnomalyThreshold(setting.CustomizableObservations.ThresholdObservations)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenSentinelAlertRuleAnomalyRequiredDataConnectors(input *[]securityinsight.SecurityMLAnalyticsSettingsDataSource) []AnomalyRuleRequiredDataConnectorModel {
	if input == nil {
		return []AnomalyRuleRequiredDataConnectorModel{}
	}

	output := make([]AnomalyRuleRequiredDataConnectorModel, 0)
	for _, v := range *input {
		if v.ConnectorID == nil || v.DataTypes == nil {
			continue
		}

		output = append(output, AnomalyRuleRequiredDataConnectorModel{
			ConnectorId: *v.ConnectorID,
			DataTypes:   *v.DataTypes,
		})
	}

	return output
}

func flattenSentinelAlertRuleAnomalyTactics(input *[]securityinsight.AttackTactic) []string {
	if input == nil {
		return []string{}
	}

	output := make([]string, 0)
	for _, v := range *input {
		output = append(output, string(v))
	}

	return output
}
