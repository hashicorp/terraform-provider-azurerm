// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription_cost_management_view

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

type SubscriptionCostManagementViewModel struct {
	Name        string                                   `tfschema:"name"`
	SubId       string                                   `tfschema:"subscription_id"`
	DisplayName string                                   `tfschema:"display_name"`
	ChartType   string                                   `tfschema:"chart_type"`
	Accumulated bool                                     `tfschema:"accumulated"`
	ReportType  string                                   `tfschema:"report_type"`
	Timeframe   string                                   `tfschema:"timeframe"`
	Dataset     []helpers.CostManagementViewDatasetModel `tfschema:"dataset"`
	Kpi         []helpers.CostManagementViewKpiModel     `tfschema:"kpi"`
	Pivot       []helpers.CostManagementViewPivotModel   `tfschema:"pivot"`
}

func (r Resource) ModelObject() any {
	return &SubscriptionCostManagementViewModel{}
}
