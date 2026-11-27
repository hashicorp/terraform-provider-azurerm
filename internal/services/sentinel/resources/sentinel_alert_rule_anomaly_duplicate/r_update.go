// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_duplicate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleAnomalyDuplicateModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := sentinelClient.From(metadata.Client).AnalyticsSettingsClient

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

			if rd.HasChanges(
				"multi_select_observation",
				"single_select_observation",
				"prioritized_exclude_observation",
				"threshold_observation",
			) {
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
