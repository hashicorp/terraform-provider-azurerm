package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2020-06-01/costmanagement"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/parse"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
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
				string(costmanagement.RecurrenceTypeDaily),
				string(costmanagement.RecurrenceTypeWeekly),
				string(costmanagement.RecurrenceTypeMonthly),
				string(costmanagement.RecurrenceTypeAnnually),
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
						ValidateFunc: storageValidate.StorageContainerResourceManagerID,
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
							string(costmanagement.ExportTypeActualCost),
							string(costmanagement.ExportTypeAmortizedCost),
							string(costmanagement.ExportTypeUsage),
						}, false),
					},

					"time_frame": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(costmanagement.Custom),
							string(costmanagement.BillingMonthToDate),
							string(costmanagement.TheLastBillingMonth),
							string(costmanagement.TheLastMonth),
							string(costmanagement.WeekToDate),
							string(costmanagement.MonthToDate),
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
			id := parse.NewCostManagementExportId(metadata.ResourceData.Get(scopeFieldName).(string), metadata.ResourceData.Get("name").(string))
			existing, err := client.Get(ctx, id.Scope, id.Name, "")
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !utils.ResponseWasNotFound(existing.Response) {
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

			id, err := parse.CostManagementExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.Scope, id.Name, "")
			if err != nil {
				if !utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.Name)
			//lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			if schedule := resp.Schedule; schedule != nil {
				if recurrencePeriod := schedule.RecurrencePeriod; recurrencePeriod != nil {
					metadata.ResourceData.Set("recurrence_period_start_date", recurrencePeriod.From.Format(time.RFC3339))
					metadata.ResourceData.Set("recurrence_period_end_date", recurrencePeriod.To.Format(time.RFC3339))
				}

				status := schedule.Status == costmanagement.Active

				metadata.ResourceData.Set("active", status)
				metadata.ResourceData.Set("recurrence_type", schedule.Recurrence)
			}

			exportDeliveryInfo, err := flattenExportDataStorageLocation(resp.DeliveryInfo)
			if err != nil {
				return fmt.Errorf("flattening `export_data_storage_location`: %+v", err)
			}

			if err := metadata.ResourceData.Set("export_data_storage_location", exportDeliveryInfo); err != nil {
				return fmt.Errorf("setting `export_data_storage_location`: %+v", err)
			}

			if err := metadata.ResourceData.Set("export_data_options", flattenExportDefinition(resp.Definition)); err != nil {
				return fmt.Errorf("setting `export_data_options`: %+v", err)
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

			id, err := parse.CostManagementExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, id.Scope, id.Name); err != nil {
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

			id, err := parse.CostManagementExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Update operation requires latest eTag to be set in the request.
			resp, err := client.Get(ctx, id.Scope, id.Name, "")
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}
			if resp.ETag == nil {
				return fmt.Errorf("add %s: etag was nil", *id)
			}
			if err := createOrUpdateCostManagementExport(ctx, client, metadata, *id, resp.ETag); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func createOrUpdateCostManagementExport(ctx context.Context, client *costmanagement.ExportsClient, metadata sdk.ResourceMetaData, id parse.CostManagementExportId, etag *string) error {
	from, _ := time.Parse(time.RFC3339, metadata.ResourceData.Get("recurrence_period_start_date").(string))
	to, _ := time.Parse(time.RFC3339, metadata.ResourceData.Get("recurrence_period_end_date").(string))

	status := costmanagement.Active
	if v := metadata.ResourceData.Get("active"); !v.(bool) {
		status = costmanagement.Inactive
	}

	deliveryInfo, err := expandExportDataStorageLocation(metadata.ResourceData.Get("export_data_storage_location").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `export_data_storage_location`: %+v", err)
	}

	props := costmanagement.Export{
		ETag: etag,
		ExportProperties: &costmanagement.ExportProperties{
			Schedule: &costmanagement.ExportSchedule{
				Recurrence: costmanagement.RecurrenceType(metadata.ResourceData.Get("recurrence_type").(string)),
				RecurrencePeriod: &costmanagement.ExportRecurrencePeriod{
					From: &date.Time{Time: from},
					To:   &date.Time{Time: to},
				},
				Status: status,
			},
			DeliveryInfo: deliveryInfo,
			Format:       costmanagement.Csv,
			Definition:   expandExportDefinition(metadata.ResourceData.Get("export_data_options").([]interface{})),
		},
	}

	_, err = client.CreateOrUpdate(ctx, id.Scope, id.Name, props)

	return err
}

func expandExportDataStorageLocation(input []interface{}) (*costmanagement.ExportDeliveryInfo, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	attrs := input[0].(map[string]interface{})

	containerId, err := storageParse.StorageContainerResourceManagerID(attrs["container_id"].(string))
	if err != nil {
		return nil, err
	}

	storageId := storageParse.NewStorageAccountID(containerId.SubscriptionId, containerId.ResourceGroup, containerId.StorageAccountName)

	deliveryInfo := &costmanagement.ExportDeliveryInfo{
		Destination: &costmanagement.ExportDeliveryDestination{
			ResourceID:     utils.String(storageId.ID()),
			Container:      utils.String(containerId.ContainerName),
			RootFolderPath: utils.String(attrs["root_folder_path"].(string)),
		},
	}

	return deliveryInfo, nil
}

func expandExportDefinition(input []interface{}) *costmanagement.ExportDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	definitionInfo := &costmanagement.ExportDefinition{
		Type:      costmanagement.ExportType(attrs["type"].(string)),
		Timeframe: costmanagement.TimeframeType(attrs["time_frame"].(string)),
	}

	return definitionInfo
}

func flattenExportDataStorageLocation(input *costmanagement.ExportDeliveryInfo) ([]interface{}, error) {
	if input == nil || input.Destination == nil {
		return []interface{}{}, nil
	}

	destination := input.Destination
	var err error
	var storageAccountId *storageParse.StorageAccountId

	if v := destination.ResourceID; v != nil {
		storageAccountId, err = storageParse.StorageAccountID(*v)
		if err != nil {
			return nil, err
		}
	}

	containerId := ""
	if v := destination.Container; v != nil && storageAccountId != nil {
		containerId = storageParse.NewStorageContainerResourceManagerID(storageAccountId.SubscriptionId, storageAccountId.ResourceGroup, storageAccountId.Name, "default", *v).ID()
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

func flattenExportDefinition(input *costmanagement.ExportDefinition) []interface{} {
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
