// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription_cost_management_export

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

type SubscriptionCostManagementExportModel struct {
	Name                      string                                                 `tfschema:"name"`
	SubscriptionId            string                                                 `tfschema:"subscription_id"`
	Active                    bool                                                   `tfschema:"active"`
	RecurrenceType            string                                                 `tfschema:"recurrence_type"`
	RecurrencePeriodStartDate string                                                 `tfschema:"recurrence_period_start_date"`
	RecurrencePeriodEndDate   string                                                 `tfschema:"recurrence_period_end_date"`
	FileFormat                string                                                 `tfschema:"file_format"`
	ExportDataStorageLocation []helpers.CostManagementExportDataStorageLocationModel `tfschema:"export_data_storage_location"`
	ExportDataOptions         []helpers.CostManagementExportDataOptionsModel         `tfschema:"export_data_options"`
}

func (r Resource) ModelObject() any {
	return &SubscriptionCostManagementExportModel{}
}
