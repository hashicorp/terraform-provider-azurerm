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
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CloudVmClusterDataSource struct{}

type CloudVmClusterDataModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// CloudVMClusterProperties
	BackupSubnetCidr             string                       `tfschema:"backup_subnet_cidr"`
	CloudExadataInfrastructureId string                       `tfschema:"cloud_exadata_infrastructure_id"`
	ClusterName                  string                       `tfschema:"cluster_name"`
	CompartmentId                string                       `tfschema:"compartment_id"`
	ComputeNodes                 []string                     `tfschema:"compute_nodes"`
	CpuCoreCount                 int64                        `tfschema:"cpu_core_count"`
	DataCollectionOptions        []DataCollectionOptionsModel `tfschema:"data_collection_options"`
	DataStoragePercentage        int64                        `tfschema:"data_storage_percentage"`
	DataStorageSizeInTbs         float64                      `tfschema:"data_storage_size_in_tbs"`
	DbNodeStorageSizeInGbs       int64                        `tfschema:"db_node_storage_size_in_gbs"`
	DbServers                    []string                     `tfschema:"db_servers"`
	DiskRedundancy               string                       `tfschema:"disk_redundancy"`
	DisplayName                  string                       `tfschema:"display_name"`
	Domain                       string                       `tfschema:"domain"`
	GiVersion                    string                       `tfschema:"gi_version"`
	Hostname                     string                       `tfschema:"hostname"`
	HostnameActual               string                       `tfschema:"hostname_actual"`
	IormConfigCache              []ExadataIormConfigModel     `tfschema:"iorm_config_cache"`
	IsLocalBackupEnabled         bool                         `tfschema:"local_backup_enabled"`
	IsSparseDiskgroupEnabled     bool                         `tfschema:"sparse_diskgroup_enabled"`
	LastUpdateHistoryEntryId     string                       `tfschema:"last_update_history_entry_id"`
	LicenseModel                 string                       `tfschema:"license_model"`
	LifecycleDetails             string                       `tfschema:"lifecycle_details"`
	LifecycleState               string                       `tfschema:"lifecycle_state"`
	ListenerPort                 int64                        `tfschema:"listener_port"`
	Location                     string                       `tfschema:"location"`
	MemorySizeInGbs              int64                        `tfschema:"memory_size_in_gbs"`
	NodeCount                    int64                        `tfschema:"node_count"`
	NsgUrl                       string                       `tfschema:"nsg_url"`
	OciUrl                       string                       `tfschema:"oci_url"`
	Ocid                         string                       `tfschema:"ocid"`
	OcpuCount                    float64                      `tfschema:"ocpu_count"`
	ScanDnsName                  string                       `tfschema:"scan_dns_name"`
	ScanDnsRecordId              string                       `tfschema:"scan_dns_record_id"`
	ScanIPIds                    []string                     `tfschema:"scan_ip_ids"`
	ScanListenerPortTcp          int64                        `tfschema:"scan_listener_port_tcp"`
	ScanListenerPortTcpSsl       int64                        `tfschema:"scan_listener_port_tcp_ssl"`
	Shape                        string                       `tfschema:"shape"`
	SshPublicKeys                []string                     `tfschema:"ssh_public_keys"`
	StorageSizeInGbs             int64                        `tfschema:"storage_size_in_gbs"`
	SubnetId                     string                       `tfschema:"subnet_id"`
	SubnetOcid                   string                       `tfschema:"subnet_ocid"`
	SystemVersion                string                       `tfschema:"system_version"`
	TimeCreated                  string                       `tfschema:"time_created"`
	TimeZone                     string                       `tfschema:"time_zone"`
	VipIds                       []string                     `tfschema:"vip_ods"`
	VnetId                       string                       `tfschema:"virtual_network_id"`
	ZoneId                       string                       `tfschema:"zone_id"`
}

type DataCollectionOptionsModel struct {
	IsDiagnosticsEventsEnabled bool `tfschema:"diagnostics_events_enabled"`
	IsHealthMonitoringEnabled  bool `tfschema:"health_monitoring_enabled"`
	IsIncidentLogsEnabled      bool `tfschema:"incident_logs_enabled"`
}

type ExadataIormConfigModel struct {
	DbPlans          []DbIormConfigModel `tfschema:"db_plans"`
	LifecycleDetails string              `tfschema:"lifecycle_details"`
	LifecycleState   string              `tfschema:"lifecycle_state"`
	Objective        string              `tfschema:"objective"`
}

type DbIormConfigModel struct {
	DbName          string `tfschema:"db_name"`
	FlashCacheLimit string `tfschema:"flash_cache_limit"`
	Share           int64  `tfschema:"share"`
}

func (d CloudVmClusterDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.CloudVMClusterName,
		},
	}
}

