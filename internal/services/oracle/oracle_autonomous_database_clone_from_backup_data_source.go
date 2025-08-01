// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = AutonomousDatabaseCloneFromBackupDataSource{}

type AutonomousDatabaseCloneFromBackupDataSource struct{}

type AutonomousDatabaseCloneFomBackupDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`

	SourceAutonomousDatabaseId string `tfschema:"source_autonomous_database_id"`

	// Base properties (computed)
	AutonomousDatabaseId                    string                          `tfschema:"autonomous_database_id"`
	AllowedIps                              []string                        `tfschema:"allowed_ips"`
	BackupRetentionPeriodInDays             int64                           `tfschema:"backup_retention_period_in_days"`
	CharacterSet                            string                          `tfschema:"character_set"`
	ComputeCount                            float64                         `tfschema:"compute_count"`
	ComputeModel                            string                          `tfschema:"compute_model"`
	CustomerContacts                        []string                        `tfschema:"customer_contacts"`
	DataStorageSizeInGbs                    int64                           `tfschema:"data_storage_size_in_gbs"`
	DataStorageSizeInTbs                    int64                           `tfschema:"data_storage_size_in_tbs"`
	DbVersion                               string                          `tfschema:"db_version"`
	DbWorkload                              string                          `tfschema:"db_workload"`
	DisplayName                             string                          `tfschema:"display_name"`
	LicenseModel                            string                          `tfschema:"license_model"`
	AutoScalingEnabled                      bool                            `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled            bool                            `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired                  bool                            `tfschema:"mtls_connection_required"`
	NationalCharacterSet                    string                          `tfschema:"national_character_set"`
	SubnetId                                string                          `tfschema:"subnet_id"`
	VnetId                                  string                          `tfschema:"virtual_network_id"`
	LifecycleState                          string                          `tfschema:"lifecycle_state"`
	PrivateEndpoint                         string                          `tfschema:"private_endpoint"`
	PrivateEndpointIp                       string                          `tfschema:"private_endpoint_ip"`
	ServiceConsoleUrl                       string                          `tfschema:"service_console_url"`
	SqlWebDeveloperUrl                      string                          `tfschema:"sql_web_developer_url"`
	TimeCreated                             string                          `tfschema:"time_created"`
	OciUrl                                  string                          `tfschema:"oci_url"`
	ActualUsedDataStorageSizeInTbs          float64                         `tfschema:"actual_used_data_storage_size_in_tbs"`
	AllocatedStorageSizeInTbs               float64                         `tfschema:"allocated_storage_size_in_tbs"`
	AvailableUpgradeVersions                []string                        `tfschema:"available_upgrade_versions"`
	CpuCoreCount                            int64                           `tfschema:"cpu_core_count"`
	FailedDataRecoveryInSeconds             int64                           `tfschema:"failed_data_recovery_in_seconds"`
	LifecycleDetails                        string                          `tfschema:"lifecycle_details"`
	LocalAdgAutoFailoverMaxDataLossLimit    int64                           `tfschema:"local_adg_auto_failover_max_data_loss_limit"`
	LocalDataGuardEnabled                   bool                            `tfschema:"local_data_guard_enabled"`
	LongTermBackupSchedule                  []LongTermBackUpScheduleDetails `tfschema:"long_term_backup_schedule"`
	MemoryAreaInGbs                         int64                           `tfschema:"in_memory_area_in_gbs"`
	MemoryPerOracleComputeUnitInGbs         int64                           `tfschema:"memory_per_oracle_compute_unit_in_gbs"`
	NextLongTermBackupTimeStamp             string                          `tfschema:"next_long_term_backup_time_stamp"`
	Ocid                                    string                          `tfschema:"ocid"`
	PeerDbId                                string                          `tfschema:"peer_db_id"`
	PeerDbIds                               []string                        `tfschema:"peer_db_ids"`
	Preview                                 bool                            `tfschema:"preview"`
	PreviewVersionWithServiceTermsAccepted  bool                            `tfschema:"preview_version_with_service_terms_accepted"`
	PrivateEndpointLabel                    string                          `tfschema:"private_endpoint_label"`
	ProvisionableCPUs                       []int64                         `tfschema:"provisionable_cpus"`
	RemoteDataGuardEnabled                  bool                            `tfschema:"remote_data_guard_enabled"`
	SupportedRegionsToCloneTo               []string                        `tfschema:"supported_regions_to_clone_to"`
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
	UsedDataStorageSizeInGbs                int64                           `tfschema:"used_data_storage_size_in_gbs"`
	UsedDataStorageSizeInTbs                int64                           `tfschema:"used_data_storage_size_in_tbs"`
}

func (AutonomousDatabaseCloneFromBackupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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

		"data_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"db_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"db_workload": {
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

		"in_memory_area_in_gbs": {
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

		"memory_per_oracle_compute_unit_in_gbs": {
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

		"used_data_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"used_data_storage_size_in_tbs": {
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
	return &AutonomousDatabaseCloneFomBackupDataSourceModel{}
}

func (AutonomousDatabaseCloneFromBackupDataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone_from_backup"
}

func (AutonomousDatabaseCloneFromBackupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			var state AutonomousDatabaseCloneFomBackupDataSourceModel
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
					return fmt.Errorf("%s is not a clone type autonomous database", id)
				}
				state.SourceAutonomousDatabaseId = props.SourceId

				// Base properties
				state.ActualUsedDataStorageSizeInTbs = pointer.From(props.ActualUsedDataStorageSizeInTbs)
				state.AllocatedStorageSizeInTbs = pointer.From(props.AllocatedStorageSizeInTbs)
				state.AllowedIps = pointer.From(props.WhitelistedIPs)
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.AutonomousDatabaseId = pointer.From(props.AutonomousDatabaseId)
				state.AvailableUpgradeVersions = pointer.From(props.AvailableUpgradeVersions)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = string(pointer.From(props.ComputeModel))
				state.CpuCoreCount = pointer.From(props.CpuCoreCount)
				state.CustomerContacts = flattenAdbsCustomerContacts(props.CustomerContacts)
				state.DataStorageSizeInGbs = pointer.From(props.DataStorageSizeInGbs)
				state.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
				state.DbVersion = pointer.From(props.DbVersion)
				state.DbWorkload = string(pointer.From(props.DbWorkload))
				state.DisplayName = pointer.From(props.DisplayName)
				state.FailedDataRecoveryInSeconds = pointer.From(props.FailedDataRecoveryInSeconds)
				state.LicenseModel = string(pointer.From(props.LicenseModel))
				state.LifecycleDetails = pointer.From(props.LifecycleDetails)
				state.LifecycleState = pointer.FromEnum(props.LifecycleState)
				state.LocalAdgAutoFailoverMaxDataLossLimit = pointer.From(props.LocalAdgAutoFailoverMaxDataLossLimit)
				state.LocalDataGuardEnabled = pointer.From(props.IsLocalDataGuardEnabled)
				state.LongTermBackupSchedule = FlattenLongTermBackUpScheduleDetails(props.LongTermBackupSchedule)
				state.MemoryAreaInGbs = pointer.From(props.InMemoryAreaInGbs)
				state.MemoryPerOracleComputeUnitInGbs = pointer.From(props.MemoryPerOracleComputeUnitInGbs)
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.NextLongTermBackupTimeStamp = pointer.From(props.NextLongTermBackupTimeStamp)
				state.OciUrl = pointer.From(props.OciURL)
				state.Ocid = pointer.From(props.Ocid)
				state.PeerDbId = pointer.From(props.PeerDbId)
				state.PeerDbIds = pointer.From(props.PeerDbIds)
				state.Preview = pointer.From(props.IsPreview)
				state.PreviewVersionWithServiceTermsAccepted = pointer.From(props.IsPreviewVersionWithServiceTermsAccepted)
				state.PrivateEndpoint = pointer.From(props.PrivateEndpoint)
				state.PrivateEndpointIp = pointer.From(props.PrivateEndpointIP)
				state.PrivateEndpointLabel = pointer.From(props.PrivateEndpointLabel)
				state.ProvisionableCPUs = pointer.From(props.ProvisionableCPUs)
				state.RemoteDataGuardEnabled = pointer.From(props.IsRemoteDataGuardEnabled)
				state.ServiceConsoleUrl = pointer.From(props.ServiceConsoleURL)
				state.SqlWebDeveloperUrl = pointer.From(props.SqlWebDeveloperURL)
				state.SubnetId = pointer.From(props.SubnetId)
				state.SupportedRegionsToCloneTo = pointer.From(props.SupportedRegionsToCloneTo)
				state.TimeCreated = pointer.From(props.TimeCreated)
				state.TimeDataGuardRoleChanged = pointer.From(props.TimeDataGuardRoleChanged)
				state.TimeDeletionOfFreeAutonomousDatabase = pointer.From(props.TimeDeletionOfFreeAutonomousDatabase)
				state.TimeLocalDataGuardEnabled = pointer.From(props.TimeLocalDataGuardEnabled)
				state.TimeMaintenanceBegin = pointer.From(props.TimeMaintenanceBegin)
				state.TimeMaintenanceEnd = pointer.From(props.TimeMaintenanceEnd)
				state.TimeOfLastFailover = pointer.From(props.TimeOfLastFailover)
				state.TimeOfLastRefresh = pointer.From(props.TimeOfLastRefresh)
				state.TimeOfLastRefreshPoint = pointer.From(props.TimeOfLastRefreshPoint)
				state.TimeOfLastSwitchover = pointer.From(props.TimeOfLastSwitchover)
				state.TimeReclamationOfFreeAutonomousDatabase = pointer.From(props.TimeReclamationOfFreeAutonomousDatabase)
				state.UsedDataStorageSizeInGbs = pointer.From(props.UsedDataStorageSizeInGbs)
				state.UsedDataStorageSizeInTbs = pointer.From(props.UsedDataStorageSizeInTbs)
				state.VnetId = pointer.From(props.VnetId)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
