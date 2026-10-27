// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_processing_rule_suppression

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AlertProcessingRuleSuppressionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.Monitor.AlertProcessingRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := alertprocessingrules.NewActionRuleID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetByName(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			alertProcessingRule := alertprocessingrules.AlertProcessingRule{
				// Location support "global" only
				Location: "global",
				Properties: &alertprocessingrules.AlertProcessingRuleProperties{
					Actions: []alertprocessingrules.Action{
						alertprocessingrules.RemoveAllActionGroups{},
					},
					Conditions:  helpers.ExpandAlertProcessingRuleConditions(model.Condition),
					Description: pointer.To(model.Description),
					Enabled:     pointer.To(model.Enabled),
					Schedule:    helpers.ExpandAlertProcessingRuleSchedule(model.Schedule),
					Scopes:      model.Scopes,
				},
				Tags: &model.Tags,
			}

			if _, err := client.CreateOrUpdate(ctx, id, alertProcessingRule); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
