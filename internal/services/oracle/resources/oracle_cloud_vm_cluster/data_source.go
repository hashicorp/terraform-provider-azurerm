// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_cloud_vm_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type CloudVmClusterDataModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// CloudVMClusterProperties
	BackupSubnetCidr             string                         `tfschema:"backup_subnet_cidr"`
	CloudExadataInfrastructureId string                         `tfschema:"cloud_exadata_infrastructure_id"`
	ClusterName                  string                         `tfschema:"cluster_name"`
	CompartmentId                string                         `tfschema:"compartment_id"`
	ComputeModel                 string                         `tfschema:"compute_model"`
	ComputeNodes                 []string                       `tfschema:"compute_nodes"`
	CpuCoreCount                 int64                          `tfschema:"cpu_core_count"`
	DataCollectionOptions        []DataCollectionOptionsModel   `tfschema:"data_collection_options"`
	DataStoragePercentage        int64                          `tfschema:"data_storage_percentage"`
	DataStorageSizeInTbs         float64                        `tfschema:"data_storage_size_in_tbs"`
	DbNodeStorageSizeInGbs       int64                          `tfschema:"db_node_storage_size_in_gbs"`
	DbServers                    []string                       `tfschema:"db_servers"`
	DiskRedundancy               string                         `tfschema:"disk_redundancy"`
	DisplayName                  string                         `tfschema:"display_name"`
	Domain                       string                         `tfschema:"domain"`
	FileSystemConfiguration      []FileSystemConfigurationModel `tfschema:"file_system_configuration"`
	GiVersion                    string                         `tfschema:"gi_version"`
	Hostname                     string                         `tfschema:"hostname"`
	HostnameActual               string                         `tfschema:"hostname_actual"`
	IormConfigCache              []ExadataIormConfigModel       `tfschema:"iorm_config_cache"`
	IsLocalBackupEnabled         bool                           `tfschema:"local_backup_enabled"`
	IsSparseDiskgroupEnabled     bool                           `tfschema:"sparse_diskgroup_enabled"`
	LastUpdateHistoryEntryId     string                         `tfschema:"last_update_history_entry_id"`
	LicenseModel                 string                         `tfschema:"license_model"`
	LifecycleDetails             string                         `tfschema:"lifecycle_details"`
	LifecycleState               string                         `tfschema:"lifecycle_state"`
	ListenerPort                 int64                          `tfschema:"listener_port"`
	Location                     string                         `tfschema:"location"`
	MemorySizeInGbs              int64                          `tfschema:"memory_size_in_gbs"`
	NodeCount                    int64                          `tfschema:"node_count"`
	NsgUrl                       string                         `tfschema:"nsg_url"`
	OciUrl                       string                         `tfschema:"oci_url"`
	Ocid                         string                         `tfschema:"ocid"`
	OcpuCount                    float64                        `tfschema:"ocpu_count"`
	ScanDnsName                  string                         `tfschema:"scan_dns_name"`
	ScanDnsRecordId              string                         `tfschema:"scan_dns_record_id"`
	ScanIPIds                    []string                       `tfschema:"scan_ip_ids"`
	ScanListenerPortTcp          int64                          `tfschema:"scan_listener_port_tcp"`
	ScanListenerPortTcpSsl       int64                          `tfschema:"scan_listener_port_tcp_ssl"`
	Shape                        string                         `tfschema:"shape"`
	SshPublicKeys                []string                       `tfschema:"ssh_public_keys"`
	StorageSizeInGbs             int64                          `tfschema:"storage_size_in_gbs"`
	SubnetId                     string                         `tfschema:"subnet_id"`
	SubnetOcid                   string                         `tfschema:"subnet_ocid"`
	SystemVersion                string                         `tfschema:"system_version"`
	TimeCreated                  string                         `tfschema:"time_created"`
	TimeZone                     string                         `tfschema:"time_zone"`
	VipIds                       []string                       `tfschema:"vip_ods"`
	VnetId                       string                         `tfschema:"virtual_network_id"`
	ZoneId                       string                         `tfschema:"zone_id"`
}

