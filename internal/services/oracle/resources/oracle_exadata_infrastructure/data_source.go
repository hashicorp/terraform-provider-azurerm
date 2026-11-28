// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exadata_infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type ExadataInfraDataModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	// CloudExadataInfrastructureProperties
	ActivatedStorageCount          int64                                 `tfschema:"activated_storage_count"`
	AdditionalStorageCount         int64                                 `tfschema:"additional_storage_count"`
	AvailableStorageSizeInGbs      int64                                 `tfschema:"available_storage_size_in_gbs"`
	ComputeCount                   int64                                 `tfschema:"compute_count"`
	ComputeModel                   string                                `tfschema:"compute_model"`
	CpuCount                       int64                                 `tfschema:"cpu_count"`
	CustomerContacts               []string                              `tfschema:"customer_contacts"`
	DatabaseServerType             string                                `tfschema:"database_server_type"`
	DataStorageSizeInTbs           float64                               `tfschema:"data_storage_size_in_tbs"`
	DbNodeStorageSizeInGbs         int64                                 `tfschema:"db_node_storage_size_in_gbs"`
	DbServerVersion                string                                `tfschema:"db_server_version"`
	DisplayName                    string                                `tfschema:"display_name"`
	EstimatedPatchingTime          []EstimatedPatchingTimeModel          `tfschema:"estimated_patching_time"`
	LastMaintenanceRunId           string                                `tfschema:"last_maintenance_run_id"`
	LifecycleDetails               string                                `tfschema:"lifecycle_details"`
	LifecycleState                 string                                `tfschema:"lifecycle_state"`
	MaintenanceWindow              []MaintenanceWindowModel              `tfschema:"maintenance_window"`
	MaxCPUCount                    int64                                 `tfschema:"max_cpu_count"`
	MaxDataStorageInTbs            float64                               `tfschema:"max_data_storage_in_tbs"`
	MaxDbNodeStorageSizeInGbs      int64                                 `tfschema:"max_db_node_storage_size_in_gbs"`
	MaxMemoryInGbs                 int64                                 `tfschema:"max_memory_in_gbs"`
	MemorySizeInGbs                int64                                 `tfschema:"memory_size_in_gbs"`
	MonthlyDbServerVersion         string                                `tfschema:"monthly_db_server_version"`
	MonthlyStorageServerVersion    string                                `tfschema:"monthly_storage_server_version"`
	NextMaintenanceRunId           string                                `tfschema:"next_maintenance_run_id"`
	OciUrl                         string                                `tfschema:"oci_url"`
	Ocid                           string                                `tfschema:"ocid"`
	Shape                          string                                `tfschema:"shape"`
	StorageCount                   int64                                 `tfschema:"storage_count"`
	StorageServerType              string                                `tfschema:"storage_server_type"`
	StorageServerVersion           string                                `tfschema:"storage_server_version"`
	TimeCreated                    string                                `tfschema:"time_created"`
	TotalStorageSizeInGbs          int64                                 `tfschema:"total_storage_size_in_gbs"`
	DefinedFileSystemConfiguration []DefinedFileSystemConfigurationModel `tfschema:"defined_file_system_configuration"`
}

type DefinedFileSystemConfigurationModel struct {
	BackupPartitionEnabled bool   `tfschema:"backup_partition_enabled"`
	ResizableEnabled       bool   `tfschema:"resizable_enabled"`
	MinimumSizeGb          int64  `tfschema:"minimum_size_in_gb"`
	MountPoint             string `tfschema:"mount_point"`
}

