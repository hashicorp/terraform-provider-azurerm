// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExadataInfraDataSource struct{}

type ExadataInfraDataModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

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
	Shape                       string                       `tfschema:"shape"`
	StorageCount                int64                        `tfschema:"storage_count"`
	StorageServerVersion        string                       `tfschema:"storage_server_version"`
	TimeCreated                 string                       `tfschema:"time_created"`
	TotalStorageSizeInGbs       int64                        `tfschema:"total_storage_size_in_gbs"`
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

func (d ExadataInfraDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ExadataName,
		},
	}
}

func (d ExadataInfraDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

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

					"custom_action_timeout_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"monthly_patching_enabled": {
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
	return &ExadataInfraDataModel{}
}

func (d ExadataInfraDataSource) ResourceType() string {
	return "azurerm_oracle_exadata_infrastructure"
}

func (d ExadataInfraDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID
}

func (d ExadataInfraDataSource) Read() sdk.ResourceFunc {
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
				state.Zones = model.Zones
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
			PatchingMode:    string(pointer.From(maintenanceWindow.PatchingMode)),
			Preference:      string(pointer.From(maintenanceWindow.Preference)),
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
