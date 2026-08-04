// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AlertRuleAnomalyDuplicateModel struct {
	Name                         string                                  `tfschema:"name"`
	DisplayName                  string                                  `tfschema:"display_name"`
	BuiltInRuleId                string                                  `tfschema:"built_in_rule_id"`
	WorkspaceId                  string                                  `tfschema:"log_analytics_workspace_id"`
	Enabled                      bool                                    `tfschema:"enabled"`
	Mode                         string                                  `tfschema:"mode"`
	AnomalyVersion               string                                  `tfschema:"anomaly_version"`
	AnomalySettingsVersion       int64                                   `tfschema:"anomaly_settings_version"`
	Description                  string                                  `tfschema:"description"`
	Frequency                    string                                  `tfschema:"frequency"`
	IsDefaultSettings            bool                                    `tfschema:"is_default_settings"`
	RequiredDataConnectors       []AnomalyRuleRequiredDataConnectorModel `tfschema:"required_data_connector"`
	SettingsDefinitionId         string                                  `tfschema:"settings_definition_id"`
	Tactics                      []string                                `tfschema:"tactics"`
	Techniques                   []string                                `tfschema:"techniques"`
	ThresholdObservation         []AnomalyRuleThresholdModel             `tfschema:"threshold_observation"`
	MultiSelectObservation       []AnomalyRuleMultiSelectModel           `tfschema:"multi_select_observation"`
	SingleSelectObservation      []AnomalyRuleSingleSelectModel          `tfschema:"single_select_observation"`
	PrioritizeExcludeObservation []AnomalyRulePriorityModel              `tfschema:"prioritized_exclude_observation"`
}

type AlertRuleAnomalyDuplicateResource struct{}

var _ sdk.ResourceWithUpdate = AlertRuleAnomalyDuplicateResource{}

func (r AlertRuleAnomalyDuplicateResource) ModelObject() interface{} {
	return &AlertRuleAnomalyDuplicateModel{}
}

func (r AlertRuleAnomalyDuplicateResource) ResourceType() string {
	return "azurerm_sentinel_alert_rule_anomaly_duplicate"
}

func (r AlertRuleAnomalyDuplicateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return securitymlanalyticssettings.ValidateSecurityMLAnalyticsSettingID
}

func (r AlertRuleAnomalyDuplicateResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"built_in_rule_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: securitymlanalyticssettings.ValidateSecurityMLAnalyticsSettingID,
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

		"multi_select_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"values": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
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
				},
			},
		},

		"single_select_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"prioritized_exclude_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"prioritize": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"exclude": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"threshold_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
						Required: true,
					},
				},
			},
		},
	}
}

func (r AlertRuleAnomalyDuplicateResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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

		"is_default_settings": {
			Type:     pluginsdk.TypeBool,
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
	}
}

func (r AlertRuleAnomalyDuplicateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleAnomalyDuplicateModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			workspaceId, err := workspaces.ParseWorkspaceID(config.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			items, err := client.ListComplete(ctx, securitymlanalyticssettings.WorkspaceId(*workspaceId))
			if err != nil {
				return fmt.Errorf("listing alerts rules on %s: %+v", workspaceId, err)
			}

			var builtInAnomalyRule *securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings
			for _, item := range items.Items {
				if strings.EqualFold(AlertRuleAnomalyIdFromWorkspaceId(*workspaceId, pointer.From(item.SecurityMLAnalyticsSetting().Name)), config.BuiltInRuleId) {
					v, ok := item.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings)
					if !ok {
						continue
					}
					builtInAnomalyRule = &v
					break
				}
			}

			if builtInAnomalyRule == nil {
				return fmt.Errorf("retrieving built-in anomaly rule (%s): not found", config.BuiltInRuleId)
			}

			if builtInAnomalyRule.Properties == nil {
				return fmt.Errorf("retrieving built-in anomaly rule (%s): `properties` was nil", config.BuiltInRuleId)
			}
			builtInAnomalyRuleProps := builtInAnomalyRule.Properties

			var duplicateRule *securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings
			for _, item := range items.Items {
				v, ok := item.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings)
				if !ok {
					continue
				}
				if v.Properties != nil && builtInAnomalyRule.Properties != nil {
					if vSettingsID := pointer.From(v.Properties.SettingsDefinitionId); strings.EqualFold(vSettingsID, pointer.From(builtInAnomalyRule.Properties.SettingsDefinitionId)) && vSettingsID != "" {
						if vName := pointer.From(v.Name); vName != pointer.From(builtInAnomalyRule.Name) && vName != "" {
							duplicateRule = &v
							break
						}
					}
				}
			}

			if duplicateRule != nil {
				parsedExistingId, err := securitymlanalyticssettings.ParseSecurityMLAnalyticsSettingID(AlertRuleAnomalyIdFromWorkspaceId(*workspaceId, pointer.From(duplicateRule.Name)))
				if err != nil {
					return err
				}
				return fmt.Errorf("only one duplicate rule of the same built-in rule is allowed, there is an existing duplicate rule with id %s", parsedExistingId.ID())
			}

			id := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, uuid.New().String())

			if builtInAnomalyRule.Properties.SettingsStatus == securitymlanalyticssettings.SettingsStatusProduction && config.Mode == string(securitymlanalyticssettings.SettingsStatusProduction) {
				return fmt.Errorf("built-in anomaly rule %s is in production mode, creating duplicate rules in production mode is not supported", builtInAnomalyRule.Properties.DisplayName)
			}

			param := securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings{
				Kind: securitymlanalyticssettings.SecurityMLAnalyticsSettingsKindAnomaly,
				Properties: &securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:            builtInAnomalyRuleProps.Description,
					DisplayName:            config.DisplayName,
					RequiredDataConnectors: builtInAnomalyRuleProps.RequiredDataConnectors,
					Tactics:                builtInAnomalyRuleProps.Tactics,
					Techniques:             builtInAnomalyRuleProps.Techniques,
					AnomalyVersion:         builtInAnomalyRuleProps.AnomalyVersion,
					Frequency:              builtInAnomalyRuleProps.Frequency,
					IsDefaultSettings:      false, // for duplicate one, it's not default settings.
					AnomalySettingsVersion: builtInAnomalyRuleProps.AnomalySettingsVersion,
					SettingsDefinitionId:   builtInAnomalyRuleProps.SettingsDefinitionId,
					Enabled:                config.Enabled,
					SettingsStatus:         securitymlanalyticssettings.SettingsStatusFlighting,
				},
			}

			customizableObservations := &securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations{}
			customizableObservations.MultiSelectObservations, err = expandAlertRuleAnomalyMultiSelectObservations(builtInAnomalyRuleProps.CustomizableObservations, config.MultiSelectObservation)
			if err != nil {
				return fmt.Errorf("expanding `multi_select_observation`: %+v", err)
			}
			customizableObservations.SingleSelectObservations, err = expandAlertRuleAnomalySingleSelectObservations(builtInAnomalyRuleProps.CustomizableObservations, config.SingleSelectObservation)
			if err != nil {
				return fmt.Errorf("expanding `single_select_observation`: %+v", err)
			}
			customizableObservations.PrioritizeExcludeObservations, err = expandAlertRuleAnomalyPrioritizeExcludeObservations(builtInAnomalyRuleProps.CustomizableObservations, config.PrioritizeExcludeObservation)
			if err != nil {
				return fmt.Errorf("expanding `prioritize_exclude_observation`: %+v", err)
			}
			customizableObservations.ThresholdObservations, err = expandAlertRuleAnomalyThresholdObservations(builtInAnomalyRuleProps.CustomizableObservations, config.ThresholdObservation)
			if err != nil {
				return fmt.Errorf("expanding `threshold_observation`: %+v", err)
			}
			param.Properties.CustomizableObservations = customizableObservations

			if _, err = client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AlertRuleAnomalyDuplicateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := securitymlanalyticssettings.ParseSecurityMLAnalyticsSettingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
			state := AlertRuleAnomalyDuplicateModel{
				WorkspaceId: workspaceId.ID(),
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
					state.IsDefaultSettings = props.IsDefaultSettings

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

					if props.SettingsDefinitionId != nil {
						state.BuiltInRuleId = AlertRuleAnomalyIdFromWorkspaceId(workspaceId, *props.SettingsDefinitionId)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AlertRuleAnomalyDuplicateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleAnomalyDuplicateModel
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
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			v, ok := existing.Model.(securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings)
			if !ok {
				return fmt.Errorf("retrieving %s: expected type `%T`, got `%T`", id, securitymlanalyticssettings.AnomalySecurityMLAnalyticsSettings{}, existing.Model)
			}

			if v.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			rd := metadata.ResourceData
			if rd.HasChange("display_name") {
				v.Properties.DisplayName = config.DisplayName
			}

			if rd.HasChange("enabled") {
				v.Properties.Enabled = config.Enabled
			}

			if rd.HasChange("mode") {
				v.Properties.SettingsStatus = securitymlanalyticssettings.SettingsStatus(config.Mode)
			}

			if rd.HasChanges("multi_select_observation", "single_select_observation", "prioritized_exclude_observation", "threshold_observation") {
				if v.Properties.CustomizableObservations == nil {
					v.Properties.CustomizableObservations = &securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations{}
				}

				if rd.HasChange("multi_select_observation") {
					v.Properties.CustomizableObservations.MultiSelectObservations, err = expandAlertRuleAnomalyMultiSelectObservations(v.Properties.CustomizableObservations, config.MultiSelectObservation)
					if err != nil {
						return fmt.Errorf("expanding `multi_select_observation`: %+v", err)
					}
				}

				if rd.HasChange("single_select_observation") {
					v.Properties.CustomizableObservations.SingleSelectObservations, err = expandAlertRuleAnomalySingleSelectObservations(v.Properties.CustomizableObservations, config.SingleSelectObservation)
					if err != nil {
						return fmt.Errorf("expanding `single_select_observation`: %+v", err)
					}
				}

				if rd.HasChange("prioritized_exclude_observation") {
					v.Properties.CustomizableObservations.PrioritizeExcludeObservations, err = expandAlertRuleAnomalyPrioritizeExcludeObservations(v.Properties.CustomizableObservations, config.PrioritizeExcludeObservation)
					if err != nil {
						return fmt.Errorf("expanding `prioritize_exclude_observation`: %+v", err)
					}
				}

				if rd.HasChange("threshold_observation") {
					v.Properties.CustomizableObservations.ThresholdObservations, err = expandAlertRuleAnomalyThresholdObservations(v.Properties.CustomizableObservations, config.ThresholdObservation)
					if err != nil {
						return fmt.Errorf("expanding `threshold_observation`: %+v", err)
					}
				}
			}

			if _, err = client.CreateOrUpdate(ctx, *id, v); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AlertRuleAnomalyDuplicateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := securitymlanalyticssettings.ParseSecurityMLAnalyticsSettingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandAlertRuleAnomalyMultiSelectObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []AnomalyRuleMultiSelectModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsMultiSelectObservations, error) {
	if builtInInput == nil || builtInInput.MultiSelectObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.MultiSelectObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `multi_select_observation` must equal or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]AnomalyRuleMultiSelectModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsMultiSelectObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := securitymlanalyticssettings.AnomalySecurityMLAnalyticsMultiSelectObservations{
			Name:               v.Name,
			Description:        v.Description,
			Values:             v.Values,
			SupportedValues:    v.SupportedValues,
			SupportedValuesKql: v.SupportedValuesKql,
			ValuesKql:          v.ValuesKql,
			SequenceNumber:     v.SequenceNumber,
			Rerun:              v.Rerun,
		}
		if in, ok := inputValueMap[strings.ToLower(*v.Name)]; ok {
			o.Values = &in.Values
			delete(inputValueMap, strings.ToLower(*v.Name))
		}
		output = append(output, o)
	}

	if len(inputValueMap) != 0 {
		keys := make([]string, 0)
		for k := range inputValueMap {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("the following `multi_select_observation` are not supported: %s", strings.Join(keys, ", "))
	}

	return &output, nil
}

func expandAlertRuleAnomalySingleSelectObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []AnomalyRuleSingleSelectModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsSingleSelectObservations, error) {
	if builtInInput == nil || builtInInput.SingleSelectObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.SingleSelectObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `single_select_observation` must equals or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]AnomalyRuleSingleSelectModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsSingleSelectObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := securitymlanalyticssettings.AnomalySecurityMLAnalyticsSingleSelectObservations{
			Name:               v.Name,
			Description:        v.Description,
			Value:              v.Value,
			SupportedValues:    v.SupportedValues,
			SupportedValuesKql: v.SupportedValuesKql,
			SequenceNumber:     v.SequenceNumber,
			Rerun:              v.Rerun,
		}
		if in, ok := inputValueMap[strings.ToLower(*v.Name)]; ok {
			o.Value = &in.Value
			delete(inputValueMap, strings.ToLower(*v.Name))
		}
		output = append(output, o)
	}

	if len(inputValueMap) != 0 {
		keys := make([]string, 0)
		for k := range inputValueMap {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("the following `single_select_observation` are not supported: %s", strings.Join(keys, ", "))
	}

	return &output, nil
}

func expandAlertRuleAnomalyPrioritizeExcludeObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []AnomalyRulePriorityModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations, error) {
	if builtInInput == nil || builtInInput.PrioritizeExcludeObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.PrioritizeExcludeObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `prioritized_exclude_observation` must equals or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]AnomalyRulePriorityModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := securitymlanalyticssettings.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations{
			Name:           v.Name,
			Description:    v.Description,
			Prioritize:     v.Prioritize,
			Exclude:        v.Exclude,
			DataType:       v.DataType,
			SequenceNumber: v.SequenceNumber,
			Rerun:          v.Rerun,
		}
		if in, ok := inputValueMap[strings.ToLower(*v.Name)]; ok {
			o.Exclude = &in.Exclude
			o.Prioritize = &in.Prioritize
			delete(inputValueMap, strings.ToLower(*v.Name))
		}
		output = append(output, o)
	}

	if len(inputValueMap) != 0 {
		keys := make([]string, 0)
		for k := range inputValueMap {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("the following `prioritized_exclude_observation` are not supported: %s", strings.Join(keys, ", "))
	}

	return &output, nil
}

func expandAlertRuleAnomalyThresholdObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []AnomalyRuleThresholdModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsThresholdObservations, error) {
	if builtInInput == nil || builtInInput.ThresholdObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.ThresholdObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `threshold_observation` must equals or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]AnomalyRuleThresholdModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsThresholdObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := securitymlanalyticssettings.AnomalySecurityMLAnalyticsThresholdObservations{
			Name:           v.Name,
			Description:    v.Description,
			Maximum:        v.Maximum,
			Minimum:        v.Minimum,
			Value:          v.Value,
			SequenceNumber: v.SequenceNumber,
			Rerun:          v.Rerun,
		}
		if in, ok := inputValueMap[strings.ToLower(*v.Name)]; ok {
			o.Value = &in.Value
			delete(inputValueMap, strings.ToLower(*v.Name))
		}
		output = append(output, o)
	}

	if len(inputValueMap) != 0 {
		keys := make([]string, 0)
		for k := range inputValueMap {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("the following `threshold_observation` are not supported: %s", strings.Join(keys, ", "))
	}

	return &output, nil
}
