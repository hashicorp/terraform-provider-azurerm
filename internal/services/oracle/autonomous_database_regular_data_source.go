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
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseRegularDataSource struct{}

type AutonomousDatabaseRegularDataModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// AutonomousDatabaseProperties
	ActualUsedDataStorageSizeInTbs          float64  `tfschema:"actual_used_data_storage_size_in_tbs"`
	AllocatedStorageSizeInTbs               float64  `tfschema:"allocated_storage_size_in_tbs"`
	AutonomousDatabaseId                    string   `tfschema:"autonomous_database_id"`
	AutoScalingEnabled                      bool     `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled            bool     `tfschema:"auto_scaling_for_storage_enabled"`
	AvailableUpgradeVersions                []string `tfschema:"available_upgrade_versions"`
	BackupRetentionPeriodInDays             int64    `tfschema:"backup_retention_period_in_days"`
	CharacterSet                            string   `tfschema:"character_set"`
	ComputeCount                            float64  `tfschema:"compute_count"`
	CpuCoreCount                            int64    `tfschema:"cpu_core_count"`
	DataStorageSizeInGbs                    int64    `tfschema:"data_storage_size_in_gbs"`
	DataStorageSizeInTbs                    int64    `tfschema:"data_storage_size_in_tbs"`
	DbVersion                               string   `tfschema:"db_version"`
	DisplayName                             string   `tfschema:"display_name"`
	FailedDataRecoveryInSeconds             int64    `tfschema:"failed_data_recovery_in_seconds"`
	LifecycleDetails                        string   `tfschema:"lifecycle_details"`
	LocalAdgAutoFailoverMaxDataLossLimit    int64    `tfschema:"local_adg_auto_failover_max_data_loss_limit"`
	LocalDataGuardEnabled                   bool     `tfschema:"local_data_guard_enabled"`
	MemoryAreaInGbs                         int64    `tfschema:"in_memory_area_in_gbs"`
	MemoryPerOracleComputeUnitInGbs         int64    `tfschema:"memory_per_oracle_compute_unit_in_gbs"`
	MtlsConnectionRequired                  bool     `tfschema:"mtls_connection_required"`
	NcharacterSet                           string   `tfschema:"national_character_set"`
	NextLongTermBackupTimeStamp             string   `tfschema:"next_long_term_backup_time_stamp"`
	Ocid                                    string   `tfschema:"ocid"`
	OciUrl                                  string   `tfschema:"oci_url"`
	PeerDbId                                string   `tfschema:"peer_db_id"`
	PeerDbIds                               []string `tfschema:"peer_db_ids"`
	Preview                                 bool     `tfschema:"preview"`
	PreviewVersionWithServiceTermsAccepted  bool     `tfschema:"preview_version_with_service_terms_accepted"`
	PrivateEndpoint                         string   `tfschema:"private_endpoint"`
	PrivateEndpointIP                       string   `tfschema:"private_endpoint_ip"`
	PrivateEndpointLabel                    string   `tfschema:"private_endpoint_label"`
	ProvisionableCPUs                       []int64  `tfschema:"provisionable_cpus"`
	RemoteDataGuardEnabled                  bool     `tfschema:"remote_data_guard_enabled"`
	ServiceConsoleUrl                       string   `tfschema:"service_console_url"`
	SqlWebDeveloperUrl                      string   `tfschema:"sql_web_developer_url"`
	SubnetId                                string   `tfschema:"subnet_id"`
	SupportedRegionsToCloneTo               []string `tfschema:"supported_regions_to_clone_to"`
	TimeCreated                             string   `tfschema:"time_created"`
	TimeDataGuardRoleChanged                string   `tfschema:"time_data_guard_role_changed"`
	TimeDeletionOfFreeAutonomousDatabase    string   `tfschema:"time_deletion_of_free_autonomous_database"`
	TimeLocalDataGuardEnabled               string   `tfschema:"time_local_data_guard_enabled_on"`
	TimeMaintenanceBegin                    string   `tfschema:"time_maintenance_begin"`
	TimeMaintenanceEnd                      string   `tfschema:"time_maintenance_end"`
	TimeOfLastFailover                      string   `tfschema:"time_of_last_failover"`
	TimeOfLastRefresh                       string   `tfschema:"time_of_last_refresh"`
	TimeOfLastRefreshPoint                  string   `tfschema:"time_of_last_refresh_point"`
	TimeOfLastSwitchover                    string   `tfschema:"time_of_last_switchover"`
	TimeReclamationOfFreeAutonomousDatabase string   `tfschema:"time_reclamation_of_free_autonomous_database"`
	UsedDataStorageSizeInGbs                int64    `tfschema:"used_data_storage_size_in_gbs"`
	UsedDataStorageSizeInTbs                int64    `tfschema:"used_data_storage_size_in_tbs"`
	VnetId                                  string   `tfschema:"virtual_network_id"`
	AllowedIps                              []string `tfschema:"allowed_ips"`
}

func (d AutonomousDatabaseRegularDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},
	}
}

func (d AutonomousDatabaseRegularDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// AutonomousDatabaseProperties
		"actual_used_data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"allocated_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"autonomous_database_id": {
			Type:     pluginsdk.TypeString,
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

		"cpu_core_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"data_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"db_node_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"db_version": {
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

		"lifecycle_details": {
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

		"preview": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"preview_version_with_service_terms_accepted": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"remote_data_guard_enabled": {
			Type:     pluginsdk.TypeBool,
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

		"allowed_ips": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeInt,
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d AutonomousDatabaseRegularDataSource) ModelObject() interface{} {
	return &AutonomousDatabaseRegularDataModel{}
}

func (d AutonomousDatabaseRegularDataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database"
}

func (d AutonomousDatabaseRegularDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}

func (d AutonomousDatabaseRegularDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state AutonomousDatabaseRegularDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					adbsProps := props.AutonomousDatabaseBaseProperties()

					state.ActualUsedDataStorageSizeInTbs = pointer.From(adbsProps.ActualUsedDataStorageSizeInTbs)
					state.AllocatedStorageSizeInTbs = pointer.From(adbsProps.AllocatedStorageSizeInTbs)
					state.AutonomousDatabaseId = pointer.From(adbsProps.AutonomousDatabaseId)
					state.AutoScalingEnabled = pointer.From(adbsProps.IsAutoScalingEnabled)
					state.AutoScalingForStorageEnabled = pointer.From(adbsProps.IsAutoScalingForStorageEnabled)
					state.AvailableUpgradeVersions = pointer.From(adbsProps.AvailableUpgradeVersions)
					state.BackupRetentionPeriodInDays = pointer.From(adbsProps.BackupRetentionPeriodInDays)
					state.CharacterSet = pointer.From(adbsProps.CharacterSet)
					state.ComputeCount = pointer.From(adbsProps.ComputeCount)
					state.CpuCoreCount = pointer.From(adbsProps.CpuCoreCount)
					state.DataStorageSizeInGbs = pointer.From(adbsProps.DataStorageSizeInGbs)
					state.DataStorageSizeInTbs = pointer.From(adbsProps.DataStorageSizeInTbs)
					state.DbVersion = pointer.From(adbsProps.DbVersion)
					state.DisplayName = pointer.From(adbsProps.DisplayName)
					state.FailedDataRecoveryInSeconds = pointer.From(adbsProps.FailedDataRecoveryInSeconds)
					state.LifecycleDetails = pointer.From(adbsProps.LifecycleDetails)
					state.LocalAdgAutoFailoverMaxDataLossLimit = pointer.From(adbsProps.LocalAdgAutoFailoverMaxDataLossLimit)
					state.LocalDataGuardEnabled = pointer.From(adbsProps.IsLocalDataGuardEnabled)
					state.MemoryAreaInGbs = pointer.From(adbsProps.InMemoryAreaInGbs)
					state.MemoryPerOracleComputeUnitInGbs = pointer.From(adbsProps.MemoryPerOracleComputeUnitInGbs)
					state.MtlsConnectionRequired = pointer.From(adbsProps.IsMtlsConnectionRequired)
					state.NcharacterSet = pointer.From(adbsProps.NcharacterSet)
					state.NextLongTermBackupTimeStamp = pointer.From(adbsProps.NextLongTermBackupTimeStamp)
					state.Ocid = pointer.From(adbsProps.Ocid)
					state.OciUrl = pointer.From(adbsProps.OciURL)
					state.PeerDbId = pointer.From(adbsProps.PeerDbId)
					state.PeerDbIds = pointer.From(adbsProps.PeerDbIds)
					state.Preview = pointer.From(adbsProps.IsPreview)
					state.PreviewVersionWithServiceTermsAccepted = pointer.From(adbsProps.IsPreviewVersionWithServiceTermsAccepted)
					state.PrivateEndpoint = pointer.From(adbsProps.PrivateEndpoint)
					state.PrivateEndpointIP = pointer.From(adbsProps.PrivateEndpointIP)
					state.PrivateEndpointLabel = pointer.From(adbsProps.PrivateEndpointLabel)
					state.ProvisionableCPUs = pointer.From(adbsProps.ProvisionableCPUs)
					state.RemoteDataGuardEnabled = pointer.From(adbsProps.IsRemoteDataGuardEnabled)
					state.ServiceConsoleUrl = pointer.From(adbsProps.ServiceConsoleURL)
					state.SqlWebDeveloperUrl = pointer.From(adbsProps.SqlWebDeveloperURL)
					state.SubnetId = pointer.From(adbsProps.SubnetId)
					state.SupportedRegionsToCloneTo = pointer.From(adbsProps.SupportedRegionsToCloneTo)
					state.TimeCreated = pointer.From(adbsProps.TimeCreated)
					state.TimeDataGuardRoleChanged = pointer.From(adbsProps.TimeDataGuardRoleChanged)
					state.TimeDeletionOfFreeAutonomousDatabase = pointer.From(adbsProps.TimeDeletionOfFreeAutonomousDatabase)
					state.TimeLocalDataGuardEnabled = pointer.From(adbsProps.TimeLocalDataGuardEnabled)
					state.TimeMaintenanceBegin = pointer.From(adbsProps.TimeMaintenanceBegin)
					state.TimeMaintenanceEnd = pointer.From(adbsProps.TimeMaintenanceEnd)
					state.TimeOfLastFailover = pointer.From(adbsProps.TimeOfLastFailover)
					state.TimeOfLastRefresh = pointer.From(adbsProps.TimeOfLastRefresh)
					state.TimeOfLastRefreshPoint = pointer.From(adbsProps.TimeOfLastRefreshPoint)
					state.TimeOfLastSwitchover = pointer.From(adbsProps.TimeOfLastSwitchover)
					state.TimeReclamationOfFreeAutonomousDatabase = pointer.From(adbsProps.TimeReclamationOfFreeAutonomousDatabase)
					state.UsedDataStorageSizeInGbs = pointer.From(adbsProps.UsedDataStorageSizeInGbs)
					state.UsedDataStorageSizeInTbs = pointer.From(adbsProps.UsedDataStorageSizeInTbs)
					state.VnetId = pointer.From(adbsProps.VnetId)
					state.AllowedIps = pointer.From(adbsProps.WhitelistedIPs)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