type EstimatedPatchingTimeModel struct {
	EstimatedDbServerPatchingTime        int64 `tfschema:"estimated_db_server_patching_time"`
	EstimatedNetworkSwitchesPatchingTime int64 `tfschema:"estimated_network_switches_patching_time"`
	EstimatedStorageServerPatchingTime   int64 `tfschema:"estimated_storage_server_patching_time"`
	TotalEstimatedPatchingTime           int64 `tfschema:"total_estimated_patching_time"`
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

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Description:  "The name of this Cloud Exadata Infrastructure.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// CloudExadataInfrastructureProperties
		"activated_storage_count": {
			Description: "The requested number of additional storage servers activated for the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"additional_storage_count": {
			Description: "The requested number of additional storage servers for the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"available_storage_size_in_gbs": {
			Description: "The available storage can be allocated to the Cloud Exadata Infrastructure resource, in gigabytes (GB).",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"compute_count": {
			Description: "The number of compute servers for the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"compute_model": {
			Description: "The compute model of the Exadata Infrastructure.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"cpu_count": {
			Description: "The total number of CPU cores allocated.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"database_server_type": {
			Description: "The database server model type of the cloud Exadata infrastructure resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"data_storage_size_in_tbs": {
			Description: "The data storage size in terabytes of the DATA disk group.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"db_node_storage_size_in_gbs": {
			Description: "The local node storage allocated in GBs.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"db_server_version": {
			Description: "The software version of the database servers (dom0) in the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The user-friendly name for the Cloud Exadata Infrastructure resource. The name does not need to be unique.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"estimated_patching_time": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"estimated_db_server_patching_time": {
						Description: "The estimated time required in minutes for database server patching.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"estimated_network_switches_patching_time": {
						Description: "The estimated time required in minutes for network switch patching.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"estimated_storage_server_patching_time": {
						Description: "The estimated time required in minutes for storage server patching.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"total_estimated_patching_time": {
						Description: "The estimated total time required in minutes for all patching operations.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
				},
			},
		},

		"last_maintenance_run_id": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the last maintenance run.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"lifecycle_details": {
			Description: "Additional information about the current lifecycle state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"lifecycle_state": {
			Description: "Cloud Exadata Infrastructure lifecycle state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
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
						Description: "Days during the week when maintenance should be performed.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"hours_of_day": {
						Description: "The window of hours during the day when maintenance should be performed.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},

					"custom_action_timeout_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"monthly_patching_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"lead_time_in_weeks": {
						Description: "Lead time window allows user to set a lead time to prepare for a down time.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"months": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"patching_mode": {
						Description: "Cloud Exadata Infrastructure node patching method.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"preference": {
						Description: "The maintenance window scheduling preference.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"weeks_of_month": {
						Description: "Weeks during the month when maintenance should be performed.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},
				},
			},
		},

		"max_cpu_count": {
			Description: "The total number of CPU cores available.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"max_data_storage_in_tbs": {
			Description: "The total available DATA disk group size.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"max_db_node_storage_size_in_gbs": {
			Description: "The total local node storage available in GBs.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"max_memory_in_gbs": {
			Description: "The total memory available in GBs.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"memory_size_in_gbs": {
			Description: "The memory allocated in GBs.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"monthly_db_server_version": {
			Description: "The monthly software version of the database servers (dom0) in the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"monthly_storage_server_version": {
			Description: "The monthly software version of the storage servers (cells) in the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"next_maintenance_run_id": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the next maintenance run.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"oci_url": {
			Description: "The URL of the resource in the OCI console.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ocid": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"shape": {
			Description: "The model name of the Cloud Exadata Infrastructure resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"storage_count": {
			Description: "The number of storage servers for the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"storage_server_type": {
			Description: "The storage server model type of the cloud Exadata infrastructure resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"storage_server_version": {
			Description: "The software version of the storage servers (cells) in the Cloud Exadata Infrastructure.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_created": {
			Description: "The date and time the Cloud Exadata Infrastructure resource was created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"total_storage_size_in_gbs": {
			Description: "The total storage allocated to the Cloud Exadata Infrastructure resource, in gigabytes (GB).",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),

		"zones": commonschema.ZonesMultipleComputed(),

		"defined_file_system_configuration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backup_partition_enabled": {
						Description: "Whether the backup partition is enabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
					"resizable_enabled": {
						Description: "Whether the resizable is enabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
					"minimum_size_in_gb": {
						Description: "The minimum size of the file system in GB.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
					"mount_point": {
						Description: "Mount path for the file system.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d DataSource) ModelObject() any {
	return &ExadataInfraDataModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_oracle_exadata_infrastructure"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Cloud Exadata Infrastructure."
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID
}

func (d DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudExadataInfrastructures
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ExadataInfraDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				state.Location = location.Normalize(model.Location)
				state.Zones = zones.Flatten(&model.Zones)
				if props := model.Properties; props != nil {
					state.ActivatedStorageCount = pointer.From(props.ActivatedStorageCount)
					state.ActivatedStorageCount = pointer.From(props.ActivatedStorageCount)
					state.AdditionalStorageCount = pointer.From(props.AdditionalStorageCount)
					state.AvailableStorageSizeInGbs = pointer.From(props.AvailableStorageSizeInGbs)
					state.CpuCount = pointer.From(props.CpuCount)
					state.ComputeCount = pointer.From(props.ComputeCount)
					state.CustomerContacts = FlattenCustomerContacts(props.CustomerContacts)
					state.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
					state.DbNodeStorageSizeInGbs = pointer.From(props.DbNodeStorageSizeInGbs)
					state.DbServerVersion = pointer.From(props.DbServerVersion)
					state.DisplayName = props.DisplayName
					state.EstimatedPatchingTime = FlattenEstimatedPatchingTimes(props.EstimatedPatchingTime)
					state.LastMaintenanceRunId = pointer.From(props.LastMaintenanceRunId)
					state.LifecycleDetails = pointer.From(props.LifecycleDetails)
					state.LifecycleState = string(*props.LifecycleState)
					state.MaintenanceWindow = FlattenMaintenanceWindow(props.MaintenanceWindow)
					state.MaxCPUCount = pointer.From(props.MaxCPUCount)
					state.MaxDataStorageInTbs = pointer.From(props.MaxDataStorageInTbs)
					state.MaxDbNodeStorageSizeInGbs = pointer.From(props.MaxDbNodeStorageSizeInGbs)
					state.MaxMemoryInGbs = pointer.From(props.MaxMemoryInGbs)
					state.MemorySizeInGbs = pointer.From(props.MemorySizeInGbs)
					state.MonthlyDbServerVersion = pointer.From(props.MonthlyDbServerVersion)
					state.MonthlyStorageServerVersion = pointer.From(props.MonthlyStorageServerVersion)
					state.NextMaintenanceRunId = pointer.From(props.NextMaintenanceRunId)
					state.OciUrl = pointer.From(props.OciURL)
					state.Ocid = pointer.From(props.Ocid)
					state.Shape = props.Shape
					state.StorageCount = pointer.From(props.StorageCount)
					state.StorageServerVersion = pointer.From(props.StorageServerVersion)
					state.TimeCreated = pointer.From(props.TimeCreated)
					state.TotalStorageSizeInGbs = pointer.From(props.TotalStorageSizeInGbs)
					state.ComputeModel = pointer.FromEnum(props.ComputeModel)
					state.DatabaseServerType = pointer.From(props.DatabaseServerType)
					state.StorageServerType = pointer.From(props.StorageServerType)
					state.DefinedFileSystemConfiguration = FlattenDefinedFileSystemConfiguration(props.DefinedFileSystemConfiguration)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
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
	estimatedPatchingTimes := make([]EstimatedPatchingTimeModel, 0)
	if estimatedPatchingTime != nil {
		return append(estimatedPatchingTimes, EstimatedPatchingTimeModel{
			EstimatedDbServerPatchingTime:        pointer.From(estimatedPatchingTime.EstimatedDbServerPatchingTime),
			EstimatedNetworkSwitchesPatchingTime: pointer.From(estimatedPatchingTime.EstimatedNetworkSwitchesPatchingTime),
			EstimatedStorageServerPatchingTime:   pointer.From(estimatedPatchingTime.EstimatedStorageServerPatchingTime),
			TotalEstimatedPatchingTime:           pointer.From(estimatedPatchingTime.TotalEstimatedPatchingTime),
		})
	}
	return estimatedPatchingTimes
}

func FlattenMaintenanceWindow(maintenanceWindow *cloudexadatainfrastructures.MaintenanceWindow) []MaintenanceWindowModel {
	output := make([]MaintenanceWindowModel, 0)
	if maintenanceWindow != nil {
		return append(output, MaintenanceWindowModel{
			DaysOfWeek:      FlattenDayOfWeek(maintenanceWindow.DaysOfWeek),
			HoursOfDay:      pointer.From(maintenanceWindow.HoursOfDay),
			LeadTimeInWeeks: pointer.From(maintenanceWindow.LeadTimeInWeeks),
			Months:          FlattenMonths(maintenanceWindow.Months),
			PatchingMode:    pointer.FromEnum(maintenanceWindow.PatchingMode),
			Preference:      pointer.FromEnum(maintenanceWindow.Preference),
			WeeksOfMonth:    pointer.From(maintenanceWindow.WeeksOfMonth),
		})
	}
	return output
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

func FlattenDefinedFileSystemConfiguration(input *[]cloudexadatainfrastructures.DefinedFileSystemConfiguration) []DefinedFileSystemConfigurationModel {
	output := make([]DefinedFileSystemConfigurationModel, 0)
	if input != nil {
		for _, item := range *input {
			output = append(output, DefinedFileSystemConfigurationModel{
				BackupPartitionEnabled: pointer.From(item.IsBackupPartition),
				ResizableEnabled:       pointer.From(item.IsResizable),
				MinimumSizeGb:          pointer.From(item.MinSizeGb),
				MountPoint:             pointer.From(item.MountPoint),
			})
		}
	}
	return output
}
