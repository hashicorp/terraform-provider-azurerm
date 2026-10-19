// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription_cost_management_export

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient

			id, err := exports.ParseScopedExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SubscriptionCostManagementExportModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Update operation requires latest eTag to be set in the request.
			var opts exports.GetOperationOptions
			resp, err := client.Get(ctx, *id, opts)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if resp.Model.ETag == nil {
				return fmt.Errorf("retrieving %s: etag was nil", *id)
			}

			model := *resp.Model
			if model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", *id)
			}

			props := model.Properties

			if metadata.ResourceData.HasChange("active") {
				status := exports.StatusTypeActive
				if !config.Active {
					status = exports.StatusTypeInactive
				}
				if props.Schedule == nil {
					props.Schedule = &exports.ExportSchedule{}
				}
				props.Schedule.Status = &status
			}

			if metadata.ResourceData.HasChange("recurrence_type") {
				if props.Schedule == nil {
					props.Schedule = &exports.ExportSchedule{}
				}
				props.Schedule.Recurrence = pointer.ToEnum[exports.RecurrenceType](config.RecurrenceType)
			}

			if metadata.ResourceData.HasChanges("recurrence_period_start_date", "recurrence_period_end_date") {
				if props.Schedule == nil {
					props.Schedule = &exports.ExportSchedule{}
				}
				props.Schedule.RecurrencePeriod = &exports.ExportRecurrencePeriod{
					From: config.RecurrencePeriodStartDate,
					To:   pointer.To(config.RecurrencePeriodEndDate),
				}
			}

			if metadata.ResourceData.HasChange("export_data_storage_location") {
				deliveryInfo, err := helpers.ExpandExportDataStorageLocationFromModel(config.ExportDataStorageLocation)
				if err != nil {
					return fmt.Errorf("expanding `export_data_storage_location`: %+v", err)
				}
				if deliveryInfo == nil {
					return fmt.Errorf("`export_data_storage_location` was empty")
				}
				props.DeliveryInfo = *deliveryInfo
			}

			if metadata.ResourceData.HasChange("file_format") {
				props.Format = pointer.ToEnum[exports.FormatType](config.FileFormat)
			}

			if metadata.ResourceData.HasChange("export_data_options") {
				definition := helpers.ExpandExportDataOptionsFromModel(config.ExportDataOptions)
				if definition == nil {
					return fmt.Errorf("`export_data_options` was empty")
				}
				props.Definition = *definition
			}

			model.Properties = props

			if _, err = client.CreateOrUpdate(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
