// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_built_in

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleAnomalyBuiltInModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := sentinelClient.From(metadata.Client).AnalyticsSettingsClient

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
