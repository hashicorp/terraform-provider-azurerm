// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseRegularDataSource struct{}

type AutonomousDatabaseRegularModel struct {
	Location string                 `tfschema:"location"`
	Name     string                 `tfschema:"name"`
	Type     string                 `tfschema:"type"`
	Tags     map[string]interface{} `tfschema:"tags"`

	// SystemData
	SystemData []SystemDataModel `tfschema:"system_data"`

	// AutonomousDatabaseProperties
	ActualUsedDataStorageSizeInTbs           float64  `tfschema:"actual_used_data_storage_size_in_tbs"`
	AllocatedStorageSizeInTbs                float64  `tfschema:"allocated_storage_size_in_tbs"`
	AutonomousDatabaseId                     string   `tfschema:"autonomous_database_id"`
	AvailableUpgradeVersions                 []string `tfschema:"available_upgrade_versions"`
	BackupRetentionPeriodInDays              int64    `tfschema:"backup_retention_period_in_days"`
	CharacterSet                             string   `tfschema:"character_set"`
	ComputeCount                             float64  `tfschema:"compute_count"`
	CpuCoreCount                             int64    `tfschema:"cpu_core_count"`
	DataStorageSizeInGbs                     int64    `tfschema:"data_storage_size_in_gbs"`
	DataStorageSizeInTbs                     int64    `tfschema:"data_storage_size_in_tbs"`
	DbVersion                                string   `tfschema:"db_version"`
	DisplayName                              string   `tfschema:"display_name"`
	FailedDataRecoveryInSeconds              int64    `tfschema:"failed_data_recovery_in_seconds"`
	InMemoryAreaInGbs                        int64    `tfschema:"in_memory_area_in_gbs"`
	IsAutoScalingEnabled                     bool     `tfschema:"is_auto_scaling_enabled"`
	IsAutoScalingForStorageEnabled           bool     `tfschema:"is_auto_scaling_for_storage_enabled"`
	IsLocalDataGuardEnabled                  bool     `tfschema:"is_local_data_guard_enabled"`
	IsMtlsConnectionRequired                 bool     `tfschema:"is_mtls_connection_required"`
	IsPreview                                bool     `tfschema:"is_preview"`
	IsPreviewVersionWithServiceTermsAccepted bool     `tfschema:"is_preview_version_with_service_terms_accepted"`
	IsRemoteDataGuardEnabled                 bool     `tfschema:"is_remote_data_guard_enabled"`
	LifecycleDetails                         string   `tfschema:"lifecycle_details"`
	LocalAdgAutoFailoverMaxDataLossLimit     int64    `tfschema:"local_adg_auto_failover_max_data_loss_limit"`
	MemoryPerOracleComputeUnitInGbs          int64    `tfschema:"memory_per_oracle_compute_unit_in_gbs"`
	NcharacterSet                            string   `tfschema:"ncharacter_set"`
	NextLongTermBackupTimeStamp              string   `tfschema:"next_long_term_backup_time_stamp"`
	OciUrl                                   string   `tfschema:"oci_url"`
	Ocid                                     string   `tfschema:"ocid"`
	PeerDbId                                 string   `tfschema:"peer_db_id"`
	PeerDbIds                                []string `tfschema:"peer_db_ids"`
	PrivateEndpoint                          string   `tfschema:"private_endpoint"`
	PrivateEndpointIP                        string   `tfschema:"private_endpoint_ip"`
	PrivateEndpointLabel                     string   `tfschema:"private_endpoint_label"`
	ProvisionableCPUs                        []int64  `tfschema:"provisionable_cpus"`
	ServiceConsoleUrl                        string   `tfschema:"service_console_url"`
	SqlWebDeveloperUrl                       string   `tfschema:"sql_web_developer_url"`
	SubnetId                                 string   `tfschema:"subnet_id"`
	SupportedRegionsToCloneTo                []string `tfschema:"supported_regions_to_clone_to"`
	TimeCreated                              string   `tfschema:"time_created"`
	TimeDataGuardRoleChanged                 string   `tfschema:"time_data_guard_role_changed"`
	TimeDeletionOfFreeAutonomousDatabase     string   `tfschema:"time_deletion_of_free_autonomous_database"`
	TimeLocalDataGuardEnabled                string   `tfschema:"time_local_data_guard_enabled"`
	TimeMaintenanceBegin                     string   `tfschema:"time_maintenance_begin"`
	TimeMaintenanceEnd                       string   `tfschema:"time_maintenance_end"`
	TimeOfLastFailover                       string   `tfschema:"time_of_last_failover"`
	TimeOfLastRefresh                        string   `tfschema:"time_of_last_refresh"`
	TimeOfLastRefreshPoint                   string   `tfschema:"time_of_last_refresh_point"`
	TimeOfLastSwitchover                     string   `tfschema:"time_of_last_switchover"`
	TimeReclamationOfFreeAutonomousDatabase  string   `tfschema:"time_reclamation_of_free_autonomous_database"`
	UsedDataStorageSizeInGbs                 int64    `tfschema:"used_data_storage_size_in_gbs"`
	UsedDataStorageSizeInTbs                 int64    `tfschema:"used_data_storage_size_in_tbs"`
	VnetId                                   string   `tfschema:"vnet_id"`
	WhitelistedIPs                           []string `tfschema:"whitelisted_ips"`
}

func (d AutonomousDatabaseRegularDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d AutonomousDatabaseRegularDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"tags": commonschema.TagsDataSource(),

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

		// AutonomousDatabaseProperties
		"actual_used_data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},
		"allocated_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
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
		"db_node_storage_size_in_gbs": {
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
		"is_auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"is_auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"is_local_data_guard_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"is_mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"is_preview": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"is_preview_version_with_service_terms_accepted": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"is_remote_data_guard_enabled": {
			Type:     pluginsdk.TypeBool,
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
		"memory_per_oracle_compute_unit_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"ncharacter_set": {
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
				Type: pluginsdk.TypeInt,
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
		"time_local_data_guard_enabled": {
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
		"vnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"whitelisted_ips": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeInt,
			},
		},
	}
}

func (d AutonomousDatabaseRegularDataSource) ModelObject() interface{} {
	return nil
}

func (d AutonomousDatabaseRegularDataSource) ResourceType() string {
	return "azurerm_oracledatabase_autonomous_database_regular"
}

func (d AutonomousDatabaseRegularDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}

func (d AutonomousDatabaseRegularDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId,
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

				var output AutonomousDatabaseRegularModel
				prop := model.Properties

				output = AutonomousDatabaseRegularModel{}

				if prop != nil {
					switch adbsPropModel := prop.(type) {
					case autonomousdatabases.AutonomousDatabaseProperties:
						output = AutonomousDatabaseRegularModel{
							ActualUsedDataStorageSizeInTbs:           pointer.From(adbsPropModel.ActualUsedDataStorageSizeInTbs),
							AllocatedStorageSizeInTbs:                pointer.From(adbsPropModel.AllocatedStorageSizeInTbs),
							AutonomousDatabaseId:                     pointer.From(adbsPropModel.AutonomousDatabaseId),
							AvailableUpgradeVersions:                 pointer.From(adbsPropModel.AvailableUpgradeVersions),
							BackupRetentionPeriodInDays:              pointer.From(adbsPropModel.BackupRetentionPeriodInDays),
							CharacterSet:                             pointer.From(adbsPropModel.CharacterSet),
							ComputeCount:                             pointer.From(adbsPropModel.ComputeCount),
							CpuCoreCount:                             pointer.From(adbsPropModel.CpuCoreCount),
							DataStorageSizeInGbs:                     pointer.From(adbsPropModel.DataStorageSizeInGbs),
							DataStorageSizeInTbs:                     pointer.From(adbsPropModel.DataStorageSizeInTbs),
							DbVersion:                                pointer.From(adbsPropModel.DbVersion),
							DisplayName:                              pointer.From(adbsPropModel.DisplayName),
							FailedDataRecoveryInSeconds:              pointer.From(adbsPropModel.FailedDataRecoveryInSeconds),
							InMemoryAreaInGbs:                        pointer.From(adbsPropModel.InMemoryAreaInGbs),
							IsAutoScalingEnabled:                     pointer.From(adbsPropModel.IsAutoScalingEnabled),
							IsAutoScalingForStorageEnabled:           pointer.From(adbsPropModel.IsAutoScalingForStorageEnabled),
							IsLocalDataGuardEnabled:                  pointer.From(adbsPropModel.IsLocalDataGuardEnabled),
							IsMtlsConnectionRequired:                 pointer.From(adbsPropModel.IsMtlsConnectionRequired),
							IsPreview:                                pointer.From(adbsPropModel.IsPreview),
							IsPreviewVersionWithServiceTermsAccepted: pointer.From(adbsPropModel.IsPreviewVersionWithServiceTermsAccepted),
							IsRemoteDataGuardEnabled:                 pointer.From(adbsPropModel.IsRemoteDataGuardEnabled),
							LifecycleDetails:                         pointer.From(adbsPropModel.LifecycleDetails),
							LocalAdgAutoFailoverMaxDataLossLimit:     pointer.From(adbsPropModel.LocalAdgAutoFailoverMaxDataLossLimit),
							MemoryPerOracleComputeUnitInGbs:          pointer.From(adbsPropModel.MemoryPerOracleComputeUnitInGbs),
							NcharacterSet:                            pointer.From(adbsPropModel.NcharacterSet),
							NextLongTermBackupTimeStamp:              pointer.From(adbsPropModel.NextLongTermBackupTimeStamp),
							OciUrl:                                   pointer.From(adbsPropModel.OciUrl),
							Ocid:                                     pointer.From(adbsPropModel.Ocid),
							PeerDbId:                                 pointer.From(adbsPropModel.PeerDbId),
							PeerDbIds:                                pointer.From(adbsPropModel.PeerDbIds),
							PrivateEndpoint:                          pointer.From(adbsPropModel.PrivateEndpoint),
							PrivateEndpointIP:                        pointer.From(adbsPropModel.PrivateEndpointIP),
							PrivateEndpointLabel:                     pointer.From(adbsPropModel.PrivateEndpointLabel),
							ProvisionableCPUs:                        pointer.From(adbsPropModel.ProvisionableCPUs),
							ServiceConsoleUrl:                        pointer.From(adbsPropModel.ServiceConsoleUrl),
							SqlWebDeveloperUrl:                       pointer.From(adbsPropModel.SqlWebDeveloperUrl),
							SubnetId:                                 pointer.From(adbsPropModel.SubnetId),
							SupportedRegionsToCloneTo:                pointer.From(adbsPropModel.SupportedRegionsToCloneTo),
							TimeCreated:                              pointer.From(adbsPropModel.TimeCreated),
							TimeDataGuardRoleChanged:                 pointer.From(adbsPropModel.TimeDataGuardRoleChanged),
							TimeDeletionOfFreeAutonomousDatabase:     pointer.From(adbsPropModel.TimeDeletionOfFreeAutonomousDatabase),
							TimeLocalDataGuardEnabled:                pointer.From(adbsPropModel.TimeLocalDataGuardEnabled),
							TimeMaintenanceBegin:                     pointer.From(adbsPropModel.TimeMaintenanceBegin),
							TimeMaintenanceEnd:                       pointer.From(adbsPropModel.TimeMaintenanceEnd),
							TimeOfLastFailover:                       pointer.From(adbsPropModel.TimeOfLastFailover),
							TimeOfLastRefresh:                        pointer.From(adbsPropModel.TimeOfLastRefresh),
							TimeOfLastRefreshPoint:                   pointer.From(adbsPropModel.TimeOfLastRefreshPoint),
							TimeOfLastSwitchover:                     pointer.From(adbsPropModel.TimeOfLastSwitchover),
							TimeReclamationOfFreeAutonomousDatabase:  pointer.From(adbsPropModel.TimeReclamationOfFreeAutonomousDatabase),
							UsedDataStorageSizeInGbs:                 pointer.From(adbsPropModel.UsedDataStorageSizeInGbs),
							UsedDataStorageSizeInTbs:                 pointer.From(adbsPropModel.UsedDataStorageSizeInTbs),
							VnetId:                                   pointer.From(adbsPropModel.VnetId),
							WhitelistedIPs:                           pointer.From(adbsPropModel.WhitelistedIPs),
						}
					default:
						return fmt.Errorf("unexpected Autonomous Database type, must be of type Regular")
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
				output.Name = id.AutonomousDatabaseName
				output.Type = pointer.From(model.Type)
				output.Tags = utils.FlattenPtrMapStringString(model.Tags)

				metadata.SetID(id)
				return metadata.Encode(&output)
			}
			return nil
		},
	}
}
