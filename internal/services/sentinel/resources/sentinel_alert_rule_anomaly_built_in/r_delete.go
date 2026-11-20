// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_built_in

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
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
