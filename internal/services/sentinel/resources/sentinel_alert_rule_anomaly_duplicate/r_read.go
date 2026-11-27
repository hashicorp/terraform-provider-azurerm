// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_duplicate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
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

					state.RequiredDataConnectors = helpers.FlattenSentinelAlertRuleAnomalyRequiredDataConnectors(props.RequiredDataConnectors)
					state.SettingsDefinitionId = pointer.From(props.SettingsDefinitionId)
					state.Tactics = pointer.FromEnumSlice(props.Tactics)
					state.Techniques = pointer.From(props.Techniques)

					if co := props.CustomizableObservations; co != nil {
						state.MultiSelectObservation = helpers.FlattenSentinelAlertRuleAnomalyMultiSelect(co.MultiSelectObservations)
						state.SingleSelectObservation = helpers.FlattenSentinelAlertRuleAnomalySingleSelect(co.SingleSelectObservations)
						state.PrioritizeExcludeObservation = helpers.FlattenSentinelAlertRuleAnomalyPriority(co.PrioritizeExcludeObservations)
						state.ThresholdObservation = helpers.FlattenSentinelAlertRuleAnomalyThreshold(co.ThresholdObservations)
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
