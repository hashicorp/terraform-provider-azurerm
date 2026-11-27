// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_clone_from_backup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

type AutonomousDatabaseCloneFromBackupDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`

	SourceAutonomousDatabaseId string `tfschema:"source_autonomous_database_id"`

	// Base properties (computed)
	AllowedIpAddresses                            []string                                `tfschema:"allowed_ip_addresses"`
	BackupRetentionPeriodInDays                   int64                                   `tfschema:"backup_retention_period_in_days"`
	CharacterSet                                  string                                  `tfschema:"character_set"`
	ComputeCount                                  float64                                 `tfschema:"compute_count"`
	ComputeModel                                  string                                  `tfschema:"compute_model"`
	ConnectionStrings                             []string                                `tfschema:"connection_strings"`
	CustomerContacts                              []string                                `tfschema:"customer_contacts"`
	DataStorageSizeInGb                           int64                                   `tfschema:"data_storage_size_in_gb"`
	DataStorageSizeInTb                           int64                                   `tfschema:"data_storage_size_in_tb"`
	DatabaseVersion                               string                                  `tfschema:"database_version"`
	DatabaseWorkload                              string                                  `tfschema:"database_workload"`
	DisplayName                                   string                                  `tfschema:"display_name"`
	LicenseModel                                  string                                  `tfschema:"license_model"`
	AutoScalingEnabled                            bool                                    `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled                  bool                                    `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired                        bool                                    `tfschema:"mtls_connection_required"`
	NationalCharacterSet                          string                                  `tfschema:"national_character_set"`
	SubnetId                                      string                                  `tfschema:"subnet_id"`
	VnetId                                        string                                  `tfschema:"virtual_network_id"`
	LifecycleState                                string                                  `tfschema:"lifecycle_state"`
	PrivateEndpointUrl                            string                                  `tfschema:"private_endpoint_url"`
	PrivateEndpointIp                             string                                  `tfschema:"private_endpoint_ip"`
	ServiceConsoleUrl                             string                                  `tfschema:"service_console_url"`
	SqlWebDeveloperUrl                            string                                  `tfschema:"sql_web_developer_url"`
	TimeCreatedUtc                                string                                  `tfschema:"time_created_in_utc"`
	OciUrl                                        string                                  `tfschema:"oci_url"`
	ActualUsedDataStorageSizeInTb                 float64                                 `tfschema:"actual_used_data_storage_size_in_tb"`
	AllocatedStorageSizeInTb                      float64                                 `tfschema:"allocated_storage_size_in_tb"`
	AvailableUpgradeVersions                      []string                                `tfschema:"available_upgrade_versions"`
	CpuCoreCount                                  int64                                   `tfschema:"cpu_core_count"`
	FailedDataRecoveryInSeconds                   int64                                   `tfschema:"failed_data_recovery_in_seconds"`
	LifecycleDetails                              string                                  `tfschema:"lifecycle_details"`
	LocalAdgAutoFailoverMaxDataLossLimitInSeconds int64                                   `tfschema:"local_adg_auto_failover_max_data_loss_limit_in_seconds"`
	LocalDataGuardEnabled                         bool                                    `tfschema:"local_data_guard_enabled"`
	LongTermBackupSchedule                        []helpers.LongTermBackUpScheduleDetails `tfschema:"long_term_backup_schedule"`
	MemoryAreaInGb                                int64                                   `tfschema:"in_memory_area_in_gb"`
	MemoryPerOracleComputeUnitInGb                int64                                   `tfschema:"memory_per_oracle_compute_unit_in_gb"`
	NextLongTermBackupTimestamp                   string                                  `tfschema:"next_long_term_backup_timestamp"`
	Ocid                                          string                                  `tfschema:"ocid"`
	PeerDatabaseIds                               []string                                `tfschema:"peer_database_ids"`
	Preview                                       bool                                    `tfschema:"preview"`
	PreviewVersionWithServiceTermsAccepted        bool                                    `tfschema:"preview_version_with_service_terms_accepted"`
	PrivateEndpointLabel                          string                                  `tfschema:"private_endpoint_label"`
	ProvisionableCPUs                             []int64                                 `tfschema:"provisionable_cpus"`
	RemoteDataGuardEnabled                        bool                                    `tfschema:"remote_data_guard_enabled"`
	SupportedRegionsToCloneTo                     []string                                `tfschema:"supported_regions_to_clone_to"`
	TimeDataGuardRoleChangedInUtc                 string                                  `tfschema:"time_data_guard_role_changed_in_utc"`
	TimeDeletionOfFreeAutonomousDatabaseInUtc     string                                  `tfschema:"time_deletion_of_free_autonomous_database_in_utc"`
	TimeLocalDataGuardEnabledInUtc                string                                  `tfschema:"time_local_data_guard_enabled_in_utc"`
	TimeMaintenanceBeginInUtc                     string                                  `tfschema:"time_maintenance_begin_in_utc"`
	TimeMaintenanceEndInUtc                       string                                  `tfschema:"time_maintenance_end_in_utc"`
	TimeOfLastFailoverInUtc                       string                                  `tfschema:"time_of_last_failover_in_utc"`
	TimeOfLastRefreshInUtc                        string                                  `tfschema:"time_of_last_refresh_in_utc"`
	TimeOfLastRefreshPointInUtc                   string                                  `tfschema:"time_of_last_refresh_point_in_utc"`
	TimeOfLastSwitchoverInUtc                     string                                  `tfschema:"time_of_last_switchover_in_utc"`
	TimeReclamationOfFreeAutonomousDatabaseInUtc  string                                  `tfschema:"time_reclamation_of_free_autonomous_database_in_utc"`
	UsedDataStorageSizeInGb                       int64                                   `tfschema:"used_data_storage_size_in_gb"`
	UsedDataStorageSizeInTb                       int64                                   `tfschema:"used_data_storage_size_in_tb"`
}

func (DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this autonomous database clone from backup.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		"source_autonomous_database_id": {
			Description: "The ID of the source Autonomous Database from which this clone was created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"actual_used_data_storage_size_in_tb": {
			Description: "The current amount of storage in use for user and system data in terabytes.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"allocated_storage_size_in_tb": {
			Description: "The amount of storage currently allocated for the database tables and billed for, rounded up. When auto-scaling is not enabled, this value is equal to the `data_storage_size_in_tb` value. You can compare this value to the `actual_used_data_storage_size_in_tb` value to determine if a manual shrink operation is appropriate for your allocated storage.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"allowed_ip_addresses": {
			Description: "A list of IP addresses on the access control list.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auto_scaling_enabled": {
			Description: "Whether auto-scaling is enabled for the Autonomous Database CPU core count.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"auto_scaling_for_storage_enabled": {
			Description: "Whether auto-scaling is enabled for the Autonomous Database storage.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"available_upgrade_versions": {
			Description: "A list of Oracle Database versions available for a database upgrade. If there are no version upgrades available, this list is empty.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backup_retention_period_in_days": {
			Description: "The backup retention period in days.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"character_set": {
			Description: "The character set for the autonomous database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"compute_count": {
			Description: "The compute amount (CPUs) available to the database.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"compute_model": {
			Description: "The compute model of the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"connection_strings": {
			Description: "The connection string used to connect to the Autonomous Database.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cpu_core_count": {
			Description: "The number of CPU cores available to the database. When the ECPU is selected, the value for cpuCoreCount is 0.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"customer_contacts": {
			Description: "A list of Customer's contact email addresses.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"data_storage_size_in_gb": {
			Description: "The quantity of data in the database in gigabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"data_storage_size_in_tb": {
			Description: "The maximum storage that can be allocated for the database in terabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"database_version": {
			Description: "The Oracle Database version for Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"database_workload": {
			Description: "The Autonomous Database workload type.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The display name for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"failed_data_recovery_in_seconds": {
			Description: "The number of seconds of data loss for Data Guard failover.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"in_memory_area_in_gb": {
			Description: "The area assigned to In-Memory tables in Autonomous Database.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"license_model": {
			Description: "The Oracle license model that applied to the Oracle Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"lifecycle_details": {
			Description: "Information about the current lifecycle state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"lifecycle_state": {
			Description: "The current state of the backup.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"local_adg_auto_failover_max_data_loss_limit_in_seconds": {
			Description: "Parameter that allows users to select an acceptable maximum data loss limit in seconds, up to which Automatic Failover will be triggered when necessary for a Local Autonomous Data Guard.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"local_data_guard_enabled": {
			Description: "Whether the Autonomous Database has local (in-region) Data Guard enabled. Not applicable to cross-region Autonomous Data Guard associations, or to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"long_term_backup_schedule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"repeat_cadence": {
						Description: "The frequency for automated long-term backups.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"time_of_backup_in_utc": {
						Description: "The timestamp in which the backup would be made.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"retention_period_in_days": {
						Description: "The retention period in days for Autonomous database backup.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
					"enabled": {
						Description: "A boolean value that indicates if long-term backup is enabled/disabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
				},
			},
		},

		"memory_per_oracle_compute_unit_in_gb": {
			Description: "The amount of memory in gigabytes per ECPU or OCPU.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"mtls_connection_required": {
			Description: "Whether the Autonomous Database requires mTLS connections.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"national_character_set": {
			Description: "The national character set for the autonomous database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"next_long_term_backup_timestamp": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oci_url": {
			Description: "The URL of the resource in the OCI console.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ocid": {
			Description: "The [OCID](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"peer_database_ids": {
			Description: "The list of [OCIDs](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of standby databases located in Autonomous Data Guard remote regions that are associated with the source database. For Autonomous Database Serverless instances, standby databases located in the same region as the source primary database do not have OCIDs.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"preview": {
			Description: "Whether the Autonomous Database version is a preview version.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"preview_version_with_service_terms_accepted": {
			Description: "Whether the Autonomous Database version is a preview version with service terms accepted.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"private_endpoint_url": {
			Description: "The Private Endpoint for the resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_endpoint_ip": {
			Description: "The Private Endpoint IP address for the resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_endpoint_label": {
			Description: "The Private Endpoint label for the resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"provisionable_cpus": {
			Description: "An array of CPU values that an Autonomous Database can be scaled to.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeInt,
			},
		},

		"remote_data_guard_enabled": {
			Description: "Whether the Autonomous Database has Cross Region Data Guard enabled. Not applicable to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"service_console_url": {
			Description: "The URL of the Service Console for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"sql_web_developer_url": {
			Description: "The URL of the SQL web developer portal.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subnet_id": {
			Description: "The ID to an Azure Resource Manager subnet the resource is associated with.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"supported_regions_to_clone_to": {
			Description: "The list of regions that support the creation of an Autonomous Database clone or an Autonomous Data Guard standby database.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"time_created_in_utc": {
			Description: "The timestamp the Autonomous Database was created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_data_guard_role_changed_in_utc": {
			Description: "The timestamp the Autonomous Data Guard role was switched for the Autonomous Database. For databases that have standbys in both the primary Data Guard region and a remote Data Guard standby region, this is the latest timestamp of either the database using the \"primary\" role in the primary Data Guard region, or database located in the remote Data Guard standby region.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_deletion_of_free_autonomous_database_in_utc": {
			Description: "The timestamp the Always Free database will be automatically deleted because of inactivity. If the database is in the STOPPED state and without activity until this time, it will be deleted.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_local_data_guard_enabled_in_utc": {
			Description: "The timestamp that Autonomous Data Guard was enabled for an Autonomous Database where the standby was provisioned in the same region as the primary database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_maintenance_begin_in_utc": {
			Description: "The timestamp when maintenance will begin.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_maintenance_end_in_utc": {
			Description: "The timestamp when maintenance will end.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_failover_in_utc": {
			Description: "The timestamp of the last failover operation.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_refresh_in_utc": {
			Description: "The timestamp when the last refresh happened.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_refresh_point_in_utc": {
			Description: "The refresh point timestamp (UTC). The refresh point is the time to which the database was most recently refreshed. Data created after the refresh point is not included in the refresh.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_switchover_in_utc": {
			Description: "The timestamp of the last switchover operation for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_reclamation_of_free_autonomous_database_in_utc": {
			Description: "The timestamp the Always Free database will be stopped because of inactivity. If this time is reached without any database activity, the database will automatically be put into the STOPPED state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"used_data_storage_size_in_gb": {
			Description: "The storage space consumed by Autonomous Database in gigabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"used_data_storage_size_in_tb": {
			Description: "The amount of storage that has been used, in terabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"virtual_network_id": {
			Description: "The ID to an Azure Resource Manager Virtual Network resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (DataSource) ModelObject() any {
	return &AutonomousDatabaseCloneFromBackupDataSourceModel{}
}

func (DataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone_from_backup"
}

func (DataSource) Description() string {
	return "Use this data source to access information about an existing autonomous database clone from backup."
}

func (DataSource) Read() sdk.ResourceFunc {
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
				state.CustomerContacts = helpers.FlattenAdbsCustomerContacts(props.CustomerContacts)
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
				state.LongTermBackupSchedule = helpers.FlattenLongTermBackUpScheduleDetails(props.LongTermBackupSchedule)
				state.MemoryAreaInGb = pointer.From(props.InMemoryAreaInGbs)
				state.MemoryPerOracleComputeUnitInGb = pointer.From(props.MemoryPerOracleComputeUnitInGbs)
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.NextLongTermBackupTimestamp = pointer.From(props.NextLongTermBackupTimeStamp)
				state.OciUrl = pointer.From(props.OciURL)
				state.Ocid = pointer.From(props.Ocid)
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
				state.ConnectionStrings = helpers.FlattenConnectionStrings(props.ConnectionStrings)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
