// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExadataInfraDataSource struct{}

type ExadataInfraDataModel struct {
	Location          string                 `tfschema:"location"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Type              string                 `tfschema:"type"`
	Tags              map[string]interface{} `tfschema:"tags"`
	Zones             zones.Schema           `tfschema:"zones"`

	// SystemData
	SystemData []SystemDataModel `tfschema:"system_data"`

	// CloudExadataInfrastructureProperties
	ActivatedStorageCount       int64                        `tfschema:"activated_storage_count"`
	AdditionalStorageCount      int64                        `tfschema:"additional_storage_count"`
	AvailableStorageSizeInGbs   int64                        `tfschema:"available_storage_size_in_gbs"`
	ComputeCount                int64                        `tfschema:"compute_count"`
	CpuCount                    int64                        `tfschema:"cpu_count"`
	CustomerContacts            []string                     `tfschema:"customer_contacts"`
	DataStorageSizeInTbs        float64                      `tfschema:"data_storage_size_in_tbs"`
	DbNodeStorageSizeInGbs      int64                        `tfschema:"db_node_storage_size_in_gbs"`
	DbServerVersion             string                       `tfschema:"db_server_version"`
	DisplayName                 string                       `tfschema:"display_name"`
	EstimatedPatchingTime       []EstimatedPatchingTimeModel `tfschema:"estimated_patching_time"`
	LastMaintenanceRunId        string                       `tfschema:"last_maintenance_run_id"`
	LifecycleDetails            string                       `tfschema:"lifecycle_details"`
	LifecycleState              string                       `tfschema:"lifecycle_state"`
	MaintenanceWindow           []MaintenanceWindowModel     `tfschema:"maintenance_window"`
	MaxCPUCount                 int64                        `tfschema:"max_cpu_count"`
	MaxDataStorageInTbs         float64                      `tfschema:"max_data_storage_in_tbs"`
	MaxDbNodeStorageSizeInGbs   int64                        `tfschema:"max_db_node_storage_size_in_gbs"`
	MaxMemoryInGbs              int64                        `tfschema:"max_memory_in_gbs"`
	MemorySizeInGbs             int64                        `tfschema:"memory_size_in_gbs"`
	MonthlyDbServerVersion      string                       `tfschema:"monthly_db_server_version"`
	MonthlyStorageServerVersion string                       `tfschema:"monthly_storage_server_version"`
	NextMaintenanceRunId        string                       `tfschema:"next_maintenance_run_id"`
	OciUrl                      string                       `tfschema:"oci_url"`
	Ocid                        string                       `tfschema:"ocid"`
	ProvisioningState           string                       `tfschema:"provisioning_state"`
	Shape                       string                       `tfschema:"shape"`
	StorageCount                int64                        `tfschema:"storage_count"`
	StorageServerVersion        string                       `tfschema:"storage_server_version"`
	TimeCreated                 string                       `tfschema:"time_created"`
	TotalStorageSizeInGbs       int64                        `tfschema:"total_storage_size_in_gbs"`
}

type SystemDataModel struct {
	CreatedBy          string `tfschema:"created_by"`
	CreatedByType      string `tfschema:"created_by_type"`
	CreatedAt          string `tfschema:"created_at"`
	LastModifiedBy     string `tfschema:"last_modified_by"`
	LastModifiedbyType string `tfschema:"last_modified_by_type"`
	LastModifiedAt     string `tfschema:"last_modified_at"`
}

type EstimatedPatchingTimeModel struct {
	EstimatedDbServerPatchingTime        *int64 `tfschema:"estimated_db_server_patching_time"`
	EstimatedNetworkSwitchesPatchingTime *int64 `tfschema:"estimated_network_switches_patching_time"`
	EstimatedStorageServerPatchingTime   *int64 `tfschema:"estimated_storage_server_patching_time"`
	TotalEstimatedPatchingTime           *int64 `tfschema:"total_estimated_patching_time"`
}

type MaintenanceWindowModel struct {
	DaysOfWeek      []string `tfschema:"days_of_week"`
	HoursOfDay      []int64  `tfschema:"hours_of_day"`
	LeadTimeInWeeks int64    `tfschema:"lead_time_in_weeks"`
	Months          []string `tfschema:"months"`
	PatchingMode    string   `tfschema:"patching_mode"`
	Preference      string   `tfschema:"preference"`
	WeeksOfMonth    []int64  `tfschema:"weeks_of_month"`
}

func (d ExadataInfraDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d ExadataInfraDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// SystemData
		"system_data": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"created_by": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"created_by_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"created_at": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_modified_by": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_modified_by_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_modified_at": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		// CloudExadataInfrastructureProperties
		"activated_storage_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"additional_storage_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"available_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"compute_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"cpu_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"db_node_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"db_server_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"estimated_patching_time": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"estimated_db_server_patching_time": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"estimated_network_switches_patching_time": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"estimated_storage_server_patching_time": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"total_estimated_patching_time": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"last_maintenance_run_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_details": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"custom_action_timeout_in_mins": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"days_of_week": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"hours_of_day": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},

					"is_custom_action_timeout_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"is_monthly_patching_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"lead_time_in_weeks": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"months": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"patching_mode": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"preference": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"weeks_of_month": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},
				},
			},
		},

		"max_cpu_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_data_storage_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"max_db_node_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_memory_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"memory_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"monthly_db_server_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"monthly_storage_server_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"next_maintenance_run_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oci_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"shape": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"storage_server_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_created": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"total_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),

		"zones": commonschema.ZonesMultipleComputed(),
	}
}

func (d ExadataInfraDataSource) ModelObject() interface{} {
	return nil
}

func (d ExadataInfraDataSource) ResourceType() string {
	return "azurerm_oracledatabase_exadata_infrastructure"
}

func (d ExadataInfraDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID
}

func (d ExadataInfraDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudExadataInfrastructures
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID(subscriptionId,
				metadata.ResourceData.Get("resource_group_name").(string),
				metadata.ResourceData.Get("name").(string))

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				err := metadata.ResourceData.Set("location", location.NormalizeNilable(&model.Location))
				if err != nil {
					return err
				}

				var output ExadataInfraDataModel

				output.Name = id.CloudExadataInfrastructureName
				output.ResourceGroupName = id.ResourceGroupName
				output.Type = pointer.From(model.Type)
				output.Tags = utils.FlattenPtrMapStringString(model.Tags)
				output.Location = model.Location
				output.Zones = model.Zones

				prop := model.Properties
				if prop != nil {
					output = ExadataInfraDataModel{
						ActivatedStorageCount:       pointer.From(prop.ActivatedStorageCount),
						AdditionalStorageCount:      pointer.From(prop.AdditionalStorageCount),
						AvailableStorageSizeInGbs:   pointer.From(prop.AvailableStorageSizeInGbs),
						CpuCount:                    pointer.From(prop.CpuCount),
						ComputeCount:                pointer.From(prop.ComputeCount),
						CustomerContacts:            FlattenCustomerContacts(prop.CustomerContacts),
						DataStorageSizeInTbs:        pointer.From(prop.DataStorageSizeInTbs),
						DbNodeStorageSizeInGbs:      pointer.From(prop.DbNodeStorageSizeInGbs),
						DbServerVersion:             pointer.From(prop.DbServerVersion),
						DisplayName:                 prop.DisplayName,
						EstimatedPatchingTime:       FlattenEstimatedPatchingTimes(prop.EstimatedPatchingTime),
						LastMaintenanceRunId:        pointer.From(prop.LastMaintenanceRunId),
						LifecycleDetails:            pointer.From(prop.LifecycleDetails),
						LifecycleState:              string(*prop.LifecycleState),
						MaintenanceWindow:           FlattenMaintenanceWindow(prop.MaintenanceWindow),
						MaxCPUCount:                 pointer.From(prop.MaxCPUCount),
						MaxDataStorageInTbs:         pointer.From(prop.MaxDataStorageInTbs),
						MaxDbNodeStorageSizeInGbs:   pointer.From(prop.MaxDbNodeStorageSizeInGbs),
						MaxMemoryInGbs:              pointer.From(prop.MaxMemoryInGbs),
						MemorySizeInGbs:             pointer.From(prop.MemorySizeInGbs),
						MonthlyDbServerVersion:      pointer.From(prop.MonthlyDbServerVersion),
						MonthlyStorageServerVersion: pointer.From(prop.MonthlyStorageServerVersion),
						NextMaintenanceRunId:        pointer.From(prop.NextMaintenanceRunId),
						OciUrl:                      pointer.From(prop.OciUrl),
						Ocid:                        pointer.From(prop.Ocid),
						ProvisioningState:           string(*prop.ProvisioningState),
						Shape:                       prop.Shape,
						StorageCount:                pointer.From(prop.StorageCount),
						StorageServerVersion:        pointer.From(prop.StorageServerVersion),
						TimeCreated:                 pointer.From(prop.TimeCreated),
						TotalStorageSizeInGbs:       pointer.From(prop.TotalStorageSizeInGbs),
					}
				}

				systemData := model.SystemData
				if systemData != nil {
					output.SystemData = []SystemDataModel{
						{
							CreatedBy:          systemData.CreatedBy,
							CreatedByType:      systemData.CreatedByType,
							CreatedAt:          systemData.CreatedAt,
							LastModifiedBy:     systemData.LastModifiedBy,
							LastModifiedbyType: systemData.LastModifiedbyType,
							LastModifiedAt:     systemData.LastModifiedAt,
						},
					}
				}

				metadata.SetID(id)
				return metadata.Encode(&output)
			}
			return nil
		},
	}
}

func FlattenCustomerContacts(customerContactsList *[]cloudexadatainfrastructures.CustomerContact) []string {
	var customerContacts []string
	if customerContactsList != nil {
		for _, customerContact := range *customerContactsList {
			customerContacts = append(customerContacts, customerContact.Email)
		}
	}
	return customerContacts
}

func FlattenEstimatedPatchingTimes(estimatedPatchingTime *cloudexadatainfrastructures.EstimatedPatchingTime) []EstimatedPatchingTimeModel {
	if estimatedPatchingTime != nil {
		return []EstimatedPatchingTimeModel{
			{
				EstimatedDbServerPatchingTime:        estimatedPatchingTime.EstimatedDbServerPatchingTime,
				EstimatedNetworkSwitchesPatchingTime: estimatedPatchingTime.EstimatedNetworkSwitchesPatchingTime,
				EstimatedStorageServerPatchingTime:   estimatedPatchingTime.EstimatedStorageServerPatchingTime,
				TotalEstimatedPatchingTime:           estimatedPatchingTime.TotalEstimatedPatchingTime,
			},
		}
	}
	return nil
}

func FlattenMaintenanceWindow(maintenanceWindow *cloudexadatainfrastructures.MaintenanceWindow) []MaintenanceWindowModel {
	if maintenanceWindow != nil {
		return []MaintenanceWindowModel{
			{
				DaysOfWeek:      FlattenDayOfWeek(maintenanceWindow.DaysOfWeek),
				HoursOfDay:      pointer.From(maintenanceWindow.HoursOfDay),
				LeadTimeInWeeks: pointer.From(maintenanceWindow.LeadTimeInWeeks),
				Months:          FlattenMonths(maintenanceWindow.Months),
				PatchingMode:    string(pointer.From(maintenanceWindow.PatchingMode)),
				Preference:      string(pointer.From(maintenanceWindow.Preference)),
				WeeksOfMonth:    pointer.From(maintenanceWindow.WeeksOfMonth),
			},
		}
	}
	return nil
}

func FlattenDayOfWeek(dayOfWeeks *[]cloudexadatainfrastructures.DayOfWeek) []string {
	var dayOfWeeksArray []string
	if dayOfWeeks != nil {
		for _, dayOfWeek := range *dayOfWeeks {
			dayOfWeeksArray = append(dayOfWeeksArray, string(dayOfWeek.Name))
		}
	}
	return dayOfWeeksArray
}

func FlattenMonths(months *[]cloudexadatainfrastructures.Month) []string {
	var monthsArray []string
	if months != nil {
		for _, month := range *months {
			monthsArray = append(monthsArray, string(month.Name))
		}
	}
	return monthsArray
}
