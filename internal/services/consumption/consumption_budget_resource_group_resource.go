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
	validateResourceGroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ResourceGroupConsumptionBudget struct {
	base consumptionBudgetBaseResource
}

type ResourceGroupConsumptionBudgetModel struct {
	Name            string                               `tfschema:"name"`
	ResourceGroupId string                               `tfschema:"resource_group_id"`
	Etag            string                               `tfschema:"etag"`
	Amount          float64                              `tfschema:"amount"`
	TimeGrain       string                               `tfschema:"time_grain"`
	TimePeriod      []ConsumptionBudgetTimePeriodModel   `tfschema:"time_period"`
	Notification    []ConsumptionBudgetNotificationModel `tfschema:"notification"`
	Filter          []ConsumptionBudgetFilterModel       `tfschema:"filter"`
}

var (
	_ sdk.Resource                   = ResourceGroupConsumptionBudget{}
	_ sdk.ResourceWithCustomImporter = ResourceGroupConsumptionBudget{}
)

func (r ResourceGroupConsumptionBudget) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateResourceGroup.ResourceGroupID,
		},
	}
	return r.base.arguments(schema)
}

func (r ResourceGroupConsumptionBudget) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ResourceGroupConsumptionBudget) ModelObject() interface{} {
	return &ResourceGroupConsumptionBudgetModel{}
}

func (r ResourceGroupConsumptionBudget) ResourceType() string {
	return "azurerm_consumption_budget_resource_group"
}

func (r ResourceGroupConsumptionBudget) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return budgets.ValidateScopedBudgetID
}

func (r ResourceGroupConsumptionBudget) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			var config ResourceGroupConsumptionBudgetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := budgets.NewScopedBudgetID(config.ResourceGroupId, config.Name)

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

func (r ResourceGroupConsumptionBudget) Read() sdk.ResourceFunc {
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

			state := ResourceGroupConsumptionBudgetModel{
				Name:            id.BudgetName,
				ResourceGroupId: id.Scope,
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

func (r ResourceGroupConsumptionBudget) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ResourceGroupConsumptionBudget) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			id, err := budgets.ParseScopedBudgetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ResourceGroupConsumptionBudgetModel
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

func (r ResourceGroupConsumptionBudget) CustomImporter() sdk.ResourceRunFunc {
	return r.base.importerFunc()
}
