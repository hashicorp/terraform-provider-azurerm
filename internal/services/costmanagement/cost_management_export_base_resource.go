package costmanagement

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2019-10-01/costmanagement"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

type costManagementExportBaseResource struct{}

func (br costManagementExportBaseResource) createFunc(resourceName, scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient
			id := parse.NewCostManagementExportId(metadata.ResourceData.Get(scopeFieldName).(string), metadata.ResourceData.Get("name").(string))
			existing, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			from, _ := time.Parse(time.RFC3339, metadata.ResourceData.Get("recurrence_period_start").(string))
			to, _ := time.Parse(time.RFC3339, metadata.ResourceData.Get("recurrence_period_end").(string))

			status := costmanagement.Active
			if v := metadata.ResourceData.Get("actuve"); !v.(bool) {
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
					Definition: expandExportQuery(metadata.ResourceData.Get("query").([]interface{})),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id.Scope, id.Name, props); err != nil {
				return fmt.Errorf("creating or updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (br costManagementExportBaseResource) readFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ExportClient

			id, err := parse.CostManagementExportID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.Scope, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
		},
	}
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

func expandExportQuery(input []interface{}) *costmanagement.QueryDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	definitionInfo := &costmanagement.QueryDefinition{
		Type:      utils.String(attrs["type"].(string)),
		Timeframe: costmanagement.TimeframeType(attrs["time_frame"].(string)),
	}

	return definitionInfo
}