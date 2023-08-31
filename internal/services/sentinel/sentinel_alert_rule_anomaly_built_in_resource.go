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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type AlertRuleAnomalyBuiltInModel struct {
	Name                         string                                  `tfschema:"name"`
	DisplayName                  string                                  `tfschema:"display_name"`
	WorkspaceId                  string                                  `tfschema:"log_analytics_workspace_id"`
	Enabled                      bool                                    `tfschema:"enabled"`
	Mode                         string                                  `tfschema:"mode"`
	AnomalyVersion               string                                  `tfschema:"anomaly_version"`
	AnomalySettingsVersion       int32                                   `tfschema:"anomaly_settings_version"`
	Description                  string                                  `tfschema:"description"`
	Frequency                    string                                  `tfschema:"frequency"`
	RequiredDataConnectors       []AnomalyRuleRequiredDataConnectorModel `tfschema:"required_data_connector"`
	SettingsDefinitionId         string                                  `tfschema:"settings_definition_id"`
	Tactics                      []string                                `tfschema:"tactics"`
	Techniques                   []string                                `tfschema:"techniques"`
	ThresholdObservation         []AnomalyRuleThresholdModel             `tfschema:"threshold_observation"`
	MultiSelectObservation       []AnomalyRuleMultiSelectModel           `tfschema:"multi_select_observation"`
	SingleSelectObservation      []AnomalyRuleSingleSelectModel          `tfschema:"single_select_observation"`
	PrioritizeExcludeObservation []AnomalyRulePriorityModel              `tfschema:"prioritized_exclude_observation"`
}

type AlertRuleAnomalyBuiltInResource struct{}

var _ sdk.ResourceWithUpdate = AlertRuleAnomalyBuiltInResource{}

func (r AlertRuleAnomalyBuiltInResource) ModelObject() interface{} {
	return &AlertRuleAnomalyBuiltInModel{}
}

func (r AlertRuleAnomalyBuiltInResource) ResourceType() string {
	return "azurerm_sentinel_alert_rule_anomaly_built_in"
}

func (r AlertRuleAnomalyBuiltInResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.MLAnalyticsSettingsID
}

func (r AlertRuleAnomalyBuiltInResource) Arguments() map[string]*schema.Schema {
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
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(securityinsight.SettingsStatusProduction),
				string(securityinsight.SettingsStatusFlighting),
			}, false),
		},
	}
}