type FileSystemConfigurationModel struct {
	MountPoint string `tfschema:"mount_point"`
	SizeInGb   int64  `tfschema:"size_in_gb"`
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

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Description:  "The name of this Cloud VM Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.CloudVMClusterName,
		},
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// CloudVMClusterProperties
		"backup_subnet_cidr": {
			Description: "Client OCI backup subnet CIDR, default is `192.168.252.0/22`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"cloud_exadata_infrastructure_id": {
			Description: "The Cloud Exadata Infrastructure ID.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"cluster_name": {
			Description: "The cluster name for Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"compartment_id": {
			Description: "The OCID of the compartment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"compute_model": {
			Description: "The compute model of the Exadata Infrastructure.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"compute_nodes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cpu_core_count": {
			Description: "The number of CPU cores enabled on the Cloud VM Cluster.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"data_collection_options": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"diagnostics_events_enabled": {
						Description: "Whether diagnostic collection is enabled for the VM Cluster/Cloud VM Cluster/VMBM DBCS. Enabling diagnostic collection allows you to receive Events service notifications for guest VM issues. Diagnostic collection also allows Oracle to provide enhanced service and proactive support for your Exadata system. You can enable diagnostic collection during VM Cluster/Cloud VM Cluster provisioning. You can also disable or enable it at any time using the `UpdateVmCluster` or `updateCloudVmCluster` API.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"health_monitoring_enabled": {
						Description: "Whether health monitoring is enabled for the VM Cluster / Cloud VM Cluster / VMBM DBCS. Enabling health monitoring allows Oracle to collect diagnostic data and share it with its operations and support personnel. You may also receive notifications for some events. Collecting health diagnostics enables Oracle to provide proactive support and enhanced service for your system. Optionally enable health monitoring while provisioning a system. You can also disable or enable health monitoring anytime using the `UpdateVmCluster`, `UpdateCloudVmCluster` or `updateDbsystem` API.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"incident_logs_enabled": {
						Description: "Whether incident logs and trace collection are enabled for the VM Cluster / Cloud VM Cluster / VMBM DBCS. Enabling incident logs collection allows Oracle to receive Events service notifications for guest VM issues, collect incident logs and traces, and use them to diagnose issues and resolve them. Optionally enable incident logs collection while provisioning a system. You can also disable or enable incident logs collection anytime using the `UpdateVmCluster`, `updateCloudVmCluster` or `updateDbsystem` API.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
				},
			},
		},

		"data_storage_percentage": {
			Description: "The percentage assigned to DATA storage (user data and database files). The remaining percentage is assigned to RECO storage (database redo logs, archive logs, and recovery manager backups). See [Storage Configuration](https://docs.oracle.com/en-us/iaas/exadatacloud/index.html#Exadata) in the Exadata documentation for details on the impact of the configuration settings on storage.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"data_storage_size_in_tbs": {
			Description: "The data disk group size to be allocated in TBs.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"db_node_storage_size_in_gbs": {
			Description: "The local node storage to be allocated in GBs.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"db_servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"disk_redundancy": {
			Description: "The type of redundancy configured for the Cloud VM Cluster. `NORMAL` is 2-way redundancy. `HIGH` is 3-way redundancy.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The user-friendly name for the Cloud VM Cluster. The name does not need to be unique.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"domain": {
			Description: "The domain name for the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"file_system_configuration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mount_point": {
						Description: "The mount path of the file system.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"size_in_gb": {
						Description: "The size of the Virtual Machine's file system.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
				},
			},
		},

		"gi_version": {
			Description: "A valid Oracle Grid Infrastructure (GI) software version.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"hostname": {
			Description: "The hostname for the Cloud VM Cluster without suffix.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"hostname_actual": {
			Description: "The hostname for the Cloud VM Cluster with suffix.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
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
									Description: "The database name. For the default `DbPlan`, the `dbName` is `default`.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"flash_cache_limit": {
									Description: "The flash cache limit for this database. This value is internally configured based on the share value assigned to the database.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"share": {
									Description: "The relative priority of this database.",
									Type:        pluginsdk.TypeInt,
									Computed:    true,
								},
							},
						},
					},

					"lifecycle_details": {
						Description: "Additional information about the current `lifecycleState`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"lifecycle_state": {
						Description: "The current state of IORM configuration for the Exadata DB system.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"objective": {
						Description: "The current value for the IORM objective.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"local_backup_enabled": {
			Description: "If true, database backup on local Exadata storage is configured for the Cloud VM Cluster. If false, database backup on local Exadata storage is not available in the Cloud VM Cluster.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"sparse_diskgroup_enabled": {
			Description: "If true, sparse disk group is configured for the Cloud VM Cluster. If false, sparse disk group is not created.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"last_update_history_entry_id": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the last maintenance update history entry. This value is updated when a maintenance update starts.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"license_model": {
			Description: "The Oracle license model that applies to the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"lifecycle_details": {
			Description: "Additional information about the current lifecycle state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"lifecycle_state": {
			Description: "Cloud VM Cluster lifecycle state enum.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"listener_port": {
			Description: "The port number configured for the listener on the Cloud VM Cluster.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"memory_size_in_gbs": {
			Description: "The memory to be allocated in GBs.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"node_count": {
			Description: "The number of nodes in the Cloud VM Cluster.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"nsg_url": {
			Description: "The list of [OCIDs](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) for the Network Security Groups (NSGs) to which this resource belongs. For more information about NSGs, see [Security Rules](https://docs.oracle.com/en-us/iaas/Content/Network/Concepts/securityrules.htm).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"oci_url": {
			Description: "The URL of the resource in the OCI console.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ocid": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ocpu_count": {
			Description: "The number of OCPU cores to enable on the Cloud VM Cluster. Only 1 decimal place is allowed for the fractional part.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},

		"scan_dns_name": {
			Description: "The FQDN of the DNS record for the SCAN IP addresses that are associated with the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"scan_dns_record_id": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DNS record for the SCAN IP addresses that are associated with the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"scan_ip_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"scan_listener_port_tcp": {
			Description: "The TCP Single Client Access Name (SCAN) port.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"scan_listener_port_tcp_ssl": {
			Description: "The TCPS Single Client Access Name (SCAN) port.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"shape": {
			Description: "The model name of the Exadata hardware running the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ssh_public_keys": {
			Description: "The public key portion of one or more key pairs used for SSH access to the Cloud VM Cluster.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"storage_size_in_gbs": {
			Description: "The storage allocation for the disk group, in gigabytes (GB).",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"subnet_id": {
			Description: "The ID of the Azure Resource Manager subnet resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subnet_ocid": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the subnet associated with the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"system_version": {
			Description: "Operating system version of the Exadata image.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_created": {
			Description: "The date and time that the Cloud VM Cluster was created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_zone": {
			Description: "The time zone of the Cloud VM Cluster. For details, see [Exadata Infrastructure Time Zones](https://docs.oracle.com/en-us/iaas/base-database/doc/manage-time-zone.html).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"vip_ods": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the virtual IP (VIP) addresses associated with the Cloud VM Cluster. The Cluster Ready Services (CRS) creates and maintains one VIP address for each node in the Exadata Cloud Service instance to enable failover. If one node fails, the VIP is reassigned to another active node in the Cluster.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_network_id": {
			Description: "The ID to an Azure Resource Manager Virtual Network resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"zone_id": {
			Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the zone the Cloud VM Cluster is associated with.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d DataSource) ModelObject() any {
	return &CloudVmClusterDataModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_oracle_cloud_vm_cluster"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Cloud VM Cluster."
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudvmclusters.ValidateCloudVMClusterID
}

func (d DataSource) Read() sdk.ResourceFunc {
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
					state.ComputeModel = pointer.FromEnum(props.ComputeModel)
					state.ComputeNodes = pointer.From(props.ComputeNodes)
					state.CpuCoreCount = props.CpuCoreCount
					state.DataStoragePercentage = pointer.From(props.DataStoragePercentage)
					state.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
					state.DbNodeStorageSizeInGbs = pointer.From(props.DbNodeStorageSizeInGbs)
					state.DbServers = pointer.From(props.DbServers)
					state.DiskRedundancy = pointer.FromEnum(props.DiskRedundancy)
					state.DisplayName = props.DisplayName
					state.Domain = pointer.From(props.Domain)
					state.FileSystemConfiguration = FlattenFileSystemConfigurationDetails(props.FileSystemConfigurationDetails)
					state.GiVersion = props.GiVersion
					state.Hostname = removeHostnameSuffix(props.Hostname)
					state.HostnameActual = props.Hostname
					state.IormConfigCache = FlattenExadataIormConfig(props.IormConfigCache)
					state.IsLocalBackupEnabled = pointer.From(props.IsLocalBackupEnabled)
					state.IsSparseDiskgroupEnabled = pointer.From(props.IsSparseDiskgroupEnabled)
					state.LastUpdateHistoryEntryId = pointer.From(props.LastUpdateHistoryEntryId)
					state.LicenseModel = pointer.FromEnum(props.LicenseModel)
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
			LifecycleState:   pointer.FromEnum(input.LifecycleState),
			Objective:        pointer.FromEnum(input.Objective),
		})
	}

	return output
}
