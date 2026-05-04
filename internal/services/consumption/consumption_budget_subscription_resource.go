// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SubscriptionConsumptionBudget struct {
	base consumptionBudgetBaseResource
}

type SubscriptionConsumptionBudgetModel struct {
	Name           string                               `tfschema:"name"`
	SubscriptionId string                               `tfschema:"subscription_id"`
	Etag           string                               `tfschema:"etag"`
	Amount         float64                              `tfschema:"amount"`
	TimeGrain      string                               `tfschema:"time_grain"`
	TimePeriod     []ConsumptionBudgetTimePeriodModel   `tfschema:"time_period"`
	Notification   []ConsumptionBudgetNotificationModel `tfschema:"notification"`
	Filter         []ConsumptionBudgetFilterModel       `tfschema:"filter"`
}

var (
	_ sdk.Resource                   = SubscriptionConsumptionBudget{}
	_ sdk.ResourceWithCustomImporter = SubscriptionConsumptionBudget{}
	_ sdk.ResourceWithStateMigration = SubscriptionConsumptionBudget{}
)

func (r SubscriptionConsumptionBudget) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConsumptionBudgetName(),
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},
	}
	return r.base.arguments(schema)
}

func (r SubscriptionConsumptionBudget) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r SubscriptionConsumptionBudget) ModelObject() interface{} {
	return &SubscriptionConsumptionBudgetModel{}
}

func (r SubscriptionConsumptionBudget) ResourceType() string {
	return "azurerm_consumption_budget_subscription"
}

func (r SubscriptionConsumptionBudget) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return budgets.ValidateScopedBudgetID
}

func (r SubscriptionConsumptionBudget) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			var config SubscriptionConsumptionBudgetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := budgets.NewScopedBudgetID(config.SubscriptionId, config.Name)

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
					Notifications: expandConsumptionBudgetNotificationsFromModel(config.Notification),
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

func (r SubscriptionConsumptionBudget) Read() sdk.ResourceFunc {
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

			state := SubscriptionConsumptionBudgetModel{
				Name:           id.BudgetName,
				SubscriptionId: id.Scope,
			}

			if model := resp.Model; model != nil {
				state.Etag = pointer.From(model.ETag)

				if props := model.Properties; props != nil {
					state.Amount = props.Amount
					state.TimeGrain = string(props.TimeGrain)
					state.TimePeriod = flattenConsumptionBudgetTimePeriodToModel(&props.TimePeriod)
					state.Notification = flattenConsumptionBudgetNotificationsToModel(props.Notifications)
					state.Filter = flattenConsumptionBudgetFilterToModel(props.Filter)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SubscriptionConsumptionBudget) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r SubscriptionConsumptionBudget) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			id, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SubscriptionConsumptionBudgetModel
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
					Notifications: expandConsumptionBudgetNotificationsFromModel(config.Notification),
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

func (r SubscriptionConsumptionBudget) CustomImporter() sdk.ResourceRunFunc {
	return r.base.importerFunc()
}

func (r SubscriptionConsumptionBudget) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 2,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.SubscriptionConsumptionBudgetV0ToV1{},
			1: migration.SubscriptionConsumptionBudgetV1ToV2{},
		},
	}
}
