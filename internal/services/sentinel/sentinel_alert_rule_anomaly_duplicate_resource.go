// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	return validate.MLAnalyticsSettingsID
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
			ValidateFunc: validate.MLAnalyticsSettingsID,
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
			var metaModel AlertRuleAnomalyDuplicateModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			workspaceId, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			builtInAnomalyRule, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, *workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.ID != nil && strings.EqualFold(AlertRuleAnomalyIdFromWorkspaceId(*workspaceId, *v.Name), metaModel.BuiltInRuleId) {
					return true
				}

				return false
			})

			if err != nil {
				return fmt.Errorf("reading built-in anomaly rule: %+v", err)
			}
			if builtInAnomalyRule == nil {
				return fmt.Errorf("built-in anomaly rule not found")
			}

			existingDuplicateRule, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, *workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.SettingsDefinitionID != nil &&
					builtInAnomalyRule.SettingsDefinitionID != nil &&
					strings.EqualFold(v.SettingsDefinitionID.String(), builtInAnomalyRule.SettingsDefinitionID.String()) &&
					v.Name != nil && builtInAnomalyRule.Name != nil && *v.Name != *builtInAnomalyRule.Name {
					return true
				}
				return false
			})
			if err != nil {
				return fmt.Errorf("checking for presence of existing duplicate rule of built-in rule: %+v", err)
			}
			if existingDuplicateRule != nil {
				parsedExistingId, err := parse.MLAnalyticsSettingsID(AlertRuleAnomalyIdFromWorkspaceId(*workspaceId, *existingDuplicateRule.Name))
				if err != nil {
					return fmt.Errorf("parsing: %+v", err)
				}
				return fmt.Errorf("only one duplicate rule of the same built-in rule is allowed, there is an existing duplicate rule of %s with id %q", *builtInAnomalyRule.DisplayName, parsedExistingId.ID())
			}

			id := parse.NewMLAnalyticsSettingsID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, uuid.New().String())
			// no need to do another existing check, it will be checked by finding existing duplicate rule of the template.

			if builtInAnomalyRule.SettingsStatus == securityinsight.SettingsStatusProduction && metaModel.Mode == string(securityinsight.SettingsStatusProduction) {
				return fmt.Errorf("built-in anomaly rule %s is in production mode, it's not allowed to create duplicate rule in production mode", *builtInAnomalyRule.DisplayName)
			}

			param := securityinsight.AnomalySecurityMLAnalyticsSettings{
				Kind: securityinsight.KindBasicSecurityMLAnalyticsSettingKindAnomaly,
				AnomalySecurityMLAnalyticsSettingsProperties: &securityinsight.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:            builtInAnomalyRule.Description,
					DisplayName:            utils.String(metaModel.DisplayName),
					RequiredDataConnectors: builtInAnomalyRule.RequiredDataConnectors,
					Tactics:                builtInAnomalyRule.Tactics,
					Techniques:             builtInAnomalyRule.Techniques,
					AnomalyVersion:         builtInAnomalyRule.AnomalyVersion,
					Frequency:              builtInAnomalyRule.Frequency,
					IsDefaultSettings:      utils.Bool(false), // for duplicate one, it's not default settings.
					AnomalySettingsVersion: builtInAnomalyRule.AnomalySettingsVersion,
					SettingsDefinitionID:   builtInAnomalyRule.SettingsDefinitionID,
					Enabled:                utils.Bool(metaModel.Enabled),
					SettingsStatus:         securityinsight.SettingsStatusFlighting,
				},
			}

			customizableObservations := &azuresdkhacks.AnomalySecurityMLAnalyticsCustomizableObservations{}
			customizableObservations.MultiSelectObservations, err = expandAlertRuleAnomalyMultiSelectObservations(builtInAnomalyRule.CustomizableObservations.MultiSelectObservations, metaModel.MultiSelectObservation)
			if err != nil {
				return fmt.Errorf("expanding `multi_select_observation`: %+v", err)
			}
			customizableObservations.SingleSelectObservations, err = expandAlertRuleAnomalySingleSelectObservations(builtInAnomalyRule.CustomizableObservations.SingleSelectObservations, metaModel.SingleSelectObservation)
			if err != nil {
				return fmt.Errorf("expanding `single_select_observation`: %+v", err)
			}
			customizableObservations.PrioritizeExcludeObservations, err = expandAlertRuleAnomalyPrioritizeExcludeObservations(builtInAnomalyRule.CustomizableObservations.PrioritizeExcludeObservations, metaModel.PrioritizeExcludeObservation)
			if err != nil {
				return fmt.Errorf("expanding `prioritize_exclude_observation`: %+v", err)
			}
			customizableObservations.ThresholdObservations, err = expandAlertRuleAnomalyThresholdObservations(builtInAnomalyRule.CustomizableObservations.ThresholdObservations, metaModel.ThresholdObservation)
			if err != nil {
				return fmt.Errorf("expanding `threshold_observation`: %+v", err)
			}

			param.AnomalySecurityMLAnalyticsSettingsProperties.CustomizableObservations = customizableObservations

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SecurityMLAnalyticsSettingName, param)
			if err != nil {
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

			state := AlertRuleAnomalyDuplicateModel{
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
				state.AnomalySettingsVersion = int64(*resp.AnomalySettingsVersion)
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
			if resp.IsDefaultSettings != nil {
				state.IsDefaultSettings = *resp.IsDefaultSettings
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

			if resp.SettingsDefinitionID != nil {
				state.BuiltInRuleId = AlertRuleAnomalyIdFromWorkspaceId(workspaceId, resp.SettingsDefinitionID.String())
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AlertRuleAnomalyDuplicateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel AlertRuleAnomalyDuplicateModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.AnalyticsSettingsClient

			id, err := parse.MLAnalyticsSettingsID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}
			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

			existing, err := AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, workspaceId, func(v *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
				if v.ID != nil && strings.EqualFold(*v.ID, id.ID()) {
					return true
				}

				return false
			})

			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing == nil {
				return fmt.Errorf("retrieving %s: Alert Rule Anomaly not found", *id)
			}

			param := securityinsight.AnomalySecurityMLAnalyticsSettings{
				Kind: securityinsight.KindBasicSecurityMLAnalyticsSettingKindAnomaly,
				AnomalySecurityMLAnalyticsSettingsProperties: &securityinsight.AnomalySecurityMLAnalyticsSettingsProperties{
					Description:            existing.Description,
					DisplayName:            existing.DisplayName,
					RequiredDataConnectors: existing.RequiredDataConnectors,
					Tactics:                existing.Tactics,
					Techniques:             existing.Techniques,
					AnomalyVersion:         existing.AnomalyVersion,
					Frequency:              existing.Frequency,
					IsDefaultSettings:      existing.IsDefaultSettings,
					AnomalySettingsVersion: existing.AnomalySettingsVersion,
					SettingsDefinitionID:   existing.SettingsDefinitionID,
					Enabled:                utils.Bool(metaModel.Enabled),
					SettingsStatus:         securityinsight.SettingsStatus(metaModel.Mode),
				},
			}

			customizableObservations := &azuresdkhacks.AnomalySecurityMLAnalyticsCustomizableObservations{}
			customizableObservations.MultiSelectObservations, err = expandAlertRuleAnomalyMultiSelectObservations(existing.CustomizableObservations.MultiSelectObservations, metaModel.MultiSelectObservation)
			if err != nil {
				return fmt.Errorf("expanding `multi_select_observation`: %+v", err)
			}
			customizableObservations.SingleSelectObservations, err = expandAlertRuleAnomalySingleSelectObservations(existing.CustomizableObservations.SingleSelectObservations, metaModel.SingleSelectObservation)
			if err != nil {
				return fmt.Errorf("expanding `single_select_observation`: %+v", err)
			}
			customizableObservations.PrioritizeExcludeObservations, err = expandAlertRuleAnomalyPrioritizeExcludeObservations(existing.CustomizableObservations.PrioritizeExcludeObservations, metaModel.PrioritizeExcludeObservation)
			if err != nil {
				return fmt.Errorf("expanding `prioritize_exclude_observation`: %+v", err)
			}
			customizableObservations.ThresholdObservations, err = expandAlertRuleAnomalyThresholdObservations(existing.CustomizableObservations.ThresholdObservations, metaModel.ThresholdObservation)
			if err != nil {
				return fmt.Errorf("expanding `threshold_observation`: %+v", err)
			}

			param.AnomalySecurityMLAnalyticsSettingsProperties.CustomizableObservations = customizableObservations

			_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SecurityMLAnalyticsSettingName, param)
			if err != nil {
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

			id, err := parse.MLAnalyticsSettingsID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			_, err = client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.SecurityMLAnalyticsSettingName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandAlertRuleAnomalyMultiSelectObservations(builtInRule *[]azuresdkhacks.AnomalySecurityMLAnalyticsMultiSelectObservations, input []AnomalyRuleMultiSelectModel) (*[]azuresdkhacks.AnomalySecurityMLAnalyticsMultiSelectObservations, error) {
	if builtInRule != nil && len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `multi_select_observation` must equal or less than %d", len(*builtInRule))
	}

	if builtInRule == nil {
		return nil, nil
	}

	inputValueMap := make(map[string]AnomalyRuleMultiSelectModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]azuresdkhacks.AnomalySecurityMLAnalyticsMultiSelectObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := azuresdkhacks.AnomalySecurityMLAnalyticsMultiSelectObservations{
			Name:               v.Name,
			Description:        v.Description,
			Values:             v.Values,
			SupportValues:      v.SupportValues,
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

func expandAlertRuleAnomalySingleSelectObservations(builtInRule *[]azuresdkhacks.AnomalySecurityMLAnalyticsSingleSelectObservations, input []AnomalyRuleSingleSelectModel) (*[]azuresdkhacks.AnomalySecurityMLAnalyticsSingleSelectObservations, error) {
	if builtInRule != nil && len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `single_select_observation` must equals or less than %d", len(*builtInRule))
	}

	if builtInRule == nil {
		return nil, nil
	}

	inputValueMap := make(map[string]AnomalyRuleSingleSelectModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]azuresdkhacks.AnomalySecurityMLAnalyticsSingleSelectObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := azuresdkhacks.AnomalySecurityMLAnalyticsSingleSelectObservations{
			Name:               v.Name,
			Description:        v.Description,
			Value:              v.Value,
			SupportValues:      v.SupportValues,
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

func expandAlertRuleAnomalyPrioritizeExcludeObservations(builtInRule *[]azuresdkhacks.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations, input []AnomalyRulePriorityModel) (*[]azuresdkhacks.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations, error) {
	if builtInRule != nil && len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `prioritized_exclude_observation` must equals or less than %d", len(*builtInRule))
	}

	if builtInRule == nil {
		return nil, nil
	}

	inputValueMap := make(map[string]AnomalyRulePriorityModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]azuresdkhacks.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := azuresdkhacks.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations{
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

func expandAlertRuleAnomalyThresholdObservations(builtInRule *[]azuresdkhacks.AnomalySecurityMLAnalyticsThresholdObservations, input []AnomalyRuleThresholdModel) (*[]azuresdkhacks.AnomalySecurityMLAnalyticsThresholdObservations, error) {
	if builtInRule != nil && len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `threshold_observation` must equals or less than %d", len(*builtInRule))
	}

	if builtInRule == nil {
		return nil, nil
	}

	inputValueMap := make(map[string]AnomalyRuleThresholdModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]azuresdkhacks.AnomalySecurityMLAnalyticsThresholdObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, fmt.Errorf("the name of built in `multi_select_observation` is nil")
		}
		// copy from built in rule
		o := azuresdkhacks.AnomalySecurityMLAnalyticsThresholdObservations{
			Name:           v.Name,
			Description:    v.Description,
			Max:            v.Max,
			Min:            v.Min,
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
