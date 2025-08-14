// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = AutonomousDatabaseCloneFromDatabaseDataSource{}

type AutonomousDatabaseCloneFromDatabaseDataSource struct{}

type AutonomousDatabaseCloneFromDatabaseDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`

	// Clone-specific properties
	SourceId                string `tfschema:"source_id"`
	CloneType               string `tfschema:"clone_type"`
	ReconnectCloneEnabled   bool   `tfschema:"reconnect_clone_enabled"`
	RefreshableClone        bool   `tfschema:"refreshable_clone"`
	RefreshableModel        string `tfschema:"refreshable_model"`
	RefreshableStatus       string `tfschema:"refreshable_status"`
	TimeUntilReconnectClone string `tfschema:"time_until_reconnect_clone"`

	// Base properties (computed)
	ActualUsedDataStorageSizeInTbs          float64                         `tfschema:"actual_used_data_storage_size_in_tbs"`
	AllocatedStorageSizeInTb                float64                         `tfschema:"allocated_storage_size_in_tb"`
	AllowedIps                              []string                        `tfschema:"allowed_ips"`
	AutoScalingEnabled                      bool                            `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled            bool                            `tfschema:"auto_scaling_for_storage_enabled"`
	AutonomousDatabaseId                    string                          `tfschema:"autonomous_database_id"`
	AvailableUpgradeVersions                []string                        `tfschema:"available_upgrade_versions"`
	BackupRetentionPeriodInDays             int64                           `tfschema:"backup_retention_period_in_days"`
	CharacterSet                            string                          `tfschema:"character_set"`
	ComputeCount                            float64                         `tfschema:"compute_count"`
	ComputeModel                            string                          `tfschema:"compute_model"`
	ConnectionStrings                       []string                        `tfschema:"connection_strings"`
	CpuCoreCount                            int64                           `tfschema:"cpu_core_count"`
	CustomerContacts                        []string                        `tfschema:"customer_contacts"`
	DataStorageSizeInGb                     int64                           `tfschema:"data_storage_size_in_gb"`
	DataStorageSizeInTb                     int64                           `tfschema:"data_storage_size_in_tb"`
	DatabaseVersion                         string                          `tfschema:"database_version"`
	DatabaseWorkload                        string                          `tfschema:"database_workload"`
	DisplayName                             string                          `tfschema:"display_name"`
	FailedDataRecoveryInSeconds             int64                           `tfschema:"failed_data_recovery_in_seconds"`
	LicenseModel                            string                          `tfschema:"license_model"`
	LifecycleDetails                        string                          `tfschema:"lifecycle_details"`
	LifecycleState                          string                          `tfschema:"lifecycle_state"`
	LocalAdgAutoFailoverMaxDataLossLimit    int64                           `tfschema:"local_adg_auto_failover_max_data_loss_limit"`
	LocalDataGuardEnabled                   bool                            `tfschema:"local_data_guard_enabled"`
	LongTermBackupSchedule                  []LongTermBackUpScheduleDetails `tfschema:"long_term_backup_schedule"`
	MemoryAreaInGb                          int64                           `tfschema:"in_memory_area_in_gbs"`
	MemoryPerOracleComputeUnitInGb          int64                           `tfschema:"memory_per_oracle_compute_unit_in_gb"`
	MtlsConnectionRequired                  bool                            `tfschema:"mtls_connection_required"`
	NationalCharacterSet                    string                          `tfschema:"national_character_set"`
	NextLongTermBackupTimeStamp             string                          `tfschema:"next_long_term_backup_time_stamp"`
	OciUrl                                  string                          `tfschema:"oci_url"`
	Ocid                                    string                          `tfschema:"ocid"`
	PeerDbId                                string                          `tfschema:"peer_db_id"`
	PeerDbIds                               []string                        `tfschema:"peer_db_ids"`
	Preview                                 bool                            `tfschema:"preview"`
	PreviewVersionWithServiceTermsAccepted  bool                            `tfschema:"preview_version_with_service_terms_accepted"`
	PrivateEndpoint                         string                          `tfschema:"private_endpoint"`
	PrivateEndpointIp                       string                          `tfschema:"private_endpoint_ip"`
	PrivateEndpointLabel                    string                          `tfschema:"private_endpoint_label"`
	ProvisionableCPUs                       []int64                         `tfschema:"provisionable_cpus"`
	RemoteDataGuardEnabled                  bool                            `tfschema:"remote_data_guard_enabled"`
	ServiceConsoleUrl                       string                          `tfschema:"service_console_url"`
	SqlWebDeveloperUrl                      string                          `tfschema:"sql_web_developer_url"`
	SubnetId                                string                          `tfschema:"subnet_id"`
	SupportedRegionsToCloneTo               []string                        `tfschema:"supported_regions_to_clone_to"`
	TimeCreated                             string                          `tfschema:"time_created"`
	TimeDataGuardRoleChanged                string                          `tfschema:"time_data_guard_role_changed"`
	TimeDeletionOfFreeAutonomousDatabase    string                          `tfschema:"time_deletion_of_free_autonomous_database"`
	TimeLocalDataGuardEnabled               string                          `tfschema:"time_local_data_guard_enabled_on"`
	TimeMaintenanceBegin                    string                          `tfschema:"time_maintenance_begin"`
	TimeMaintenanceEnd                      string                          `tfschema:"time_maintenance_end"`
	TimeOfLastFailover                      string                          `tfschema:"time_of_last_failover"`
	TimeOfLastRefresh                       string                          `tfschema:"time_of_last_refresh"`
	TimeOfLastRefreshPoint                  string                          `tfschema:"time_of_last_refresh_point"`
	TimeOfLastSwitchover                    string                          `tfschema:"time_of_last_switchover"`
	TimeReclamationOfFreeAutonomousDatabase string                          `tfschema:"time_reclamation_of_free_autonomous_database"`
	UsedDataStorageSizeInGb                 int64                           `tfschema:"used_data_storage_size_in_gb"`
	UsedDataStorageSizeInTb                 int64                           `tfschema:"used_data_storage_size_in_tb"`
	VnetId                                  string                          `tfschema:"virtual_network_id"`
}

func (AutonomousDatabaseCloneFromDatabaseDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (AutonomousDatabaseCloneFromDatabaseDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		// Clone-specific properties
		"source_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"clone_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reconnect_clone_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"refreshable_clone": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"refreshable_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"refreshable_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_until_reconnect_clone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// Base properties
		"autonomous_database_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"actual_used_data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"allocated_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"allowed_ips": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeInt,
			},
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"available_upgrade_versions": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backup_retention_period_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"compute_count": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"compute_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"connection_strings": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cpu_core_count": {
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

		"data_storage_size_in_gb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"data_storage_size_in_tb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"database_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"database_workload": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"failed_data_recovery_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"in_memory_area_in_gb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"license_model": {
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

		"local_adg_auto_failover_max_data_loss_limit": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"local_data_guard_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"long_term_backup_schedule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"repeat_cadence": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"time_of_backup": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"retention_period_in_days": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"memory_per_oracle_compute_unit_in_gb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"national_character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"next_long_term_backup_time_stamp": {
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

		"peer_db_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"peer_db_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"preview": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"preview_version_with_service_terms_accepted": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"private_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_endpoint_ip": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_endpoint_label": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"provisionable_cpus": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeInt,
			},
		},

		"remote_data_guard_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"service_console_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sql_web_developer_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"supported_regions_to_clone_to": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"time_created": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_data_guard_role_changed": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_deletion_of_free_autonomous_database": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_local_data_guard_enabled_on": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_maintenance_begin": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_maintenance_end": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_failover": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_refresh": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_refresh_point": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_switchover": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_reclamation_of_free_autonomous_database": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"used_data_storage_size_in_gb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"used_data_storage_size_in_tb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (AutonomousDatabaseCloneFromDatabaseDataSource) ModelObject() interface{} {
	return &AutonomousDatabaseCloneFromDatabaseDataSourceModel{}
}

func (AutonomousDatabaseCloneFromDatabaseDataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone_from_database"
}

func (AutonomousDatabaseCloneFromDatabaseDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			var state AutonomousDatabaseCloneFromDatabaseDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseCloneProperties)
				if !ok {
					return fmt.Errorf("%s is not a clone type autonomous database", id)
				}
				state.CloneType = string(props.CloneType)
				state.SourceId = props.SourceId
				state.ReconnectCloneEnabled = pointer.From(props.IsReconnectCloneEnabled)
				state.RefreshableClone = pointer.From(props.IsRefreshableClone)
				state.TimeUntilReconnectClone = pointer.From(props.TimeUntilReconnectCloneEnabled)
				state.RefreshableModel = pointer.FromEnum(props.RefreshableModel)
				state.RefreshableStatus = pointer.FromEnum(props.RefreshableStatus)
				state.AutonomousDatabaseId = pointer.From(props.AutonomousDatabaseId)
				state.AllowedIps = pointer.From(props.WhitelistedIPs)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = string(pointer.From(props.ComputeModel))
				state.CustomerContacts = flattenAdbsCustomerContacts(props.CustomerContacts)
				state.DataStorageSizeInGb = pointer.From(props.DataStorageSizeInGbs)
				state.DataStorageSizeInTb = pointer.From(props.DataStorageSizeInTbs)
				state.DatabaseVersion = pointer.From(props.DbVersion)
				state.DatabaseWorkload = string(pointer.From(props.DbWorkload))
				state.DisplayName = pointer.From(props.DisplayName)
				state.LicenseModel = string(pointer.From(props.LicenseModel))
				state.LifecycleState = pointer.FromEnum(props.LifecycleState)
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.SubnetId = pointer.From(props.SubnetId)
				state.VnetId = pointer.From(props.VnetId)
				state.PrivateEndpoint = pointer.From(props.PrivateEndpoint)
				state.PrivateEndpointIp = pointer.From(props.PrivateEndpointIP)
				state.ServiceConsoleUrl = pointer.From(props.ServiceConsoleURL)
				state.SqlWebDeveloperUrl = pointer.From(props.SqlWebDeveloperURL)
				state.TimeCreated = pointer.From(props.TimeCreated)
				state.OciUrl = pointer.From(props.OciURL)
				if props.ConnectionStrings != nil && props.ConnectionStrings.AllConnectionStrings != nil {
					connStrings := make([]string, 0)
					allConnStrings := *props.ConnectionStrings.AllConnectionStrings
					if allConnStrings.High != nil {
						connStrings = append(connStrings, *allConnStrings.High)
					}
					if allConnStrings.Medium != nil {
						connStrings = append(connStrings, *allConnStrings.Medium)
					}
					if allConnStrings.Low != nil {
						connStrings = append(connStrings, *allConnStrings.Low)
					}
					state.ConnectionStrings = connStrings
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
