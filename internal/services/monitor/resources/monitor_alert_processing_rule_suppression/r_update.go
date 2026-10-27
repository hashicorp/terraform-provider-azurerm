// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_processing_rule_suppression

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertProcessingRulesClient

			id, err := alertprocessingrules.ParseActionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var resourceModel AlertProcessingRuleSuppressionModel
			if err := metadata.Decode(&resourceModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.GetByName(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			model := resp.Model
			if model.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			if metadata.ResourceData.HasChange("condition") {
				model.Properties.Conditions = helpers.ExpandAlertProcessingRuleConditions(resourceModel.Condition)
			}

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = pointer.To(resourceModel.Description)
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = pointer.To(resourceModel.Enabled)
			}

			if metadata.ResourceData.HasChange("schedule") {
				model.Properties.Schedule = helpers.ExpandAlertProcessingRuleSchedule(resourceModel.Schedule)
			}

			if metadata.ResourceData.HasChange("scopes") {
				model.Properties.Scopes = resourceModel.Scopes
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
