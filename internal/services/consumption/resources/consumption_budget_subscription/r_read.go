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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
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
					state.TimePeriod = helpers.FlattenConsumptionBudgetTimePeriodToModel(&props.TimePeriod)
					state.Notification = helpers.FlattenConsumptionBudgetNotificationsToModel(props.Notifications)
					state.Filter = helpers.FlattenConsumptionBudgetFilterToModel(props.Filter)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
