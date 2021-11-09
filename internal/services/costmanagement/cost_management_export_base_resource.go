package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2020-06-01/costmanagement"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
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

		"recurrence_period_start": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"recurrence_period_end": {
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
					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"container_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.ExportContainerName,
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

		"export_data_definition": {
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
							string(costmanagement.MonthToDate),
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

			if err := createOrUpdateCostManagementExport(ctx, client, metadata, id); err != nil {
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
				if !utils.ResponseWasNotFound(resp.Response){
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.Name)
			// lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			if schedule := resp.Schedule; schedule != nil {
				if recurrencePeriod := schedule.RecurrencePeriod; recurrencePeriod != nil {
					metadata.ResourceData.Set("recurrence_period_start", recurrencePeriod.From.Format(time.RFC3339))
					metadata.ResourceData.Set("recurrence_period_end", recurrencePeriod.To.Format(time.RFC3339))
				}
				status := false
				if schedule.Status == costmanagement.Active {
					status = true
				}
				metadata.ResourceData.Set("active", status)
				metadata.ResourceData.Set("recurrence_type", schedule.Recurrence)
			}
			if err := metadata.ResourceData.Set("export_data_storage_location", flattenExportDeliveryInfo(resp.DeliveryInfo)); err != nil {
				return fmt.Errorf("setting `export_data_storage_location`: %+v", err)
			}

			if err := metadata.ResourceData.Set("export_data_definition", flattenExportDefinition(resp.Definition)); err != nil {
				return fmt.Errorf("setting `export_data_definition`: %+v", err)
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

			if err := createOrUpdateCostManagementExport(ctx, client, metadata, *id); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func createOrUpdateCostManagementExport(ctx context.Context, client *costmanagement.ExportsClient, metadata sdk.ResourceMetaData, id parse.CostManagementExportId) error {
	from, _ := time.Parse(time.RFC3339, metadata.ResourceData.Get("recurrence_period_start").(string))
	to, _ := time.Parse(time.RFC3339, metadata.ResourceData.Get("recurrence_period_end").(string))

	status := costmanagement.Active
	if v := metadata.ResourceData.Get("active"); !v.(bool) {
		status = costmanagement.Inactive
	}

	props := costmanagement.Export{
		ExportProperties: &costmanagement.ExportProperties{
			Schedule: &costmanagement.ExportSchedule{
				Recurrence: costmanagement.RecurrenceType(metadata.ResourceData.Get("recurrence_type").(string)),
				RecurrencePeriod: &costmanagement.ExportRecurrencePeriod{
					From: &date.Time{Time: from},
					To: &date.Time{Time: to},
				},
				Status: status,
			},
			DeliveryInfo: expandExportDeliveryInfo(metadata.ResourceData.Get("export_data_storage_location").([]interface{})),
			Format:	costmanagement.Csv,
			Definition: expandExportDefinition(metadata.ResourceData.Get("export_data_definition").([]interface{})),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.Scope, id.Name, props); err != nil {
		return err
	}
	return nil
}

func expandExportDeliveryInfo(input []interface{}) *costmanagement.ExportDeliveryInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	deliveryInfo := &costmanagement.ExportDeliveryInfo{
		Destination: &costmanagement.ExportDeliveryDestination{
			ResourceID:     utils.String(attrs["storage_account_id"].(string)),
			Container:      utils.String(attrs["container_name"].(string)),
			RootFolderPath: utils.String(attrs["root_folder_path"].(string)),
		},
	}

	return deliveryInfo
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

func flattenExportDeliveryInfo(input *costmanagement.ExportDeliveryInfo) []interface{} {
	if input == nil || input.Destination == nil {
		return []interface{}{}
	}

	destination := input.Destination
	attrs := make(map[string]interface{})
	if resourceID := destination.ResourceID; resourceID != nil {
		attrs["storage_account_id"] = *resourceID
	}
	if containerName := destination.Container; containerName != nil {
		attrs["container_name"] = *containerName
	}
	if rootFolderPath := destination.RootFolderPath; rootFolderPath != nil {
		attrs["root_folder_path"] = *rootFolderPath
	}

	return []interface{}{attrs}
}

func flattenExportDefinition(input *costmanagement.ExportDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	attrs := make(map[string]interface{})
	if queryType := input.Type; queryType != "" {
		attrs["type"] = queryType
	}
	attrs["time_frame"] = string(input.Timeframe)

	return []interface{}{attrs}
}