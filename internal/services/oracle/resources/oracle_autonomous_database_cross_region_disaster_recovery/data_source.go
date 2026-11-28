// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_cross_region_disaster_recovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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

type AutonomousDatabaseCrossRegionDisasterRecoveryDataModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`

	ActualUsedDataStorageSizeInTb                     float64           `tfschema:"actual_used_data_storage_size_in_tb"`
	AllocatedStorageSizeInTb                          float64           `tfschema:"allocated_storage_size_in_tb"`
	AutoScalingEnabled                                bool              `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled                      bool              `tfschema:"auto_scaling_for_storage_enabled"`
	AvailableUpgradeVersions                          []string          `tfschema:"available_upgrade_versions"`
	BackupRetentionPeriodInDays                       int64             `tfschema:"backup_retention_period_in_days"`
	CharacterSet                                      string            `tfschema:"character_set"`
	ComputeCount                                      float64           `tfschema:"compute_count"`
	CpuCoreCount                                      int64             `tfschema:"cpu_core_count"`
	CustomerContacts                                  []string          `tfschema:"customer_contacts"`
	DataBaseType                                      string            `tfschema:"database_type"`
	DataStorageSizeInGb                               int64             `tfschema:"data_storage_size_in_gb"`
	DataStorageSizeInTb                               int64             `tfschema:"data_storage_size_in_tb"`
	DbVersion                                         string            `tfschema:"database_version"`
	DatabaseWorkload                                  string            `tfschema:"database_workload"`
	DisplayName                                       string            `tfschema:"display_name"`
	FailedDataRecoveryInSeconds                       int64             `tfschema:"failed_data_recovery_in_seconds"`
	MemoryAreaInGb                                    int64             `tfschema:"in_memory_area_in_gb"`
	LicenseModel                                      string            `tfschema:"license_model"`
	LifecycleDetails                                  string            `tfschema:"lifecycle_details"`
	LocalAdgAutoFailoverMaximumDataLossLimitInSeconds int64             `tfschema:"local_adg_auto_failover_maximum_data_loss_limit_in_seconds"`
	LocalDataGuardEnabled                             bool              `tfschema:"local_data_guard_enabled"`
	MemoryPerOracleComputeUnitInGb                    int64             `tfschema:"memory_per_oracle_compute_unit_in_gb"`
	MtlsConnectionRequired                            bool              `tfschema:"mtls_connection_required"`
	NcharacterSet                                     string            `tfschema:"national_character_set"`
	NextLongTermBackupTimeStampInUtc                  string            `tfschema:"next_long_term_backup_timestamp_in_utc"`
	OciUrl                                            string            `tfschema:"oci_url"`
	Ocid                                              string            `tfschema:"ocid"`
	PeerDbIds                                         []string          `tfschema:"peer_database_ids"`
	Preview                                           bool              `tfschema:"preview"`
	PreviewVersionWithServiceTermsAccepted            bool              `tfschema:"preview_version_with_service_terms_accepted"`
	PrivateEndpointIP                                 string            `tfschema:"private_endpoint_ip"`
	PrivateEndpointLabel                              string            `tfschema:"private_endpoint_label"`
	PrivateEndpoint                                   string            `tfschema:"private_endpoint_url"`
	ProvisionableCPUs                                 []int64           `tfschema:"provisionable_cpus"`
	RemoteDataGuardEnabled                            bool              `tfschema:"remote_data_guard_enabled"`
	RemoteDisasterRecoveryType                        string            `tfschema:"remote_disaster_recovery_type"`
	ReplicateAutomaticBackupsEnabled                  bool              `tfschema:"replicate_automatic_backups_enabled"`
	ServiceConsoleUrl                                 string            `tfschema:"service_console_url"`
	SourceType                                        string            `tfschema:"source_type"`
	SourceAutonomousDatabaseId                        string            `tfschema:"source_autonomous_database_id"`
	SourceLocation                                    string            `tfschema:"source_location"`
	SourceOcid                                        string            `tfschema:"source_ocid"`
	SqlWebDeveloperUrl                                string            `tfschema:"sql_web_developer_url"`
	SubnetId                                          string            `tfschema:"subnet_id"`
	TimeCreatedInUtc                                  string            `tfschema:"time_created_in_utc"`
	TimeDataGuardRoleChangedInUtc                     string            `tfschema:"time_data_guard_role_changed_in_utc"`
	TimeDeletionOfFreeAutonomousDatabaseInUtc         string            `tfschema:"time_deletion_of_free_autonomous_database_in_utc"`
	TimeLocalDataGuardEnabledInUtc                    string            `tfschema:"time_local_data_guard_enabled_in_utc"`
	TimeMaintenanceBeginInUtc                         string            `tfschema:"time_maintenance_begin_in_utc"`
	TimeMaintenanceEndInUtc                           string            `tfschema:"time_maintenance_end_in_utc"`
	TimeOfLastFailoverInUtc                           string            `tfschema:"time_of_last_failover_in_utc"`
	TimeOfLastRefreshInUtc                            string            `tfschema:"time_of_last_refresh_in_utc"`
	TimeOfLastRefreshPointInUtc                       string            `tfschema:"time_of_last_refresh_point_in_utc"`
	TimeOfLastSwitchoverInUtc                         string            `tfschema:"time_of_last_switchover_in_utc"`
	TimeReclamationOfFreeAutonomousDatabaseInUtc      string            `tfschema:"time_reclamation_of_free_autonomous_database_in_utc"`
	UsedDataStorageSizeInGb                           int64             `tfschema:"used_data_storage_size_in_gb"`
	UsedDataStorageSizeInTb                           int64             `tfschema:"used_data_storage_size_in_tb"`
	VnetId                                            string            `tfschema:"virtual_network_id"`
	Tags                                              map[string]string `tfschema:"tags"`
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Cross Region Disaster Recovery Autonomous Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"actual_used_data_storage_size_in_tb": {
			Description: "The current amount of storage in use for user and system data in terabytes.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"allocated_storage_size_in_tb": {
			Description: "The amount of storage currently allocated for the database tables and billed for, rounded up.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
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
			Description: "A list of Oracle Database versions available for a database upgrade.",
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
			Description: "The character set for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"compute_count": {
			Description: "The compute amount available to the database.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"cpu_core_count": {
			Description: "The number of CPU cores available to the database.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"customer_contacts": {
			Description: "A list of customer contact email addresses.",
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

		"database_type": {
			Description: "The database type.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"database_version": {
			Description: "The Oracle Database version for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"database_workload": {
			Description: "The Autonomous Database workload type.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The user-friendly name for the Cross Region Disaster Recovery Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"failed_data_recovery_in_seconds": {
			Description: "The number of seconds of data loss for Data Guard failover.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"in_memory_area_in_gb": {
			Description: "The area assigned to in-memory tables in the Autonomous Database.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"license_model": {
			Description: "The Oracle license model that applies to the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"lifecycle_details": {
			Description: "Information about the current lifecycle state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"local_adg_auto_failover_maximum_data_loss_limit_in_seconds": {
			Description: "The maximum data loss limit in seconds for local ADG automatic failover.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"local_data_guard_enabled": {
			Description: "Whether local (in-region) Data Guard is enabled.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
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
			Description: "The national character set for the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"next_long_term_backup_timestamp_in_utc": {
			Description: "The timestamp when the next long-term backup would be created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"oci_url": {
			Description: "The URL of the resource in the OCI console.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ocid": {
			Description: "The OCID of the Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"peer_database_ids": {
			Description: "The IDs of standby databases associated with the source database.",
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
			Description: "Whether the preview version terms were accepted.",
			Type:        pluginsdk.TypeBool,
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

		"private_endpoint_url": {
			Description: "The Private Endpoint URL for the resource.",
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
			Description: "Whether cross-region Data Guard is enabled.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"remote_disaster_recovery_type": {
			Description: "Type of recovery.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"replicate_automatic_backups_enabled": {
			Description: "If true, backups are replicated across regions.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"service_console_url": {
			Description: "The URL of the service console.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"source_type": {
			Description: "The source type used for the Cross Region Disaster Recovery Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"source_autonomous_database_id": {
			Description: "The ID of the source (primary) Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"source_location": {
			Description: "The Azure Region where the source Autonomous Database is located.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"source_ocid": {
			Description: "The OCID of the source Autonomous Database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"sql_web_developer_url": {
			Description: "The URL of SQL Web Developer.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subnet_id": {
			Description: "The ID of the subnet associated with this database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_created_in_utc": {
			Description: "The date and time the Autonomous Database was created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_data_guard_role_changed_in_utc": {
			Description: "The timestamp when the Data Guard role was last switched.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_deletion_of_free_autonomous_database_in_utc": {
			Description: "The timestamp when an inactive Always Free database is deleted.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_local_data_guard_enabled_in_utc": {
			Description: "The timestamp when local Data Guard was enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_maintenance_begin_in_utc": {
			Description: "The timestamp when maintenance begins.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_maintenance_end_in_utc": {
			Description: "The timestamp when maintenance ends.",
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
			Description: "The refresh point timestamp.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_of_last_switchover_in_utc": {
			Description: "The timestamp of the last switchover operation.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_reclamation_of_free_autonomous_database_in_utc": {
			Description: "The timestamp when an inactive Always Free database is stopped.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"used_data_storage_size_in_gb": {
			Description: "The amount of used data storage in gigabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"used_data_storage_size_in_tb": {
			Description: "The amount of used data storage in terabytes.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"virtual_network_id": {
			Description: "The ID of an Azure Resource Manager Virtual Network resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d DataSource) ModelObject() any {
	return &AutonomousDatabaseCrossRegionDisasterRecoveryDataModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Cross Region Disaster Recovery Autonomous Database."
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

			var state AutonomousDatabaseCrossRegionDisasterRecoveryDataModel
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

				adbsProps, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseCrossRegionDisasterRecoveryProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `CrossRegionDisasterRecovery`", id)
				}

				if adbsProps.RemoteDisasterRecoveryType != autonomousdatabases.DisasterRecoveryTypeAdg {
					return fmt.Errorf("%s expected disaster recovery type `%s`, got `%s`", id, autonomousdatabases.DisasterRecoveryTypeAdg, adbsProps.RemoteDisasterRecoveryType)
				}

				state.ReplicateAutomaticBackupsEnabled = pointer.From(adbsProps.IsReplicateAutomaticBackups)
				state.RemoteDisasterRecoveryType = string(adbsProps.RemoteDisasterRecoveryType)
				state.DataBaseType = string(adbsProps.DataBaseType)
				state.SourceType = string(adbsProps.Source)

				if adbsProps.SourceId != "" {
					sourceDatabaseID, err := autonomousdatabases.ParseAutonomousDatabaseIDInsensitively(adbsProps.SourceId)
					if err != nil {
						return err
					}
					state.SourceAutonomousDatabaseId = sourceDatabaseID.ID()
				}

				state.SourceLocation = location.NormalizeNilable(adbsProps.SourceLocation)
				state.SourceOcid = pointer.From(adbsProps.SourceOcid)
				state.ActualUsedDataStorageSizeInTb = pointer.From(adbsProps.ActualUsedDataStorageSizeInTbs)
				state.AllocatedStorageSizeInTb = pointer.From(adbsProps.AllocatedStorageSizeInTbs)
				state.AutoScalingEnabled = pointer.From(adbsProps.IsAutoScalingEnabled)
				state.AutoScalingForStorageEnabled = pointer.From(adbsProps.IsAutoScalingForStorageEnabled)
				state.AvailableUpgradeVersions = pointer.From(adbsProps.AvailableUpgradeVersions)
				state.BackupRetentionPeriodInDays = pointer.From(adbsProps.BackupRetentionPeriodInDays)
				state.CharacterSet = pointer.From(adbsProps.CharacterSet)
				state.ComputeCount = pointer.From(adbsProps.ComputeCount)
				state.CpuCoreCount = pointer.From(adbsProps.CpuCoreCount)
				state.CustomerContacts = helpers.FlattenAdbsCustomerContacts(adbsProps.CustomerContacts)
				state.DataStorageSizeInGb = pointer.From(adbsProps.DataStorageSizeInGbs)
				state.DataStorageSizeInTb = pointer.From(adbsProps.DataStorageSizeInTbs)
				state.DatabaseWorkload = pointer.FromEnum(adbsProps.DbWorkload)
				state.DbVersion = pointer.From(adbsProps.DbVersion)
				state.DisplayName = pointer.From(adbsProps.DisplayName)
				state.FailedDataRecoveryInSeconds = pointer.From(adbsProps.FailedDataRecoveryInSeconds)
				state.LifecycleDetails = pointer.From(adbsProps.LifecycleDetails)
				state.LocalAdgAutoFailoverMaximumDataLossLimitInSeconds = pointer.From(adbsProps.LocalAdgAutoFailoverMaxDataLossLimit)
				state.LocalDataGuardEnabled = pointer.From(adbsProps.IsLocalDataGuardEnabled)
				state.LicenseModel = pointer.FromEnum(adbsProps.LicenseModel)
				state.MemoryAreaInGb = pointer.From(adbsProps.InMemoryAreaInGbs)
				state.MemoryPerOracleComputeUnitInGb = pointer.From(adbsProps.MemoryPerOracleComputeUnitInGbs)
				state.MtlsConnectionRequired = pointer.From(adbsProps.IsMtlsConnectionRequired)
				state.NcharacterSet = pointer.From(adbsProps.NcharacterSet)
				state.NextLongTermBackupTimeStampInUtc = pointer.From(adbsProps.NextLongTermBackupTimeStamp)
				state.Ocid = pointer.From(adbsProps.Ocid)
				state.OciUrl = pointer.From(adbsProps.OciURL)
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

				if adbsProps.SubnetId != nil {
					subnetID, err := commonids.ParseSubnetIDInsensitively(*adbsProps.SubnetId)
					if err != nil {
						return err
					}
					state.SubnetId = subnetID.ID()
				}

				state.TimeCreatedInUtc = pointer.From(adbsProps.TimeCreated)
				state.TimeDataGuardRoleChangedInUtc = pointer.From(adbsProps.TimeDataGuardRoleChanged)
				state.TimeDeletionOfFreeAutonomousDatabaseInUtc = pointer.From(adbsProps.TimeDeletionOfFreeAutonomousDatabase)
				state.TimeLocalDataGuardEnabledInUtc = pointer.From(adbsProps.TimeLocalDataGuardEnabled)
				state.TimeMaintenanceBeginInUtc = pointer.From(adbsProps.TimeMaintenanceBegin)
				state.TimeMaintenanceEndInUtc = pointer.From(adbsProps.TimeMaintenanceEnd)
				state.TimeOfLastFailoverInUtc = pointer.From(adbsProps.TimeOfLastFailover)
				state.TimeOfLastRefreshInUtc = pointer.From(adbsProps.TimeOfLastRefresh)
				state.TimeOfLastRefreshPointInUtc = pointer.From(adbsProps.TimeOfLastRefreshPoint)
				state.TimeOfLastSwitchoverInUtc = pointer.From(adbsProps.TimeOfLastSwitchover)
				state.TimeReclamationOfFreeAutonomousDatabaseInUtc = pointer.From(adbsProps.TimeReclamationOfFreeAutonomousDatabase)
				state.UsedDataStorageSizeInGb = pointer.From(adbsProps.UsedDataStorageSizeInGbs)
				state.UsedDataStorageSizeInTb = pointer.From(adbsProps.UsedDataStorageSizeInTbs)

				if vnetIdRaw := pointer.From(adbsProps.VnetId); vnetIdRaw != "" {
					vnetId, err := commonids.ParseVirtualNetworkIDInsensitively(vnetIdRaw)
					if err != nil {
						return err
					}
					state.VnetId = vnetId.ID()
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