func (d CloudVmClusterDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// CloudVMClusterProperties
		"backup_subnet_cidr": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cloud_exadata_infrastructure_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"compartment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"compute_nodes": {
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

		"data_collection_options": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"diagnostics_events_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"health_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"incident_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"data_storage_percentage": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"db_node_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"db_servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"disk_redundancy": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"domain": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"gi_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname_actual": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"iorm_config_cache": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"db_plans": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"db_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"flash_cache_limit": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"share": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},

					"lifecycle_details": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"lifecycle_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"objective": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"local_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"sparse_diskgroup_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"last_update_history_entry_id": {
			Type:     pluginsdk.TypeString,
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

		"listener_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"memory_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"nsg_url": {
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

		"ocpu_count": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"scan_dns_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scan_dns_record_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scan_ip_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"scan_listener_port_tcp": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"scan_listener_port_tcp_ssl": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"shape": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ssh_public_keys": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet_ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"system_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_created": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vip_ods": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"zone_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d CloudVmClusterDataSource) ModelObject() interface{} {
	return &CloudVmClusterDataModel{}
}

func (d CloudVmClusterDataSource) ResourceType() string {
	return "azurerm_oracle_cloud_vm_cluster"
}

func (d CloudVmClusterDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudvmclusters.ValidateCloudVMClusterID
}

func (d CloudVmClusterDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state CloudVmClusterDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := cloudvmclusters.NewCloudVMClusterID(subscriptionId, state.ResourceGroupName, state.Name)

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
					state.BackupSubnetCidr = pointer.From(props.BackupSubnetCidr)
					state.CloudExadataInfrastructureId = props.CloudExadataInfrastructureId
					state.ClusterName = pointer.From(props.ClusterName)
					state.CompartmentId = pointer.From(props.CompartmentId)
					state.ComputeNodes = pointer.From(props.ComputeNodes)
					state.CpuCoreCount = props.CpuCoreCount
					state.DataStoragePercentage = pointer.From(props.DataStoragePercentage)
					state.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
					state.DbNodeStorageSizeInGbs = pointer.From(props.DbNodeStorageSizeInGbs)
					state.DbServers = pointer.From(props.DbServers)
					state.DiskRedundancy = string(pointer.From(props.DiskRedundancy))
					state.DisplayName = props.DisplayName
					state.Domain = pointer.From(props.Domain)
					state.GiVersion = props.GiVersion
					state.Hostname = removeHostnameSuffix(props.Hostname)
					state.HostnameActual = props.Hostname
					state.IormConfigCache = FlattenExadataIormConfig(props.IormConfigCache)
					state.IsLocalBackupEnabled = pointer.From(props.IsLocalBackupEnabled)
					state.IsSparseDiskgroupEnabled = pointer.From(props.IsSparseDiskgroupEnabled)
					state.LastUpdateHistoryEntryId = pointer.From(props.LastUpdateHistoryEntryId)
					state.LicenseModel = string(pointer.From(props.LicenseModel))
					state.LifecycleDetails = pointer.From(props.LifecycleDetails)
					state.LifecycleState = string(*props.LifecycleState)
					state.ListenerPort = pointer.From(props.ListenerPort)
					state.MemorySizeInGbs = pointer.From(props.MemorySizeInGbs)
					state.NodeCount = pointer.From(props.NodeCount)
					state.NsgUrl = pointer.From(props.NsgURL)
					state.OciUrl = pointer.From(props.OciURL)
					state.Ocid = pointer.From(props.Ocid)
					state.Shape = pointer.From(props.Shape)
					state.StorageSizeInGbs = pointer.From(props.StorageSizeInGbs)
					state.SubnetId = props.SubnetId
					state.SubnetOcid = pointer.From(props.SubnetOcid)
					state.SystemVersion = pointer.From(props.SystemVersion)
					state.TimeCreated = pointer.From(props.TimeCreated)
					state.TimeZone = pointer.From(props.TimeZone)
					state.VnetId = props.VnetId
					state.ZoneId = pointer.From(props.ZoneId)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func FlattenExadataIormConfig(input *cloudvmclusters.ExadataIormConfig) []ExadataIormConfigModel {
	output := make([]ExadataIormConfigModel, 0)

	if input != nil {
		var dbIormConfigModel []DbIormConfigModel
		if input.DbPlans != nil {
			dbPlans := *input.DbPlans
			for _, dbPlan := range dbPlans {
				dbIormConfigModel = append(dbIormConfigModel, DbIormConfigModel{
					DbName:          pointer.From(dbPlan.DbName),
					FlashCacheLimit: pointer.From(dbPlan.FlashCacheLimit),
					Share:           pointer.From(dbPlan.Share),
				})
			}
		}
		return append(output, ExadataIormConfigModel{
			DbPlans:          dbIormConfigModel,
			LifecycleDetails: pointer.From(input.LifecycleDetails),
			LifecycleState:   string(pointer.From(input.LifecycleState)),
			Objective:        string(pointer.From(input.Objective)),
		})
	}

	return output
}
