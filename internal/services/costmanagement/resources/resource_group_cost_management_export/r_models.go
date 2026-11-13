// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_cost_management_export

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

type ResourceGroupCostManagementExportModel struct {
	Name                      string                                                 `tfschema:"name"`
	Description               string                                                 `tfschema:"description"`
	Location                  string                                                 `tfschema:"location"`
	ResourceGroupId           string                                                 `tfschema:"resource_group_id"`
	Active                    bool                                                   `tfschema:"active"`
	RecurrenceType            string                                                 `tfschema:"recurrence_type"`
	RecurrencePeriodStartDate string                                                 `tfschema:"recurrence_period_start_date"`
	RecurrencePeriodEndDate   string                                                 `tfschema:"recurrence_period_end_date"`
	FileFormat                string                                                 `tfschema:"file_format"`
	ExportDataStorageLocation []helpers.CostManagementExportDataStorageLocationModel `tfschema:"export_data_storage_location"`
	ExportDataOptions         []helpers.CostManagementExportDataOptionsModel         `tfschema:"export_data_options"`
	CompressionMode           string                                                 `tfschema:"compression_mode"`
	PartitionData             bool                                                   `tfschema:"partition_data"`
	DataOverwriteBehavior     string                                                 `tfschema:"data_overwrite_behavior"`
	Identity                  []identity.ModelSystemAssigned                         `tfschema:"identity"`
}

func (r Resource) ModelObject() any {
	return &ResourceGroupCostManagementExportModel{}
}
