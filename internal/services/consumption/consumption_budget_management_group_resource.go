// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	validateManagementGroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagementGroupConsumptionBudget struct {
	base consumptionBudgetBaseResource
}

type ManagementGroupConsumptionBudgetModel struct {
	Name              string                                   `tfschema:"name"`
	ManagementGroupId string                                   `tfschema:"management_group_id"`
	Etag              string                                   `tfschema:"etag"`
	Amount            float64                                  `tfschema:"amount"`
	TimeGrain         string                                   `tfschema:"time_grain"`
	TimePeriod        []ConsumptionBudgetTimePeriodModel       `tfschema:"time_period"`
	Notification      []ConsumptionBudgetMgmtNotificationModel `tfschema:"notification"`
	Filter            []ConsumptionBudgetFilterModel           `tfschema:"filter"`
}

// Management group notification model (no contact_groups/contact_roles)
type ConsumptionBudgetMgmtNotificationModel struct {
	Enabled       bool     `tfschema:"enabled"`
	Threshold     int64    `tfschema:"threshold"`
	ThresholdType string   `tfschema:"threshold_type"`
	Operator      string   `tfschema:"operator"`
	ContactEmails []string `tfschema:"contact_emails"`
}

var (
	_ sdk.Resource                   = ManagementGroupConsumptionBudget{}
	_ sdk.ResourceWithCustomImporter = ManagementGroupConsumptionBudget{}
)

func (r ManagementGroupConsumptionBudget) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateManagementGroup.ManagementGroupID,
		},

		// Consumption Budgets for Management Groups have a different notification schema,
		// here we override the notification schema in the base resource
		"notification": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"threshold": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"threshold_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(budgets.ThresholdTypeActual),
						ValidateFunc: validation.StringInSlice([]string{
							string(budgets.ThresholdTypeActual),
							"Forecasted",
						}, false),
					},
					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(budgets.OperatorTypeEqualTo),
							string(budgets.OperatorTypeGreaterThan),
							string(budgets.OperatorTypeGreaterThanOrEqualTo),
						}, false),
					},

					"contact_emails": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
	return r.base.arguments(schema)
}

func (r ManagementGroupConsumptionBudget) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ManagementGroupConsumptionBudget) ModelObject() interface{} {
	return &ManagementGroupConsumptionBudgetModel{}
}

func (r ManagementGroupConsumptionBudget) ResourceType() string {
	return "azurerm_consumption_budget_management_group"
}

func (r ManagementGroupConsumptionBudget) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return budgets.ValidateScopedBudgetID
}

func (r ManagementGroupConsumptionBudget) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			var config ManagementGroupConsumptionBudgetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := budgets.NewScopedBudgetID(config.ManagementGroupId, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			timePeriod, err := expandConsumptionBudgetTimePeriodFromModel(config.TimePeriod)
			if err != nil {
				return fmt.Errorf("expanding `time_period`: %+v", err)
			}

			parameters := budgets.Budget{
				Name: pointer.To(id.BudgetName),
				Properties: &budgets.BudgetProperties{
					Amount:        config.Amount,
					Category:      budgets.CategoryTypeCost,
					Filter:        expandConsumptionBudgetFilterFromModel(config.Filter),
					Notifications: expandConsumptionBudgetMgmtNotificationsFromModel(config.Notification),
					TimeGrain:     budgets.TimeGrainType(config.TimeGrain),
					TimePeriod:    *timePeriod,
				},
			}

			if config.Etag != "" {
				parameters.ETag = pointer.To(config.Etag)
			}

			if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagementGroupConsumptionBudget) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient
			id, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s, %+v", *id, err)
			}

			state := ManagementGroupConsumptionBudgetModel{
				Name:              id.BudgetName,
				ManagementGroupId: id.Scope,
			}

			if model := resp.Model; model != nil {
				state.Etag = pointer.From(model.ETag)

				if props := model.Properties; props != nil {
					state.Amount = props.Amount
					state.TimeGrain = string(props.TimeGrain)
					state.TimePeriod = flattenConsumptionBudgetTimePeriodToModel(&props.TimePeriod)
					state.Notification = flattenConsumptionBudgetMgmtNotificationsToModel(props.Notifications)
					state.Filter = flattenConsumptionBudgetFilterToModel(props.Filter)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagementGroupConsumptionBudget) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ManagementGroupConsumptionBudget) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			id, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ManagementGroupConsumptionBudgetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			timePeriod, err := expandConsumptionBudgetTimePeriodFromModel(config.TimePeriod)
			if err != nil {
				return fmt.Errorf("expanding `time_period`: %+v", err)
			}

			parameters := budgets.Budget{
				Name: pointer.To(id.BudgetName),
				Properties: &budgets.BudgetProperties{
					Amount:        config.Amount,
					Category:      budgets.CategoryTypeCost,
					Filter:        expandConsumptionBudgetFilterFromModel(config.Filter),
					Notifications: expandConsumptionBudgetMgmtNotificationsFromModel(config.Notification),
					TimeGrain:     budgets.TimeGrainType(config.TimeGrain),
					TimePeriod:    *timePeriod,
				},
			}

			if config.Etag != "" {
				parameters.ETag = pointer.To(config.Etag)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagementGroupConsumptionBudget) CustomImporter() sdk.ResourceRunFunc {
	return r.base.importerFunc()
}

func expandConsumptionBudgetMgmtNotificationsFromModel(input []ConsumptionBudgetMgmtNotificationModel) *map[string]budgets.Notification {
	if len(input) == 0 {
		return nil
	}

	notifications := make(map[string]budgets.Notification)
	for _, n := range input {
		notification := budgets.Notification{
			Enabled:  n.Enabled,
			Operator: budgets.OperatorType(n.Operator),
			// nolint: gosec
			Threshold: float64(n.Threshold),
		}

		thresholdType := budgets.ThresholdType(n.ThresholdType)
		notification.ThresholdType = &thresholdType

		notification.ContactEmails = n.ContactEmails

		notificationKey := fmt.Sprintf("%s_%s_%f_Percent", string(thresholdType), string(notification.Operator), notification.Threshold)
		notifications[notificationKey] = notification
	}

	return &notifications
}

func flattenConsumptionBudgetMgmtNotificationsToModel(input *map[string]budgets.Notification) []ConsumptionBudgetMgmtNotificationModel {
	if input == nil {
		return []ConsumptionBudgetMgmtNotificationModel{}
	}

	result := make([]ConsumptionBudgetMgmtNotificationModel, 0)
	for _, n := range *input {
		thresholdType := string(budgets.ThresholdTypeActual)
		if v := n.ThresholdType; v != nil {
			thresholdType = string(*v)
		}

		result = append(result, ConsumptionBudgetMgmtNotificationModel{
			Enabled:       n.Enabled,
			Operator:      string(n.Operator),
			Threshold:     int64(n.Threshold),
			ThresholdType: thresholdType,
			ContactEmails: n.ContactEmails,
		})
	}

	return result
}
