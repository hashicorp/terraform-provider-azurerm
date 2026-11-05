// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Consumption.BudgetsClient

			var config SubscriptionConsumptionBudgetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := budgets.NewScopedBudgetID(config.SubscriptionId, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			timePeriod, err := helpers.ExpandConsumptionBudgetTimePeriodFromModel(config.TimePeriod)
			if err != nil {
				return fmt.Errorf("expanding `time_period`: %+v", err)
			}

			parameters := budgets.Budget{
				Name: pointer.To(id.BudgetName),
				Properties: &budgets.BudgetProperties{
					Amount:        config.Amount,
					Category:      budgets.CategoryTypeCost,
					Filter:        helpers.ExpandConsumptionBudgetFilterFromModel(config.Filter),
					Notifications: helpers.ExpandConsumptionBudgetNotificationsFromModel(config.Notification),
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
