// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AlertProcessingRuleSuppressionModel struct {
	Name              string                              `tfschema:"name"`
	ResourceGroupName string                              `tfschema:"resource_group_name"`
	Scopes            []string                            `tfschema:"scopes"`
	Description       string                              `tfschema:"description"`
	Enabled           bool                                `tfschema:"enabled"`
	Condition         []AlertProcessingRuleConditionModel `tfschema:"condition"`
	Schedule          []AlertProcessingRuleScheduleModel  `tfschema:"schedule"`
	Tags              map[string]string                   `tfschema:"tags"`
}

type AlertProcessingRuleSuppressionResource struct{}

var _ sdk.ResourceWithUpdate = AlertProcessingRuleSuppressionResource{}

func (r AlertProcessingRuleSuppressionResource) ResourceType() string {
	return "azurerm_monitor_alert_processing_rule_suppression"
}

func (r AlertProcessingRuleSuppressionResource) ModelObject() interface{} {
	return &AlertProcessingRuleSuppressionModel{}
}

func (r AlertProcessingRuleSuppressionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return alertprocessingrules.ValidateActionRuleID
}

func (r AlertProcessingRuleSuppressionResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := schemaAlertProcessingRule()
	return arguments
}

func (r AlertProcessingRuleSuppressionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AlertProcessingRuleSuppressionResource) Create() sdk.ResourceFunc {
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
			existing, err := client.GetByName(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			alertProcessingRule := alertprocessingrules.AlertProcessingRule{
				// Location support "global" only
				Location: "global",
				Properties: &alertprocessingrules.AlertProcessingRuleProperties{
					Actions: []alertprocessingrules.Action{
						alertprocessingrules.RemoveAllActionGroups{},
					},
					Conditions:  expandAlertProcessingRuleConditions(model.Condition),
					Description: utils.String(model.Description),
					Enabled:     utils.Bool(model.Enabled),
					Schedule:    expandAlertProcessingRuleSchedule(model.Schedule),
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

func (r AlertProcessingRuleSuppressionResource) Update() sdk.ResourceFunc {
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

			if metadata.ResourceData.HasChange("add_action_group_ids") {
				model.Properties.Actions = []alertprocessingrules.Action{
					alertprocessingrules.RemoveAllActionGroups{}}
			}

			if metadata.ResourceData.HasChange("condition") {
				model.Properties.Conditions = expandAlertProcessingRuleConditions(resourceModel.Condition)
			}

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = utils.String(resourceModel.Description)
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = utils.Bool(resourceModel.Enabled)
			}

			if metadata.ResourceData.HasChange("schedule") {
				model.Properties.Schedule = expandAlertProcessingRuleSchedule(resourceModel.Schedule)
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

func (r AlertProcessingRuleSuppressionResource) Read() sdk.ResourceFunc {
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

			state := AlertProcessingRuleSuppressionModel{
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

			state.Condition = flattenAlertProcessingRuleConditions(properties.Conditions)
			state.Schedule = flattenAlertProcessingRuleSchedule(properties.Schedule)

			return metadata.Encode(&state)
		},
	}
}
func (r AlertProcessingRuleSuppressionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertProcessingRulesClient

			id, err := alertprocessingrules.ParseActionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
