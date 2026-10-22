// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_processing_rule_action_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertProcessingRulesClient

			id, err := alertprocessingrules.ParseActionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetByName(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := AlertProcessingRuleActionGroupModel{
				Name:              id.ActionRuleName,
				ResourceGroupName: id.ResourceGroupName,
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model is null", *id)
			}
			properties := model.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: property is null", *id)
			}

			addActionGroupID, err := flattenAlertProcessingRuleAddActionGroupId(properties.Actions)
			if err != nil {
				return err
			}
			state.AddActionGroupIds = addActionGroupID

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.Enabled != nil {
				state.Enabled = *properties.Enabled
			}

			state.Scopes = properties.Scopes

			if properties.Conditions != nil {
				state.Condition = helpers.FlattenAlertProcessingRuleConditions(properties.Conditions)
			}

			if properties.Schedule != nil {
				state.Schedule = helpers.FlattenAlertProcessingRuleSchedule(properties.Schedule)
			}

			return metadata.Encode(&state)
		},
	}
}
