// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = AutonomousDatabaseCloneFromBackupDataSource{}

type AutonomousDatabaseCloneFromBackupDataSource struct{}

type AutonomousDatabaseCloneFromBackupDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`

	SourceAutonomousDatabaseId string `tfschema:"source_autonomous_database_id"`

	// Base properties (computed)
	AllowedIpAddresses                            []string                        `tfschema:"allowed_ip_addresses"`
	BackupRetentionPeriodInDays                   int64                           `tfschema:"backup_retention_period_in_days"`
	CharacterSet                                  string                          `tfschema:"character_set"`
	ComputeCount                                  float64                         `tfschema:"compute_count"`
	ComputeModel                                  string                          `tfschema:"compute_model"`
	ConnectionStrings                             []string                        `tfschema:"connection_strings"`
	CustomerContacts                              []string                        `tfschema:"customer_contacts"`
	DataStorageSizeInGb                           int64                           `tfschema:"data_storage_size_in_gb"`
	DataStorageSizeInTb                           int64                           `tfschema:"data_storage_size_in_tb"`
	DatabaseVersion                               string                          `tfschema:"database_version"`
	DatabaseWorkload                              string                          `tfschema:"database_workload"`
	DisplayName                                   string                          `tfschema:"display_name"`
	LicenseModel                                  string                          `tfschema:"license_model"`
	AutoScalingEnabled                            bool                            `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled                  bool                            `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired                        bool                            `tfschema:"mtls_connection_required"`
	NationalCharacterSet                          string                          `tfschema:"national_character_set"`
	SubnetId                                      string                          `tfschema:"subnet_id"`
	VnetId                                        string                          `tfschema:"virtual_network_id"`
	LifecycleState                                string                          `tfschema:"lifecycle_state"`
	PrivateEndpointUrl                            string                          `tfschema:"private_endpoint_url"`
	PrivateEndpointIp                             string                          `tfschema:"private_endpoint_ip"`
	ServiceConsoleUrl                             string                          `tfschema:"service_console_url"`
	SqlWebDeveloperUrl                            string                          `tfschema:"sql_web_developer_url"`
	TimeCreatedUtc                                string                          `tfschema:"time_created_in_utc"`
	OciUrl                                        string                          `tfschema:"oci_url"`
	ActualUsedDataStorageSizeInTb                 float64                         `tfschema:"actual_used_data_storage_size_in_tb"`
	AllocatedStorageSizeInTb                      float64                         `tfschema:"allocated_storage_size_in_tb"`
	AvailableUpgradeVersions                      []string                        `tfschema:"available_upgrade_versions"`
	CpuCoreCount                                  int64                           `tfschema:"cpu_core_count"`
	FailedDataRecoveryInSeconds                   int64                           `tfschema:"failed_data_recovery_in_seconds"`
	LifecycleDetails                              string                          `tfschema:"lifecycle_details"`
	LocalAdgAutoFailoverMaxDataLossLimitInSeconds int64                           `tfschema:"local_adg_auto_failover_max_data_loss_limit_in_seconds"`
	LocalDataGuardEnabled                         bool                            `tfschema:"local_data_guard_enabled"`
	LongTermBackupSchedule                        []LongTermBackUpScheduleDetails `tfschema:"long_term_backup_schedule"`
	MemoryAreaInGb                                int64                           `tfschema:"in_memory_area_in_gb"`
	MemoryPerOracleComputeUnitInGb                int64                           `tfschema:"memory_per_oracle_compute_unit_in_gb"`
	NextLongTermBackupTimestamp                   string                          `tfschema:"next_long_term_backup_timestamp"`
	Ocid                                          string                          `tfschema:"ocid"`
	PeerDatabaseId                                string                          `tfschema:"peer_database_id"`
	PeerDatabaseIds                               []string                        `tfschema:"peer_database_ids"`
	Preview                                       bool                            `tfschema:"preview"`
	PreviewVersionWithServiceTermsAccepted        bool                            `tfschema:"preview_version_with_service_terms_accepted"`
	PrivateEndpointLabel                          string                          `tfschema:"private_endpoint_label"`
	ProvisionableCPUs                             []int64                         `tfschema:"provisionable_cpus"`
	RemoteDataGuardEnabled                        bool                            `tfschema:"remote_data_guard_enabled"`
	SupportedRegionsToCloneTo                     []string                        `tfschema:"supported_regions_to_clone_to"`
	TimeDataGuardRoleChangedInUtc                 string                          `tfschema:"time_data_guard_role_changed_in_utc"`
	TimeDeletionOfFreeAutonomousDatabaseInUtc     string                          `tfschema:"time_deletion_of_free_autonomous_database_in_utc"`
	TimeLocalDataGuardEnabledInUtc                string                          `tfschema:"time_local_data_guard_enabled_in_utc"`
	TimeMaintenanceBeginInUtc                     string                          `tfschema:"time_maintenance_begin_in_utc"`
	TimeMaintenanceEndInUtc                       string                          `tfschema:"time_maintenance_end_in_utc"`
	TimeOfLastFailoverInUtc                       string                          `tfschema:"time_of_last_failover_in_utc"`
	TimeOfLastRefreshInUtc                        string                          `tfschema:"time_of_last_refresh_in_utc"`
	TimeOfLastRefreshPointInUtc                   string                          `tfschema:"time_of_last_refresh_point_in_utc"`
	TimeOfLastSwitchoverInUtc                     string                          `tfschema:"time_of_last_switchover_in_utc"`
	TimeReclamationOfFreeAutonomousDatabaseInUtc  string                          `tfschema:"time_reclamation_of_free_autonomous_database_in_utc"`
	UsedDataStorageSizeInGb                       int64                           `tfschema:"used_data_storage_size_in_gb"`
	UsedDataStorageSizeInTb                       int64                           `tfschema:"used_data_storage_size_in_tb"`
}

func (AutonomousDatabaseCloneFromBackupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (AutonomousDatabaseCloneFromBackupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		"source_autonomous_database_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"actual_used_data_storage_size_in_tb": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"allocated_storage_size_in_tb": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"allowed_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
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

		"local_adg_auto_failover_max_data_loss_limit_in_seconds": {
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
					"time_of_backup_in_utc": {
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

		"next_long_term_backup_timestamp": {
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

		"peer_database_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"peer_database_ids": {
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

		"private_endpoint_url": {
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

		"time_created_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_data_guard_role_changed_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_deletion_of_free_autonomous_database_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_local_data_guard_enabled_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_maintenance_begin_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_maintenance_end_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_failover_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_refresh_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_refresh_point_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_of_last_switchover_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_reclamation_of_free_autonomous_database_in_utc": {
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

func (AutonomousDatabaseCloneFromBackupDataSource) ModelObject() interface{} {
	return &AutonomousDatabaseCloneFromBackupDataSourceModel{}
}

func (AutonomousDatabaseCloneFromBackupDataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone_from_backup"
}

func (AutonomousDatabaseCloneFromBackupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state AutonomousDatabaseCloneFromBackupDataSourceModel
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
				state.Name = id.AutonomousDatabaseName
				state.ResourceGroupName = id.ResourceGroupName

				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `CloneFromBackupTimestamp`", id)
				}
				state.SourceAutonomousDatabaseId = props.SourceId

				// Base properties
				state.ActualUsedDataStorageSizeInTb = pointer.From(props.ActualUsedDataStorageSizeInTbs)
				state.AllocatedStorageSizeInTb = pointer.From(props.AllocatedStorageSizeInTbs)
				state.AllowedIpAddresses = pointer.From(props.WhitelistedIPs)
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.AvailableUpgradeVersions = pointer.From(props.AvailableUpgradeVersions)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = pointer.FromEnum(props.ComputeModel)
				state.CpuCoreCount = pointer.From(props.CpuCoreCount)
				state.CustomerContacts = flattenAdbsCustomerContacts(props.CustomerContacts)
				state.DataStorageSizeInGb = pointer.From(props.DataStorageSizeInGbs)
				state.DataStorageSizeInTb = pointer.From(props.DataStorageSizeInTbs)
				state.DatabaseVersion = pointer.From(props.DbVersion)
				state.DatabaseWorkload = pointer.FromEnum(props.DbWorkload)
				state.DisplayName = pointer.From(props.DisplayName)
				state.FailedDataRecoveryInSeconds = pointer.From(props.FailedDataRecoveryInSeconds)
				state.LicenseModel = pointer.FromEnum(props.LicenseModel)
				state.LifecycleDetails = pointer.From(props.LifecycleDetails)
				state.LifecycleState = pointer.FromEnum(props.LifecycleState)
				state.LocalAdgAutoFailoverMaxDataLossLimitInSeconds = pointer.From(props.LocalAdgAutoFailoverMaxDataLossLimit)
				state.LocalDataGuardEnabled = pointer.From(props.IsLocalDataGuardEnabled)
				state.LongTermBackupSchedule = FlattenLongTermBackUpScheduleDetails(props.LongTermBackupSchedule)
				state.MemoryAreaInGb = pointer.From(props.InMemoryAreaInGbs)
				state.MemoryPerOracleComputeUnitInGb = pointer.From(props.MemoryPerOracleComputeUnitInGbs)
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.NextLongTermBackupTimestamp = pointer.From(props.NextLongTermBackupTimeStamp)
				state.OciUrl = pointer.From(props.OciURL)
				state.Ocid = pointer.From(props.Ocid)
				state.PeerDatabaseId = pointer.From(props.PeerDbId)
				state.PeerDatabaseIds = pointer.From(props.PeerDbIds)
				state.Preview = pointer.From(props.IsPreview)
				state.PreviewVersionWithServiceTermsAccepted = pointer.From(props.IsPreviewVersionWithServiceTermsAccepted)
				state.PrivateEndpointUrl = pointer.From(props.PrivateEndpoint)
				state.PrivateEndpointIp = pointer.From(props.PrivateEndpointIP)
				state.PrivateEndpointLabel = pointer.From(props.PrivateEndpointLabel)
				state.ProvisionableCPUs = pointer.From(props.ProvisionableCPUs)
				state.RemoteDataGuardEnabled = pointer.From(props.IsRemoteDataGuardEnabled)
				state.ServiceConsoleUrl = pointer.From(props.ServiceConsoleURL)
				state.SqlWebDeveloperUrl = pointer.From(props.SqlWebDeveloperURL)
				state.SubnetId = pointer.From(props.SubnetId)
				state.SupportedRegionsToCloneTo = pointer.From(props.SupportedRegionsToCloneTo)
				state.TimeCreatedUtc = pointer.From(props.TimeCreated)
				state.TimeDataGuardRoleChangedInUtc = pointer.From(props.TimeDataGuardRoleChanged)
				state.TimeDeletionOfFreeAutonomousDatabaseInUtc = pointer.From(props.TimeDeletionOfFreeAutonomousDatabase)
				state.TimeLocalDataGuardEnabledInUtc = pointer.From(props.TimeLocalDataGuardEnabled)
				state.TimeMaintenanceBeginInUtc = pointer.From(props.TimeMaintenanceBegin)
				state.TimeMaintenanceEndInUtc = pointer.From(props.TimeMaintenanceEnd)
				state.TimeOfLastFailoverInUtc = pointer.From(props.TimeOfLastFailover)
				state.TimeOfLastRefreshInUtc = pointer.From(props.TimeOfLastRefresh)
				state.TimeOfLastRefreshPointInUtc = pointer.From(props.TimeOfLastRefreshPoint)
				state.TimeOfLastSwitchoverInUtc = pointer.From(props.TimeOfLastSwitchover)
				state.TimeReclamationOfFreeAutonomousDatabaseInUtc = pointer.From(props.TimeReclamationOfFreeAutonomousDatabase)
				state.UsedDataStorageSizeInGb = pointer.From(props.UsedDataStorageSizeInGbs)
				state.UsedDataStorageSizeInTb = pointer.From(props.UsedDataStorageSizeInTbs)
				state.VnetId = pointer.From(props.VnetId)
				state.ConnectionStrings = flattenConnectionStrings(props.ConnectionStrings)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
