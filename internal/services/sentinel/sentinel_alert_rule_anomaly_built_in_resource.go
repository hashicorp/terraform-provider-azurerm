// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AlertRuleAnomalyBuiltInModel struct {
	Name                         string                                  `tfschema:"name"`
	DisplayName                  string                                  `tfschema:"display_name"`
	WorkspaceId                  string                                  `tfschema:"log_analytics_workspace_id"`
	Enabled                      bool                                    `tfschema:"enabled"`
	Mode                         string                                  `tfschema:"mode"`
	AnomalyVersion               string                                  `tfschema:"anomaly_version"`
	AnomalySettingsVersion       int64                                   `tfschema:"anomaly_settings_version"`
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
	return securitymlanalyticssettings.ValidateSecurityMLAnalyticsSettingID
}

func (r AlertRuleAnomalyBuiltInResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
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
				string(securitymlanalyticssettings.SettingsStatusProduction),
				string(securitymlanalyticssettings.SettingsStatusFlighting),
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
			var config AlertRuleAnomalyBuiltInModel
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

			var builtInAnomalyRule *securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings
			var builtInRuleName string
			for _, item := range items.Items {
				v, ok := item.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings)
				if !ok {
					continue
				}

				switch {
				case config.Name != "":
					if strings.EqualFold(pointer.From(v.Name), config.Name) {
						builtInRuleName = config.Name
						builtInAnomalyRule = &v
						break
					}
				case config.DisplayName != "":
					if v.Properties != nil && strings.EqualFold(v.Properties.DisplayName, config.DisplayName) {
						builtInRuleName = config.DisplayName
						builtInAnomalyRule = &v
						break
					}
				}
			}

			if builtInAnomalyRule == nil {
				return fmt.Errorf("retrieving built-in anomaly rule (%s): not found", builtInRuleName)
			}

			if builtInAnomalyRule.Properties == nil {
				return fmt.Errorf("retrieving built-in anomaly rule (%s): `properties` was nil", builtInRuleName)
			}
			builtInAnomalyRuleProps := builtInAnomalyRule.Properties

			id := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, pointer.From(builtInAnomalyRule.Name))

			param := securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings{
				Kind: securitymlanalyticssettings.SecurityMLAnalyticsSettingsKindAnomaly,
				Properties: &securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:              builtInAnomalyRuleProps.Description,
					DisplayName:              builtInAnomalyRuleProps.DisplayName,
					RequiredDataConnectors:   builtInAnomalyRuleProps.RequiredDataConnectors,
					Tactics:                  builtInAnomalyRuleProps.Tactics,
					Techniques:               builtInAnomalyRuleProps.Techniques,
					AnomalyVersion:           builtInAnomalyRuleProps.AnomalyVersion,
					Frequency:                builtInAnomalyRuleProps.Frequency,
					IsDefaultSettings:        builtInAnomalyRuleProps.IsDefaultSettings,
					AnomalySettingsVersion:   builtInAnomalyRuleProps.AnomalySettingsVersion,
					SettingsDefinitionId:     builtInAnomalyRuleProps.SettingsDefinitionId,
					Enabled:                  config.Enabled,
					SettingsStatus:           securitymlanalyticssettings.SettingsStatus(config.Mode),
					CustomizableObservations: builtInAnomalyRuleProps.CustomizableObservations,
				},
			}

			if _, err = client.CreateOrUpdate(ctx, id, param); err != nil {
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

			id, err := securitymlanalyticssettings.ParseSecurityMLAnalyticsSettingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AlertRuleAnomalyBuiltInModel{
				WorkspaceId: workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Name = pointer.From(model.SecurityMLAnalyticsSetting().Name)
				if v, ok := model.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings); ok && v.Properties != nil {
					props := v.Properties

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
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AlertRuleAnomalyBuiltInResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleAnomalyBuiltInModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := securitymlanalyticssettings.ParseSecurityMLAnalyticsSettingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			v, ok := existing.Model.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings)
			if !ok {
				return fmt.Errorf("retrieving %s: expected type `%T`, got `%T`", id, securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings{}, existing.Model)
			}

			if v.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := v.Properties

			rd := metadata.ResourceData

			if rd.HasChange("enabled") {
				props.Enabled = config.Enabled
			}

			if rd.HasChange("mode") {
				props.SettingsStatus = securitymlanalyticssettings.SettingsStatus(config.Mode)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, v); err != nil {
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
			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := securitymlanalyticssettings.ParseSecurityMLAnalyticsSettingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			v, ok := existing.Model.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings)
			if !ok {
				return fmt.Errorf("retrieving %s: expected type `%T`, got `%T`", id, securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings{}, existing.Model)
			}

			if v.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := v.Properties

			param := securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings{
				Kind: securitymlanalyticssettings.SecurityMLAnalyticsSettingsKindAnomaly,
				Properties: &securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:              props.Description,
					DisplayName:              props.DisplayName,
					RequiredDataConnectors:   props.RequiredDataConnectors,
					Tactics:                  props.Tactics,
					Techniques:               props.Techniques,
					AnomalyVersion:           props.AnomalyVersion,
					Frequency:                props.Frequency,
					IsDefaultSettings:        props.IsDefaultSettings,
					AnomalySettingsVersion:   props.AnomalySettingsVersion,
					SettingsDefinitionId:     props.SettingsDefinitionId,
					Enabled:                  false,
					SettingsStatus:           props.SettingsStatus,
					CustomizableObservations: props.CustomizableObservations,
				},
			}

			if _, err = client.CreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
