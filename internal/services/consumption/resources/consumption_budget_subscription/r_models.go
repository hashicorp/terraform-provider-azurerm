// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_subscription

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/helpers"
)

type SubscriptionConsumptionBudgetModel struct {
	Name           string                                       `tfschema:"name"`
	SubscriptionId string                                       `tfschema:"subscription_id"`
	Etag           string                                       `tfschema:"etag"`
	Amount         float64                                      `tfschema:"amount"`
	TimeGrain      string                                       `tfschema:"time_grain"`
	TimePeriod     []helpers.ConsumptionBudgetTimePeriodModel   `tfschema:"time_period"`
	Notification   []helpers.ConsumptionBudgetNotificationModel `tfschema:"notification"`
	Filter         []helpers.ConsumptionBudgetFilterModel       `tfschema:"filter"`
}

func (r Resource) ModelObject() any {
	return &SubscriptionConsumptionBudgetModel{}
}
