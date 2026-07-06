// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Shared nested model structs for cost management export resources

type CostManagementExportDataStorageLocationModel struct {
	ContainerId    string `tfschema:"container_id"`
	RootFolderPath string `tfschema:"root_folder_path"`
}

type CostManagementExportDataOptionsModel struct {
	Type      string `tfschema:"type"`
	TimeFrame string `tfschema:"time_frame"`
}

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

		"file_format": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(exports.FormatTypeCsv),
			ValidateFunc: validation.StringInSlice([]string{
				string(exports.FormatTypeCsv),
				// TODO add support for Parquet once added to the SDK
			}, false),
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

// expandExportDataStorageLocationFromModel converts the typed model to the SDK type
func expandExportDataStorageLocationFromModel(input []CostManagementExportDataStorageLocationModel) (*exports.ExportDeliveryInfo, error) {
	if len(input) == 0 {
		return nil, nil
	}

	loc := input[0]

	containerId, err := commonids.ParseStorageContainerID(loc.ContainerId)
	if err != nil {
		return nil, err
	}

	storageId := commonids.NewStorageAccountID(containerId.SubscriptionId, containerId.ResourceGroupName, containerId.StorageAccountName)

	return &exports.ExportDeliveryInfo{
		Destination: exports.ExportDeliveryDestination{
			ResourceId:     pointer.To(storageId.ID()),
			Container:      containerId.ContainerName,
			RootFolderPath: pointer.To(loc.RootFolderPath),
		},
	}, nil
}

// flattenExportDataStorageLocationToModel converts the SDK type to the typed model
func flattenExportDataStorageLocationToModel(input exports.ExportDeliveryInfo) ([]CostManagementExportDataStorageLocationModel, error) {
	destination := input.Destination

	var storageAccountId *commonids.StorageAccountId

	if v := destination.ResourceId; v != nil {
		var err error
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

	return []CostManagementExportDataStorageLocationModel{
		{
			ContainerId:    containerId,
			RootFolderPath: rootFolderPath,
		},
	}, nil
}

// expandExportDataOptionsFromModel converts the typed model to the SDK type
func expandExportDataOptionsFromModel(input []CostManagementExportDataOptionsModel) *exports.ExportDefinition {
	if len(input) == 0 {
		return nil
	}

	opt := input[0]
	return &exports.ExportDefinition{
		Type:      exports.ExportType(opt.Type),
		Timeframe: exports.TimeframeType(opt.TimeFrame),
	}
}

// flattenExportDataOptionsToModel converts the SDK type to the typed model
func flattenExportDataOptionsToModel(input exports.ExportDefinition) []CostManagementExportDataOptionsModel {
	queryType := ""
	if v := input.Type; v != "" {
		queryType = string(input.Type)
	}

	return []CostManagementExportDataOptionsModel{
		{
			TimeFrame: string(input.Timeframe),
			Type:      queryType,
		},
	}
}
