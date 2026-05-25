// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SubscriptionCostManagementExportResource struct {
	base costManagementExportBaseResource
}

type SubscriptionCostManagementExportModel struct {
	Name                      string                                         `tfschema:"name"`
	SubscriptionId            string                                         `tfschema:"subscription_id"`
	Active                    bool                                           `tfschema:"active"`
	RecurrenceType            string                                         `tfschema:"recurrence_type"`
	RecurrencePeriodStartDate string                                         `tfschema:"recurrence_period_start_date"`
	RecurrencePeriodEndDate   string                                         `tfschema:"recurrence_period_end_date"`
	FileFormat                string                                         `tfschema:"file_format"`
	ExportDataStorageLocation []CostManagementExportDataStorageLocationModel `tfschema:"export_data_storage_location"`
	ExportDataOptions         []CostManagementExportDataOptionsModel         `tfschema:"export_data_options"`
}

var _ sdk.Resource = SubscriptionCostManagementExportResource{}

func (r SubscriptionCostManagementExportResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},
	}
	return r.base.arguments(schema)
}

func (r SubscriptionCostManagementExportResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r SubscriptionCostManagementExportResource) ModelObject() interface{} {
	return &SubscriptionCostManagementExportModel{}
}

func (r SubscriptionCostManagementExportResource) ResourceType() string {
	return "azurerm_subscription_cost_management_export"
}

func (r SubscriptionCostManagementExportResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SubscriptionCostManagementExportID
}

func (r SubscriptionCostManagementExportResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient

			var config SubscriptionCostManagementExportModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := exports.NewScopedExportID(config.SubscriptionId, config.Name)

			var opts exports.GetOperationOptions
			existing, err := client.Get(ctx, id, opts)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			deliveryInfo, err := expandExportDataStorageLocationFromModel(config.ExportDataStorageLocation)
			if err != nil {
				return fmt.Errorf("expanding `export_data_storage_location`: %+v", err)
			}
			if deliveryInfo == nil {
				return fmt.Errorf("`export_data_storage_location` was empty")
			}

			definition := expandExportDataOptionsFromModel(config.ExportDataOptions)
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

func (r SubscriptionCostManagementExportResource) Read() sdk.ResourceFunc {
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

			state := SubscriptionCostManagementExportModel{
				Name:           id.ExportName,
				SubscriptionId: id.Scope,
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
						state.RecurrenceType = string(pointer.From(schedule.Recurrence))
					}

					storageLocation, err := flattenExportDataStorageLocationToModel(props.DeliveryInfo)
					if err != nil {
						return fmt.Errorf("flattening `export_data_storage_location`: %+v", err)
					}
					state.ExportDataStorageLocation = storageLocation
					state.ExportDataOptions = flattenExportDataOptionsToModel(props.Definition)
					state.FileFormat = string(pointer.From(props.Format))
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SubscriptionCostManagementExportResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r SubscriptionCostManagementExportResource) Update() sdk.ResourceFunc {
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
				deliveryInfo, err := expandExportDataStorageLocationFromModel(config.ExportDataStorageLocation)
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
				definition := expandExportDataOptionsFromModel(config.ExportDataOptions)
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
