// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(exports.PossibleValuesForRecurrenceType(), false),
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

		"file_format": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(exports.FormatTypeCsv),
			ValidateFunc: validation.StringInSlice(exports.PossibleValuesForFormatType(), false),
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: func(val interface{}, key string) ([]string, []error) {
							warnings, errors := validation.StringIsNotEmpty(val, key)

							// Since API 2025-03-01, root_folder_path cannot start with a slash.
							if v := val.(string); strings.HasPrefix(v, "/") {
								warnings = append(warnings, fmt.Sprintf("%q should not start with a slash", key))
							}
							return warnings, errors
						},
						DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
							return strings.TrimPrefix(new, "/") == strings.TrimPrefix(old, "/")
						},
						DiffSuppressOnRefresh: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(exports.PossibleValuesForExportType(), false),
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

					"time_period": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"from": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},

								"to": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
							},
						},
					},

					"data_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
							// Ignore changes to data_version if it is not set
							return new == ""
						},
						DiffSuppressOnRefresh: true,
					},

					"data_granularity": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(exports.GranularityTypeDaily),
						ValidateFunc: validation.StringInSlice(exports.PossibleValuesForGranularityType(), false),
					},

					"filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(exports.PossibleValuesForFilterItemNames(), false),
								},
								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										// for ReservationScope
										"Single",
										"Shared",

										// for LookBackPeriod
										"Last7Days",
										"Last30Days",
										"Last60Days",

										// for ResourceType
										"VirtualMachines",
										"SQLDatabases",
										"PostgreSQL",
										"ManagedDisk",
										"MySQL",
										"RedHat",
										"MariaDB",
										"RedisCache",
										"CosmosDB",
										"SqlDataWarehouse",
										"SUSELinux",
										"AppService",
										"BlockBlob",
										"AzureDataExplorer",
										"VMwareCloudSimple",
									}, false),
								},
							},
						},
					},
				},
			},
		},

		"compression_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(exports.CompressionModeTypeNone),
			ValidateFunc: validation.StringInSlice(exports.PossibleValuesForCompressionModeType(), false),
		},

		"partition_data": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"data_overwrite_behavior": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(exports.DataOverwriteBehaviorTypeCreateNewReport),
			ValidateFunc: validation.StringInSlice(exports.PossibleValuesForDataOverwriteBehaviorType(), false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "",
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
			// lintignore:R001
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
					metadata.ResourceData.Set("description", pointer.From(props.ExportDescription))
					metadata.ResourceData.Set("file_format", string(pointer.From(props.Format)))
					metadata.ResourceData.Set("compression_mode", string(pointer.From(props.CompressionMode)))
					metadata.ResourceData.Set("partition_data", pointer.From(props.PartitionData))
					metadata.ResourceData.Set("data_overwrite_behavior", string(pointer.From(props.DataOverwriteBehavior)))
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

	// PartitionData cannot be set to false otherwise API will return an error.
	var partitionData *bool
	if v, ok := metadata.ResourceData.Get("partition_data").(bool); ok && v {
		partitionData = utils.Bool(true)
	}

	format := exports.FormatType(metadata.ResourceData.Get("file_format").(string))
	compressionMode := exports.CompressionModeType(metadata.ResourceData.Get("compression_mode").(string))
	dataOverwriteBehavior := exports.DataOverwriteBehaviorType(metadata.ResourceData.Get("data_overwrite_behavior").(string))

	recurrenceType := exports.RecurrenceType(metadata.ResourceData.Get("recurrence_type").(string))
	props := exports.Export{
		ETag: etag,
		Properties: &exports.ExportProperties{
			Schedule: &exports.ExportSchedule{
				Recurrence: &recurrenceType,
				RecurrencePeriod: &exports.ExportRecurrencePeriod{
					From: metadata.ResourceData.Get("recurrence_period_start_date").(string),
					To:   pointer.To(metadata.ResourceData.Get("recurrence_period_end_date").(string)),
				},
				Status: &status,
			},
			DeliveryInfo:          *deliveryInfo,
			ExportDescription:     utils.String(metadata.ResourceData.Get("description").(string)),
			Format:                &format,
			CompressionMode:       &compressionMode,
			PartitionData:         partitionData,
			Definition:            *expandExportDefinition(metadata.ResourceData.Get("export_data_options").([]interface{})),
			DataOverwriteBehavior: &dataOverwriteBehavior,
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
	rootFolderPath := strings.TrimPrefix(attrs["root_folder_path"].(string), "/") // Ensure no leading slash

	deliveryInfo := &exports.ExportDeliveryInfo{
		Destination: exports.ExportDeliveryDestination{
			ResourceId:     pointer.To(storageId.ID()),
			Container:      containerId.ContainerName,
			RootFolderPath: pointer.To(rootFolderPath),
		},
	}

	return deliveryInfo, nil
}

func expandExportDefinition(input []interface{}) *exports.ExportDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	granularity := exports.GranularityType(attrs["data_granularity"].(string))

	var timePeriod *exports.ExportTimePeriod
	if v, ok := attrs["time_period"].([]map[string]string); ok && len(v) > 0 {
		timePeriod = &exports.ExportTimePeriod{
			From: v[0]["from"],
			To:   v[0]["to"],
		}
	}

	definitionInfo := &exports.ExportDefinition{
		Type:       exports.ExportType(attrs["type"].(string)),
		Timeframe:  exports.TimeframeType(attrs["time_frame"].(string)),
		TimePeriod: timePeriod,
		DataSet: &exports.ExportDataset{
			Granularity:   &granularity,
			Configuration: expandExportDatasetConfiguration(attrs),
		},
	}

	return definitionInfo
}

func expandExportDatasetConfiguration(attrs map[string]interface{}) *exports.ExportDatasetConfiguration {
	var filters *[]exports.FilterItems
	if filterInput, ok := attrs["filter"].([]map[string]string); ok && len(filterInput) > 0 {
		filters = &[]exports.FilterItems{}
		for _, item := range filterInput {
			name := exports.FilterItemNames(item["name"])
			*filters = append(*filters, exports.FilterItems{
				Name:  &name,
				Value: utils.String(item["value"]),
			})
		}
	}

	var dataVersion *string
	if v, ok := attrs["data_version"].(string); ok && v != "" {
		dataVersion = utils.String(v)
	}

	return &exports.ExportDatasetConfiguration{
		DataVersion: dataVersion,
		Filters:     filters,
	}
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

	var timePeriod []map[string]string
	if input.TimePeriod != nil {
		timePeriod = []map[string]string{
			{
				"from": input.TimePeriod.From,
				"to":   input.TimePeriod.To,
			},
		}
	}

	dataVersion := ""
	dataGranularity := ""
	filters := &[]map[string]string{}

	if input.DataSet != nil {
		dataGranularity = string(pointer.From(input.DataSet.Granularity))

		if c := input.DataSet.Configuration; c != nil {
			dataVersion = pointer.From(c.DataVersion)

			if c.Filters != nil && len(*c.Filters) > 0 {
				for _, item := range *c.Filters {
					*filters = append(*filters, map[string]string{
						"name":  string(pointer.From(item.Name)),
						"value": pointer.From(item.Value),
					})
				}
			}
		}
	}

	return []interface{}{
		map[string]interface{}{
			"time_frame":       string(input.Timeframe),
			"type":             queryType,
			"time_period":      timePeriod,
			"data_version":     dataVersion,
			"data_granularity": dataGranularity,
			"filter":           filters,
		},
	}
}
