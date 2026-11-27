// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_duplicate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
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
				return err
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

func expandAlertRuleAnomalyMultiSelectObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []helpers.AnomalyRuleMultiSelectModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsMultiSelectObservations, error) {
	if builtInInput == nil || builtInInput.MultiSelectObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.MultiSelectObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `multi_select_observation` must equal or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]helpers.AnomalyRuleMultiSelectModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsMultiSelectObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, errors.New("the name of built in `multi_select_observation` is nil")
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

func expandAlertRuleAnomalySingleSelectObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []helpers.AnomalyRuleSingleSelectModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsSingleSelectObservations, error) {
	if builtInInput == nil || builtInInput.SingleSelectObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.SingleSelectObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `single_select_observation` must equals or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]helpers.AnomalyRuleSingleSelectModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsSingleSelectObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, errors.New("the name of built in `multi_select_observation` is nil")
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

func expandAlertRuleAnomalyPrioritizeExcludeObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []helpers.AnomalyRulePriorityModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations, error) {
	if builtInInput == nil || builtInInput.PrioritizeExcludeObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.PrioritizeExcludeObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `prioritized_exclude_observation` must equals or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]helpers.AnomalyRulePriorityModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsPrioritizeExcludeObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, errors.New("the name of built in `multi_select_observation` is nil")
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

func expandAlertRuleAnomalyThresholdObservations(builtInInput *securitymlanalyticssettings.AnomalySecurityMLAnalyticsCustomizableObservations, input []helpers.AnomalyRuleThresholdModel) (*[]securitymlanalyticssettings.AnomalySecurityMLAnalyticsThresholdObservations, error) {
	if builtInInput == nil || builtInInput.ThresholdObservations == nil {
		return nil, nil
	}
	builtInRule := builtInInput.ThresholdObservations

	if len(*builtInRule) < len(input) {
		return nil, fmt.Errorf("the number of `threshold_observation` must equals or less than %d", len(*builtInRule))
	}

	inputValueMap := make(map[string]helpers.AnomalyRuleThresholdModel)
	for _, v := range input {
		inputValueMap[strings.ToLower(v.Name)] = v
	}

	output := make([]securitymlanalyticssettings.AnomalySecurityMLAnalyticsThresholdObservations, 0)
	for _, v := range *builtInRule {
		if v.Name == nil {
			return nil, errors.New("the name of built in `multi_select_observation` is nil")
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
