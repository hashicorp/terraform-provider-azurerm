// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing_account_cost_management_export

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient

			id, err := exports.ParseScopedExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var opts exports.GetOperationOptions
			resp, err := client.Get(ctx, *id, opts)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BillingAccountCostManagementExportModel{
				Name:             id.ExportName,
				BillingAccountId: id.Scope,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if schedule := props.Schedule; schedule != nil {
						if recurrencePeriod := schedule.RecurrencePeriod; recurrencePeriod != nil {
							state.RecurrencePeriodStartDate = recurrencePeriod.From
							state.RecurrencePeriodEndDate = pointer.From(recurrencePeriod.To)
						}
						if schedule.Status != nil {
							state.Active = *schedule.Status == exports.StatusTypeActive
						}
						state.RecurrenceType = pointer.FromEnum(schedule.Recurrence)
					}

					storageLocation, err := helpers.FlattenExportDataStorageLocationToModel(props.DeliveryInfo)
					if err != nil {
						return fmt.Errorf("flattening `export_data_storage_location`: %+v", err)
					}
					state.ExportDataStorageLocation = storageLocation
					state.ExportDataOptions = helpers.FlattenExportDataOptionsToModel(props.Definition)
					state.FileFormat = pointer.FromEnum(props.Format)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
