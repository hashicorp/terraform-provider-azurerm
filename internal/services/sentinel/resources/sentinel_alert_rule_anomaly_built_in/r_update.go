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

func (r Resource) Update() sdk.ResourceFunc {
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
