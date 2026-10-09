// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/exports"
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
	Type            string                                  `tfschema:"type"`
	TimeFrame       string                                  `tfschema:"time_frame"`
	TimePeriod      CostManagementExportDataTimePeriodModel `tfschema:"time_period"`
	DataGranularity string                                  `tfschema:"data_granularity"`
	DataVersion     string                                  `tfschema:"data_version"`
	Filter          []CostManagementExportDataFilterModel   `tfschema:"filter"`
}

type CostManagementExportDataTimePeriodModel struct {
	From string `tfschema:"from"`
	To   string `tfschema:"to"`
}

type CostManagementExportDataFilterModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type costManagementExportBaseResource struct{}

func (br costManagementExportBaseResource) arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		"active": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "",
		},

		"location": commonschema.LocationOptional(),

		"identity": commonschema.SystemAssignedIdentityOptional(),

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
	}

	maps.Copy(output, fields)

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
	rootFolderPath := strings.TrimPrefix(loc.RootFolderPath, "/") // Ensure no leading slash

	return &exports.ExportDeliveryInfo{
		Destination: exports.ExportDeliveryDestination{
			ResourceId:     pointer.To(storageId.ID()),
			Container:      containerId.ContainerName,
			RootFolderPath: pointer.To(rootFolderPath),
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

	rootFolderPath := pointer.From(destination.RootFolderPath)

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

	timePeriod := &exports.ExportTimePeriod{
		From: opt.TimePeriod.From,
		To:   opt.TimePeriod.To,
	}

	filters := []exports.FilterItems{}
	for _, item := range opt.Filter {
		filters = append(filters, exports.FilterItems{
			Name:  pointer.ToEnum[exports.FilterItemNames](item.Name),
			Value: pointer.To(item.Value),
		})
	}

	dataset := &exports.ExportDataset{
		Granularity: pointer.ToEnum[exports.GranularityType](opt.DataGranularity),
		Configuration: &exports.ExportDatasetConfiguration{
			DataVersion: pointer.To(opt.DataVersion),
			Filters:     pointer.To(filters),
		},
	}

	return &exports.ExportDefinition{
		Type:       exports.ExportType(opt.Type),
		Timeframe:  exports.TimeframeType(opt.TimeFrame),
		TimePeriod: timePeriod,
		DataSet:    dataset,
	}
}

// flattenExportDataOptionsToModel converts the SDK type to the typed model
func flattenExportDataOptionsToModel(input exports.ExportDefinition) []CostManagementExportDataOptionsModel {
	queryType := ""
	if v := input.Type; v != "" {
		queryType = string(input.Type)
	}

	timePeriod := CostManagementExportDataTimePeriodModel{}
	if v := input.TimePeriod; v != nil {
		timePeriod.From = v.From
		timePeriod.To = v.To
	}

	dataVersion := ""
	dataGranularity := ""
	filters := []CostManagementExportDataFilterModel{}

	if input.DataSet != nil {
		dataGranularity = string(pointer.From(input.DataSet.Granularity))

		if c := input.DataSet.Configuration; c != nil {
			dataVersion = pointer.From(c.DataVersion)

			if v := c.Filters; v != nil {
				for _, item := range *v {
					filters = append(filters, CostManagementExportDataFilterModel{
						Name:  string(pointer.From(item.Name)),
						Value: pointer.From(item.Value),
					})
				}
			}
		}
	}

	return []CostManagementExportDataOptionsModel{
		{
			TimeFrame:       string(input.Timeframe),
			Type:            queryType,
			TimePeriod:      timePeriod,
			DataGranularity: dataGranularity,
			DataVersion:     dataVersion,
			Filter:          filters,
		},
	}
}
