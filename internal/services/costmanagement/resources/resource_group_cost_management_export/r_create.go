// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_group_cost_management_export

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient

			var config ResourceGroupCostManagementExportModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := exports.NewScopedExportID(config.ResourceGroupId, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, exports.DefaultGetOperationOptions())
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			deliveryInfo, err := helpers.ExpandExportDataStorageLocationFromModel(config.ExportDataStorageLocation)
			if err != nil {
				return fmt.Errorf("expanding `export_data_storage_location`: %+v", err)
			}
			if deliveryInfo == nil {
				return fmt.Errorf("`export_data_storage_location` was empty")
			}

			definition := helpers.ExpandExportDataOptionsFromModel(config.ExportDataOptions)
			if definition == nil {
				return fmt.Errorf("`export_data_options` was empty")
			}

			status := exports.StatusTypeActive
			if !config.Active {
				status = exports.StatusTypeInactive
			}

			props := exports.Export{
				Properties: &exports.ExportProperties{
					Schedule: &exports.ExportSchedule{
						Recurrence: pointer.ToEnum[exports.RecurrenceType](config.RecurrenceType),
						RecurrencePeriod: &exports.ExportRecurrencePeriod{
							From: config.RecurrencePeriodStartDate,
							To:   pointer.To(config.RecurrencePeriodEndDate),
						},
						Status: &status,
					},
					DeliveryInfo: *deliveryInfo,
					Format:       pointer.ToEnum[exports.FormatType](config.FileFormat),
					Definition:   *definition,
				},
			}

			if _, err = client.CreateOrUpdate(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
