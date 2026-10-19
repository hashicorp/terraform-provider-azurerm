// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_management_group

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/helpers"
)

type ManagementGroupConsumptionBudgetModel struct {
	Name              string                                     `tfschema:"name"`
	ManagementGroupId string                                     `tfschema:"management_group_id"`
	Etag              string                                     `tfschema:"etag"`
	Amount            float64                                    `tfschema:"amount"`
	TimeGrain         string                                     `tfschema:"time_grain"`
	TimePeriod        []helpers.ConsumptionBudgetTimePeriodModel `tfschema:"time_period"`
	Notification      []ConsumptionBudgetMgmtNotificationModel   `tfschema:"notification"`
	Filter            []helpers.ConsumptionBudgetFilterModel     `tfschema:"filter"`
}

// Management group notification model (no contact_groups/contact_roles)
type ConsumptionBudgetMgmtNotificationModel struct {
	Enabled       bool     `tfschema:"enabled"`
	Threshold     int64    `tfschema:"threshold"`
	ThresholdType string   `tfschema:"threshold_type"`
	Operator      string   `tfschema:"operator"`
	ContactEmails []string `tfschema:"contact_emails"`
}

func (r Resource) ModelObject() any {
	return &ManagementGroupConsumptionBudgetModel{}
}
