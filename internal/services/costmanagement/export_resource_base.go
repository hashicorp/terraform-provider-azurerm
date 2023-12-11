// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2021-10-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type costManagementExportBaseResource struct{}

func (br costManagementExportBaseResource) arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		"active": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"recurrence_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(exports.RecurrenceTypeDaily),
				string(exports.RecurrenceTypeWeekly),
				string(exports.RecurrenceTypeMonthly),
				string(exports.RecurrenceTypeAnnually),
			}, false),
		},

		"recurrence_period_start_date": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"recurrence_period_end_date": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"export_data_storage_location": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageContainerID,
					},
					"root_folder_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"export_data_options": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(exports.ExportTypeActualCost),
							string(exports.ExportTypeAmortizedCost),
							string(exports.ExportTypeUsage),
						}, false),
					},

					"time_frame": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(exports.TimeframeTypeCustom),
							string(exports.TimeframeTypeBillingMonthToDate),
							string(exports.TimeframeTypeTheLastBillingMonth),
							string(exports.TimeframeTypeTheLastMonth),
							string(exports.TimeframeTypeWeekToDate),
							string(exports.TimeframeTypeMonthToDate),
							// TODO Use value from SDK after https://github.com/Azure/azure-rest-api-specs/issues/23707 is fixed
							"TheLast7Days",
						}, false),
					},
				},
			},
		},
	}

	for k, v := range fields {
		output[k] = v
	}

	return output
}

func (br costManagementExportBaseResource) attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br costManagementExportBaseResource) createFunc(resourceName, scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient
			id := exports.NewScopedExportID(metadata.ResourceData.Get(scopeFieldName).(string), metadata.ResourceData.Get("name").(string))
			var opts exports.GetOperationOptions
			existing, err := client.Get(ctx, id, opts)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			if err := createOrUpdateCostManagementExport(ctx, client, metadata, id, nil); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (br costManagementExportBaseResource) readFunc(scopeFieldName string) sdk.ResourceFunc {
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
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.ExportName)
			//lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if schedule := props.Schedule; schedule != nil {
						if recurrencePeriod := schedule.RecurrencePeriod; recurrencePeriod != nil {
							metadata.ResourceData.Set("recurrence_period_start_date", recurrencePeriod.From)
							metadata.ResourceData.Set("recurrence_period_end_date", recurrencePeriod.To)
						}
						status := *schedule.Status == exports.StatusTypeActive

						metadata.ResourceData.Set("active", status)
						metadata.ResourceData.Set("recurrence_type", string(pointer.From(schedule.Recurrence)))
					}

					exportDeliveryInfo, err := flattenExportDataStorageLocation(&props.DeliveryInfo)
					if err != nil {
						return fmt.Errorf("flattening `export_data_storage_location`: %+v", err)
					}
					if err := metadata.ResourceData.Set("export_data_storage_location", exportDeliveryInfo); err != nil {
						return fmt.Errorf("setting `export_data_storage_location`: %+v", err)
					}
					if err := metadata.ResourceData.Set("export_data_options", flattenExportDefinition(&props.Definition)); err != nil {
						return fmt.Errorf("setting `export_data_options`: %+v", err)
					}
				}
			}

			return nil
		},
	}
}

func (br costManagementExportBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient

			id, err := exports.ParseScopedExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (br costManagementExportBaseResource) updateFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient

			id, err := exports.ParseScopedExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Update operation requires latest eTag to be set in the request.
			var opts exports.GetOperationOptions
			resp, err := client.Get(ctx, *id, opts)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				if model.ETag == nil {
					return fmt.Errorf("add %s: etag was nil", *id)
				}
			}

			if err := createOrUpdateCostManagementExport(ctx, client, metadata, *id, resp.Model.ETag); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func createOrUpdateCostManagementExport(ctx context.Context, client *exports.ExportsClient, metadata sdk.ResourceMetaData, id exports.ScopedExportId, etag *string) error {
	status := exports.StatusTypeActive
	if v := metadata.ResourceData.Get("active"); !v.(bool) {
		status = exports.StatusTypeInactive
	}

	deliveryInfo, err := expandExportDataStorageLocation(metadata.ResourceData.Get("export_data_storage_location").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `export_data_storage_location`: %+v", err)
	}

	format := exports.FormatTypeCsv
	recurrenceType := exports.RecurrenceType(metadata.ResourceData.Get("recurrence_type").(string))
	props := exports.Export{
		ETag: etag,
		Properties: &exports.ExportProperties{
			Schedule: &exports.ExportSchedule{
				Recurrence: &recurrenceType,
				RecurrencePeriod: &exports.ExportRecurrencePeriod{
					From: metadata.ResourceData.Get("recurrence_period_start_date").(string),
					To:   utils.String(metadata.ResourceData.Get("recurrence_period_end_date").(string)),
				},
				Status: &status,
			},
			DeliveryInfo: *deliveryInfo,
			Format:       &format,
			Definition:   *expandExportDefinition(metadata.ResourceData.Get("export_data_options").([]interface{})),
		},
	}

	_, err = client.CreateOrUpdate(ctx, id, props)

	return err
}

func expandExportDataStorageLocation(input []interface{}) (*exports.ExportDeliveryInfo, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	attrs := input[0].(map[string]interface{})

	containerId, err := commonids.ParseStorageContainerID(attrs["container_id"].(string))
	if err != nil {
		return nil, err
	}

	storageId := commonids.NewStorageAccountID(containerId.SubscriptionId, containerId.ResourceGroupName, containerId.StorageAccountName)

	deliveryInfo := &exports.ExportDeliveryInfo{
		Destination: exports.ExportDeliveryDestination{
			ResourceId:     utils.String(storageId.ID()),
			Container:      containerId.ContainerName,
			RootFolderPath: utils.String(attrs["root_folder_path"].(string)),
		},
	}

	return deliveryInfo, nil
}

func expandExportDefinition(input []interface{}) *exports.ExportDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	definitionInfo := &exports.ExportDefinition{
		Type:      exports.ExportType(attrs["type"].(string)),
		Timeframe: exports.TimeframeType(attrs["time_frame"].(string)),
	}

	return definitionInfo
}

func flattenExportDataStorageLocation(input *exports.ExportDeliveryInfo) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	destination := input.Destination
	var err error
	var storageAccountId *commonids.StorageAccountId

	if v := destination.ResourceId; v != nil {
		storageAccountId, err = commonids.ParseStorageAccountIDInsensitively(*v)
		if err != nil {
			return nil, err
		}
	}

	containerId := ""
	if v := destination.Container; v != "" && storageAccountId != nil {
		containerId = commonids.NewStorageContainerID(storageAccountId.SubscriptionId, storageAccountId.ResourceGroupName, storageAccountId.StorageAccountName, v).ID()
	}

	rootFolderPath := ""
	if v := destination.RootFolderPath; v != nil {
		rootFolderPath = *v
	}

	return []interface{}{
		map[string]interface{}{
			"container_id":     containerId,
			"root_folder_path": rootFolderPath,
		},
	}, nil
}

func flattenExportDefinition(input *exports.ExportDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	queryType := ""
	if v := input.Type; v != "" {
		queryType = string(input.Type)
	}

	return []interface{}{
		map[string]interface{}{
			"time_frame": string(input.Timeframe),
			"type":       queryType,
		},
	}
}
