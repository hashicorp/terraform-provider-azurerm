// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type AutonomousDatabaseRegularDataModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// AutonomousDatabaseProperties
	ActualUsedDataStorageSizeInTbs          float64                                 `tfschema:"actual_used_data_storage_size_in_tbs"`
	AllocatedStorageSizeInTbs               float64                                 `tfschema:"allocated_storage_size_in_tbs"`
	AutonomousDatabaseId                    string                                  `tfschema:"autonomous_database_id"`
	AutoScalingEnabled                      bool                                    `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled            bool                                    `tfschema:"auto_scaling_for_storage_enabled"`
	AvailableUpgradeVersions                []string                                `tfschema:"available_upgrade_versions"`
	BackupRetentionPeriodInDays             int64                                   `tfschema:"backup_retention_period_in_days"`
	CharacterSet                            string                                  `tfschema:"character_set"`
	ComputeCount                            float64                                 `tfschema:"compute_count"`
	CpuCoreCount                            int64                                   `tfschema:"cpu_core_count"`
	DataStorageSizeInGbs                    int64                                   `tfschema:"data_storage_size_in_gbs"`
	DataStorageSizeInTbs                    int64                                   `tfschema:"data_storage_size_in_tbs"`
	DbVersion                               string                                  `tfschema:"db_version"`
	DisplayName                             string                                  `tfschema:"display_name"`
	FailedDataRecoveryInSeconds             int64                                   `tfschema:"failed_data_recovery_in_seconds"`
	LifecycleDetails                        string                                  `tfschema:"lifecycle_details"`
	LocalAdgAutoFailoverMaxDataLossLimit    int64                                   `tfschema:"local_adg_auto_failover_max_data_loss_limit"`
	LocalDataGuardEnabled                   bool                                    `tfschema:"local_data_guard_enabled"`
	LongTermBackupSchedule                  []helpers.LongTermBackUpScheduleDetails `tfschema:"long_term_backup_schedule"`
	MemoryAreaInGbs                         int64                                   `tfschema:"in_memory_area_in_gbs"`
	MemoryPerOracleComputeUnitInGbs         int64                                   `tfschema:"memory_per_oracle_compute_unit_in_gbs"`
	MtlsConnectionRequired                  bool                                    `tfschema:"mtls_connection_required"`
	NcharacterSet                           string                                  `tfschema:"national_character_set"`
	NextLongTermBackupTimeStamp             string                                  `tfschema:"next_long_term_backup_time_stamp"`
	Ocid                                    string                                  `tfschema:"ocid"`
	OciUrl                                  string                                  `tfschema:"oci_url"`
	PeerDbId                                string                                  `tfschema:"peer_db_id"`
	PeerDbIds                               []string                                `tfschema:"peer_db_ids"`
	Preview                                 bool                                    `tfschema:"preview"`
	PreviewVersionWithServiceTermsAccepted  bool                                    `tfschema:"preview_version_with_service_terms_accepted"`
	PrivateEndpoint                         string                                  `tfschema:"private_endpoint"`
	PrivateEndpointIP                       string                                  `tfschema:"private_endpoint_ip"`
	PrivateEndpointLabel                    string                                  `tfschema:"private_endpoint_label"`
	ProvisionableCPUs                       []int64                                 `tfschema:"provisionable_cpus"`
	RemoteDataGuardEnabled                  bool                                    `tfschema:"remote_data_guard_enabled"`
	ServiceConsoleUrl                       string                                  `tfschema:"service_console_url"`
	SqlWebDeveloperUrl                      string                                  `tfschema:"sql_web_developer_url"`
	SubnetId                                string                                  `tfschema:"subnet_id"`
	SupportedRegionsToCloneTo               []string                                `tfschema:"supported_regions_to_clone_to"`
	TimeCreated                             string                                  `tfschema:"time_created"`
	TimeDataGuardRoleChanged                string                                  `tfschema:"time_data_guard_role_changed"`
	TimeDeletionOfFreeAutonomousDatabase    string                                  `tfschema:"time_deletion_of_free_autonomous_database"`
	TimeLocalDataGuardEnabled               string                                  `tfschema:"time_local_data_guard_enabled_on"`
	TimeMaintenanceBegin                    string                                  `tfschema:"time_maintenance_begin"`
	TimeMaintenanceEnd                      string                                  `tfschema:"time_maintenance_end"`
	TimeOfLastFailover                      string                                  `tfschema:"time_of_last_failover"`
	TimeOfLastRefresh                       string                                  `tfschema:"time_of_last_refresh"`
	TimeOfLastRefreshPoint                  string                                  `tfschema:"time_of_last_refresh_point"`
	TimeOfLastSwitchover                    string                                  `tfschema:"time_of_last_switchover"`
	TimeReclamationOfFreeAutonomousDatabase string                                  `tfschema:"time_reclamation_of_free_autonomous_database"`
	UsedDataStorageSizeInGbs                int64                                   `tfschema:"used_data_storage_size_in_gbs"`
	UsedDataStorageSizeInTbs                int64                                   `tfschema:"used_data_storage_size_in_tbs"`
	VnetId                                  string                                  `tfschema:"virtual_network_id"`
	AllowedIps                              []string                                `tfschema:"allowed_ips"`
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Description:  "The name of this Autonomous Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// AutonomousDatabaseProperties
		"actual_used_data_storage_size_in_tbs": {
			Description: "The current amount of storage in use for user and system data, in terabytes (TB).",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"allocated_storage_size_in_tbs": {
			Description: "The amount of storage currently allocated for the database tables and billed for, rounded up. When auto-scaling is not enabled, this value is equal to the `dataStorageSizeInTBs` value. You can compare this value to the `actualUsedDataStorageSizeInTBs` value to determine if a manual shrink operation is appropriate for your allocated storage.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"auto_scaling_enabled": {
			Description: "Whether auto scaling is enabled for the Autonomous Database CPU core count.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"auto_scaling_for_storage_enabled": {
			Description: "Whether auto scaling is enabled for the Autonomous Database storage.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"autonomous_database_id": {
			Description: "The database [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"available_upgrade_versions": {
			Description: "List of Oracle Database versions available for a database upgrade. If there are no version upgrades available, this list is empty.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backup_retention_period_in_days": {
			Description: "Retention period, in days, for backups.",
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

		"cpu_core_count": {
			Description: "The number of CPU cores to be made available to the database. When the ECPU is selected, the value for cpuCoreCount is 0. For Autonomous Database on Dedicated Exadata infrastructure, the maximum number of cores is determined by the infrastructure shape. See [Characteristics of Infrastructure Shapes](https://docs.oracle.com/en/cloud/paas/autonomous-database/dedicated/adbaa/characteristics-of-infrastructure-shapes.html) for shape details.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"data_storage_size_in_gbs": {
			Description: "The quantity of data in the database, in gigabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"data_storage_size_in_tbs": {
			Description: "The maximum storage that can be allocated for the database, in terabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"db_node_storage_size_in_gbs": {
			Description: "The DB node storage size, in gigabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"db_version": {
			Description: "A valid Oracle Database version for Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The user-friendly name for the Autonomous Database. The name does not have to be unique.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"failed_data_recovery_in_seconds": {
			Description: "The number of seconds of data loss for a Data Guard failover.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"in_memory_area_in_gbs": {
			Description: "The area assigned to In-Memory tables in Autonomous Database.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"lifecycle_details": {
			Description: "Information about the current lifecycle state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"local_adg_auto_failover_max_data_loss_limit": {
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
					"time_of_backup": {
						Description: "The date and time in which the backup would be made.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"retention_period_in_days": {
						Description: "The retention period in days for Autonomous database backup.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
					"enabled": {
						Description: "A boolean value that indicates if long term backup is enabled/disabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
				},
			},
		},

		"memory_per_oracle_compute_unit_in_gbs": {
			Description: "The amount of memory (in GBs) enabled per ECPU or OCPU.",
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

		"next_long_term_backup_time_stamp": {
			Description: "The date and time when the next long-term backup would be created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"oci_url": {
			Description: "The URL of the resource in the OCI console.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ocid": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
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

		"remote_data_guard_enabled": {
			Description: "Whether the Autonomous Database has Cross Region Data Guard enabled. Not applicable to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"peer_db_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"peer_db_ids": {
			Description: "The list of [OCIDs](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of standby databases located in Autonomous Data Guard remote regions that are associated with the source database. Note that for Autonomous Database Serverless instances, standby databases located in the same region as the source primary database do not have OCIDs.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_endpoint": {
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

		"service_console_url": {
			Description: "The URL of the Service Console for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"sql_web_developer_url": {
			Description: "The URL of the SQL web developer.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subnet_id": {
			Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the subnet the resource is associated with.",
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

		"time_created": {
			Description: "The date and time the Autonomous Database was created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_data_guard_role_changed": {
			Description: "The date and time the Autonomous Data Guard role was switched for the Autonomous Database. For databases that have standbys in both the primary Data Guard region and a remote Data Guard standby region, this is the latest timestamp of either the database using the \"primary\" role in the primary Data Guard region, or database located in the remote Data Guard standby region.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_deletion_of_free_autonomous_database": {
			Description: "The date and time the Always Free database will be automatically deleted because of inactivity. If the database is in the STOPPED state and without activity until this time, it will be deleted.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_local_data_guard_enabled_on": {
			Description: "The date and time that Autonomous Data Guard was enabled for an Autonomous Database where the standby was provisioned in the same region as the primary database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_maintenance_begin": {
			Description: "The date and time when maintenance will begin.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_maintenance_end": {
			Description: "The date and time when maintenance will end.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_failover": {
			Description: "The timestamp of the last failover operation.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_refresh": {
			Description: "The date and time when last refresh happened.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_refresh_point": {
			Description: "The refresh point timestamp (UTC). The refresh point is the time to which the database was most recently refreshed. Data created after the refresh point is not included in the refresh.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_switchover": {
			Description: "The timestamp of the last switchover operation for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_reclamation_of_free_autonomous_database": {
			Description: "The date and time the Always Free database will be stopped because of inactivity. If this time is reached without any database activity, the database will automatically be put into the STOPPED state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"used_data_storage_size_in_gbs": {
			Description: "The storage space consumed by Autonomous Database in GBs.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"used_data_storage_size_in_tbs": {
			Description: "The amount of storage that has been used, in terabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"virtual_network_id": {
			Description: "The ID to an Azure Resource Manager VNet resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"allowed_ips": {
			Description: "The client IP access control list (ACL). This feature is available for [Autonomous Database Serverless](https://docs.oracle.com/en/cloud/paas/autonomous-database/index.html) and on Exadata Cloud@Customer. Only clients connecting from an IP address included in the ACL may access the Autonomous Database instance. If `arePrimaryWhitelistedIpsUsed` is 'TRUE' then Autonomous Database uses this primary's IP access control list (ACL) for the disaster recovery peer called `standbywhitelistedips`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeInt,
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d DataSource) ModelObject() any {
	return &AutonomousDatabaseRegularDataModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Autonomous Database."
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}

func (d DataSource) Read() sdk.ResourceFunc {
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
					state.LongTermBackupSchedule = helpers.FlattenLongTermBackUpScheduleDetails(adbsProps.LongTermBackupSchedule)
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