func (r AlertRuleAnomalyBuiltInResource) Attributes() map[string]*schema.Schema {
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

func (r AlertRuleAnomalyBuiltInResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel AlertRuleAnomalyBuiltInModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			workspaceId, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			builtinRule, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, *workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.Name != nil && strings.EqualFold(*v.Name, metaModel.Name) {
					return true
				}

				if v.DisplayName != nil && strings.EqualFold(*v.DisplayName, metaModel.DisplayName) {
					return true
				}

				return false
			})

			if err != nil {
				return fmt.Errorf("reading: %+v", err)
			}
			if builtinRule == nil {
				if metaModel.DisplayName != "" {
					return fmt.Errorf("built in rule (Display Name %q) was not found", metaModel.DisplayName)
				}
				return fmt.Errorf("built in rule (Display Name %q) was not found", metaModel.Name)
			}

			id, err := parse.MLAnalyticsSettingsID(AlertRuleAnomalyIdFromWorkspaceId(*workspaceId, *builtinRule.Name))
			if err != nil {
				return fmt.Errorf("parsing: %+v", err)
			}

			param := securityinsight.AnomalySecurityMLAnalyticsSettings{
				Kind: securityinsight.KindBasicSecurityMLAnalyticsSettingKindAnomaly,
				AnomalySecurityMLAnalyticsSettingsProperties: &securityinsight.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:              builtinRule.Description,
					DisplayName:              builtinRule.DisplayName,
					RequiredDataConnectors:   builtinRule.RequiredDataConnectors,
					Tactics:                  builtinRule.Tactics,
					Techniques:               builtinRule.Techniques,
					AnomalyVersion:           builtinRule.AnomalyVersion,
					Frequency:                builtinRule.Frequency,
					IsDefaultSettings:        builtinRule.IsDefaultSettings,
					AnomalySettingsVersion:   builtinRule.AnomalySettingsVersion,
					SettingsDefinitionID:     builtinRule.SettingsDefinitionID,
					Enabled:                  utils.Bool(metaModel.Enabled),
					SettingsStatus:           securityinsight.SettingsStatus(metaModel.Mode),
					CustomizableObservations: builtinRule.CustomizableObservations,
				},
			}

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SecurityMLAnalyticsSettingName, param)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AlertRuleAnomalyBuiltInResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := parse.MLAnalyticsSettingsID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}
			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

			resp, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.ID != nil && strings.EqualFold(*v.ID, id.ID()) {
					return true
				}
				return false
			})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp == nil {
				return metadata.MarkAsGone(id)
			}

			state := AlertRuleAnomalyBuiltInModel{
				WorkspaceId: workspaceId.ID(),
				Mode:        string(resp.SettingsStatus),
			}

			if resp.Name != nil {
				state.Name = *resp.Name
			}

			if resp.DisplayName != nil {
				state.DisplayName = *resp.DisplayName
			}

			if resp.AnomalyVersion != nil {
				state.AnomalyVersion = *resp.AnomalyVersion
			}

			if resp.AnomalySettingsVersion != nil {
				state.AnomalySettingsVersion = *resp.AnomalySettingsVersion
			}

			if resp.Description != nil {
				state.Description = *resp.Description
			}

			if resp.Enabled != nil {
				state.Enabled = *resp.Enabled
			}

			if resp.Frequency != nil {
				state.Frequency = *resp.Frequency
			}

			state.RequiredDataConnectors = flattenSentinelAlertRuleAnomalyRequiredDataConnectors(resp.RequiredDataConnectors)

			if resp.SettingsDefinitionID != nil {
				state.SettingsDefinitionId = resp.SettingsDefinitionID.String()
			}

			state.Tactics = flattenSentinelAlertRuleAnomalyTactics(resp.Tactics)

			if resp.Techniques != nil {
				state.Techniques = *resp.Techniques
			}

			if resp.CustomizableObservations != nil {
				state.MultiSelectObservation = flattenSentinelAlertRuleAnomalyMultiSelect(resp.CustomizableObservations.MultiSelectObservations)
				state.SingleSelectObservation = flattenSentinelAlertRuleAnomalySingleSelect(resp.CustomizableObservations.SingleSelectObservations)
				state.PrioritizeExcludeObservation = flattenSentinelAlertRuleAnomalyPriority(resp.CustomizableObservations.PrioritizeExcludeObservations)
				state.ThresholdObservation = flattenSentinelAlertRuleAnomalyThreshold(resp.CustomizableObservations.ThresholdObservations)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AlertRuleAnomalyBuiltInResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel AlertRuleAnomalyBuiltInModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := parse.MLAnalyticsSettingsID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing: %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			existing, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, *workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.ID != nil && strings.EqualFold(*v.ID, id.ID()) {
					return true
				}
				return false
			})

			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing == nil {
				return fmt.Errorf("retrieving %s: not found", *id)
			}

			param := securityinsight.AnomalySecurityMLAnalyticsSettings{
				Kind: securityinsight.KindBasicSecurityMLAnalyticsSettingKindAnomaly,
				AnomalySecurityMLAnalyticsSettingsProperties: &securityinsight.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:              existing.Description,
					DisplayName:              existing.DisplayName,
					RequiredDataConnectors:   existing.RequiredDataConnectors,
					Tactics:                  existing.Tactics,
					Techniques:               existing.Techniques,
					AnomalyVersion:           existing.AnomalyVersion,
					Frequency:                existing.Frequency,
					IsDefaultSettings:        existing.IsDefaultSettings,
					AnomalySettingsVersion:   existing.AnomalySettingsVersion,
					SettingsDefinitionID:     existing.SettingsDefinitionID,
					Enabled:                  utils.Bool(metaModel.Enabled),
					SettingsStatus:           securityinsight.SettingsStatus(metaModel.Mode),
					CustomizableObservations: existing.CustomizableObservations,
				},
			}

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SecurityMLAnalyticsSettingName, param)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AlertRuleAnomalyBuiltInResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// it's not able to delete built-in rules.
			var metaModel AlertRuleAnomalyBuiltInModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := parse.MLAnalyticsSettingsID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing: %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			existing, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, *workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.ID != nil && strings.EqualFold(*v.ID, id.ID()) {
					return true
				}
				return false
			})

			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing == nil {
				return fmt.Errorf("retrieving %s: not found", *id)
			}

			param := securityinsight.AnomalySecurityMLAnalyticsSettings{
				Kind: securityinsight.KindBasicSecurityMLAnalyticsSettingKindAnomaly,
				AnomalySecurityMLAnalyticsSettingsProperties: &securityinsight.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:              existing.Description,
					DisplayName:              existing.DisplayName,
					RequiredDataConnectors:   existing.RequiredDataConnectors,
					Tactics:                  existing.Tactics,
					Techniques:               existing.Techniques,
					AnomalyVersion:           existing.AnomalyVersion,
					Frequency:                existing.Frequency,
					IsDefaultSettings:        existing.IsDefaultSettings,
					AnomalySettingsVersion:   existing.AnomalySettingsVersion,
					SettingsDefinitionID:     existing.SettingsDefinitionID,
					Enabled:                  utils.Bool(false),
					SettingsStatus:           securityinsight.SettingsStatus(metaModel.Mode),
					CustomizableObservations: existing.CustomizableObservations,
				},
			}

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SecurityMLAnalyticsSettingName, param)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
