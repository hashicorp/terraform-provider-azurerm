// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CloudVmClusterDataSource struct{}

type CloudVmClusterDataModel struct {
	Name string                 `tfschema:"name"`
	Type string                 `tfschema:"type"`
	Tags map[string]interface{} `tfschema:"tags"`

	// SystemData
	SystemData []SystemDataModel `tfschema:"system_data"`

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
	IormConfigCache              []ExadataIormConfigModel     `tfschema:"iorm_config_cache"`
	IsLocalBackupEnabled         bool                         `tfschema:"is_local_backup_enabled"`
	IsSparseDiskgroupEnabled     bool                         `tfschema:"is_sparse_diskgroup_enabled"`
	LastUpdateHistoryEntryId     string                       `tfschema:"last_update_history_entry_id"`
	LicenseModel                 string                       `tfschema:"license_model"`
	LifecycleDetails             string                       `tfschema:"lifecycle_details"`
	LifecycleState               string                       `tfschema:"lifecycle_state"`
	ListenerPort                 int64                        `tfschema:"listener_port"`
	MemorySizeInGbs              int64                        `tfschema:"memory_size_in_gbs"`
	NodeCount                    int64                        `tfschema:"node_count"`
	NsgUrl                       string                       `tfschema:"nsg_url"`
	OciUrl                       string                       `tfschema:"oci_url"`
	Ocid                         string                       `tfschema:"ocid"`
	OcpuCount                    float64                      `tfschema:"ocpu_count"`
	ProvisioningState            string                       `tfschema:"provisioning_state"`
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
	VnetId                       string                       `tfschema:"vnet_id"`
	ZoneId                       string                       `tfschema:"zone_id"`
}

type SystemDataModel struct {
	CreatedBy          string `tfschema:"created_by"`
	CreatedByType      string `tfschema:"created_by_type"`
	CreatedAt          string `tfschema:"created_at"`
	LastModifiedBy     string `tfschema:"last_modified_by"`
	LastModifiedbyType string `tfschema:"last_modified_by_type"`
	LastModifiedAt     string `tfschema:"last_modified_at"`
}

type DataCollectionOptionsModel struct {
	IsDiagnosticsEventsEnabled bool `tfschema:"is_diagnostics_events_enabled"`
	IsHealthMonitoringEnabled  bool `tfschema:"is_health_monitoring_enabled"`
	IsIncidentLogsEnabled      bool `tfschema:"is_incident_logs_enabled"`
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
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d CloudVmClusterDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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
					"is_diagnostics_events_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"is_health_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"is_incident_logs_enabled": {
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

		"is_local_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_sparse_diskgroup_enabled": {
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

		"provisioning_state": {
			Type:     pluginsdk.TypeString,
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

		"vnet_id": {
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
	return nil
}

func (d CloudVmClusterDataSource) ResourceType() string {
	return "azurerm_oracledatabase_cloud_vm_cluster"
}

func (d CloudVmClusterDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudvmclusters.ValidateCloudVMClusterID
}

func (d CloudVmClusterDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := cloudvmclusters.NewCloudVMClusterID(subscriptionId,
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

				var output CloudVmClusterDataModel
				prop := model.Properties
				if prop != nil {
					output = CloudVmClusterDataModel{
						BackupSubnetCidr:             pointer.From(prop.BackupSubnetCidr),
						CloudExadataInfrastructureId: prop.CloudExadataInfrastructureId,
						ClusterName:                  pointer.From(prop.ClusterName),
						CompartmentId:                pointer.From(prop.CompartmentId),
						ComputeNodes:                 pointer.From(prop.ComputeNodes),
						CpuCoreCount:                 prop.CpuCoreCount,
						DataStoragePercentage:        pointer.From(prop.DataStoragePercentage),
						DataStorageSizeInTbs:         pointer.From(prop.DataStorageSizeInTbs),
						DbNodeStorageSizeInGbs:       pointer.From(prop.DbNodeStorageSizeInGbs),
						DbServers:                    pointer.From(prop.DbServers),
						DiskRedundancy:               string(pointer.From(prop.DiskRedundancy)),
						DisplayName:                  prop.DisplayName,
						Domain:                       pointer.From(prop.Domain),
						GiVersion:                    prop.GiVersion,
						Hostname:                     prop.Hostname,
						IormConfigCache:              ConvertExadataIormConfigToInternal(prop.IormConfigCache),
						IsLocalBackupEnabled:         pointer.From(prop.IsLocalBackupEnabled),
						IsSparseDiskgroupEnabled:     pointer.From(prop.IsSparseDiskgroupEnabled),
						LastUpdateHistoryEntryId:     pointer.From(prop.LastUpdateHistoryEntryId),
						LicenseModel:                 string(pointer.From(prop.LicenseModel)),
						LifecycleDetails:             pointer.From(prop.LifecycleDetails),
						LifecycleState:               string(*prop.LifecycleState),
						ListenerPort:                 pointer.From(prop.ListenerPort),
						MemorySizeInGbs:              pointer.From(prop.MemorySizeInGbs),
						NodeCount:                    pointer.From(prop.NodeCount),
						NsgUrl:                       pointer.From(prop.NsgUrl),
						OciUrl:                       pointer.From(prop.OciUrl),
						Ocid:                         pointer.From(prop.Ocid),
						ProvisioningState:            string(pointer.From(prop.ProvisioningState)),
						Shape:                        pointer.From(prop.Shape),
						StorageSizeInGbs:             pointer.From(prop.StorageSizeInGbs),
						SubnetId:                     prop.SubnetId,
						SubnetOcid:                   pointer.From(prop.SubnetOcid),
						SystemVersion:                pointer.From(prop.SystemVersion),
						TimeCreated:                  pointer.From(prop.TimeCreated),
						TimeZone:                     pointer.From(prop.TimeZone),
						ZoneId:                       pointer.From(prop.ZoneId),
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
				output.Name = id.CloudVmClusterName
				output.Type = pointer.From(model.Type)
				output.Tags = utils.FlattenPtrMapStringString(model.Tags)

				metadata.SetID(id)
				return metadata.Encode(&output)
			}
			return nil
		},
	}
}
