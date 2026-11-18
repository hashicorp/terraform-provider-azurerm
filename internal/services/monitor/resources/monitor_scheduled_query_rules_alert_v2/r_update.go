// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert_v2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.ScheduledQueryRulesV2Client

			id, err := scheduledqueryrules.ParseScheduledQueryRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var resourceModel ScheduledQueryRulesAlertV2Model
			if err := metadata.Decode(&resourceModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if metadata.ResourceData.HasChange("action") {
				model.Properties.Actions = expandScheduledQueryRulesAlertV2ActionsModel(resourceModel.Actions)
			}

			if metadata.ResourceData.HasChange("auto_mitigation_enabled") {
				model.Properties.AutoMitigate = &resourceModel.AutoMitigate
			}

			if metadata.ResourceData.HasChange("workspace_alerts_storage_enabled") {
				model.Properties.CheckWorkspaceAlertsStorageConfigured = &resourceModel.CheckWorkspaceAlertsStorageConfigured
			}

			if metadata.ResourceData.HasChange("criteria") {
				model.Properties.Criteria = expandScheduledQueryRulesAlertV2CriteriaModel(resourceModel.Criteria)
			}

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = &resourceModel.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				if resourceModel.DisplayName != "" {
					model.Properties.DisplayName = &resourceModel.DisplayName
				} else {
					model.Properties.DisplayName = nil
				}
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = &resourceModel.Enabled
			}

			if metadata.ResourceData.HasChange("evaluation_frequency") {
				model.Properties.EvaluationFrequency = &resourceModel.EvaluationFrequency
			}

			if metadata.ResourceData.HasChange("mute_actions_after_alert_duration") {
				if resourceModel.MuteActionsDuration != "" {
					if resourceModel.AutoMitigate {
						return fmt.Errorf("auto mitigation must be disabled when mute action duration is set")
					}
					model.Properties.MuteActionsDuration = &resourceModel.MuteActionsDuration
				} else {
					model.Properties.MuteActionsDuration = nil
				}
			}

			if metadata.ResourceData.HasChange("query_time_range_override") {
				if resourceModel.OverrideQueryTimeRange != "" {
					model.Properties.OverrideQueryTimeRange = &resourceModel.OverrideQueryTimeRange
				} else {
					model.Properties.OverrideQueryTimeRange = nil
				}
			}

			if metadata.ResourceData.HasChange("severity") {
				model.Properties.Severity = &resourceModel.Severity
			}

			if metadata.ResourceData.HasChange("skip_query_validation") {
				model.Properties.SkipQueryValidation = &resourceModel.SkipQueryValidation
			}

			if metadata.ResourceData.HasChange("target_resource_types") {
				model.Properties.TargetResourceTypes = &resourceModel.TargetResourceTypes
			}

			if metadata.ResourceData.HasChange("window_duration") {
				if resourceModel.WindowSize != "" {
					model.Properties.WindowSize = &resourceModel.WindowSize
				} else {
					model.Properties.WindowSize = nil
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				if len(resourceModel.Identity) != 0 {
					model.Identity, err = identity.ExpandSystemOrUserAssignedMapFromModel(resourceModel.Identity)
					if err != nil {
						return fmt.Errorf("expanding SystemOrUserAssigned Identity: %+v", err)
					}
				} else {
					model.Identity = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &resourceModel.Tags
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
