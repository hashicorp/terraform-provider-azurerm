// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_management_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
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
