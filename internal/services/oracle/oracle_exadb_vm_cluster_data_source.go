// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exadbvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExadbVmClusterDataSource struct{}

type ExadbVmClusterDataModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	// ExadbVMClusterProperties
	BackupSubnetCidr          string                            `tfschema:"backup_subnet_cidr"`
	ClusterName               string                            `tfschema:"cluster_name"`
	DataCollectionOptions     []ExadbDataCollectionOptionsModel `tfschema:"data_collection_options"`
	DisplayName               string                            `tfschema:"display_name"`
	Domain                    string                            `tfschema:"domain"`
	EnabledEcpuCount          int64                             `tfschema:"enabled_ecpu_count"`
	ExascaleDbStorageVaultId  string                            `tfschema:"exascale_db_storage_vault_id"`
	GiVersion                 string                            `tfschema:"gi_version"`
	GridImageOcid             string                            `tfschema:"grid_image_ocid"`
	GridImageType             string                            `tfschema:"grid_image_type"`
	Hostname                  string                            `tfschema:"hostname"`
	HostnameActual            string                            `tfschema:"hostname_actual"`
	IormConfigCache           []IormConfigModel                 `tfschema:"iorm_config_cache"`
	LicenseModel              string                            `tfschema:"license_model"`
	LifecycleDetails          string                            `tfschema:"lifecycle_details"`
	LifecycleState            string                            `tfschema:"lifecycle_state"`
	ListenerPort              int64                             `tfschema:"listener_port"`
	Location                  string                            `tfschema:"location"`
	MemorySizeInGbs           int64                             `tfschema:"memory_size_in_gbs"`
	NodeCount                 int64                             `tfschema:"node_count"`
	NsgCidrs                  []NsgCidrModel                    `tfschema:"nsg_cidrs"`
	NsgUrl                    string                            `tfschema:"nsg_url"`
	OciUrl                    string                            `tfschema:"oci_url"`
	Ocid                      string                            `tfschema:"ocid"`
	PrivateZoneOcid           string                            `tfschema:"private_zone_ocid"`
	ScanDnsName               string                            `tfschema:"scan_dns_name"`
	ScanDnsRecordId           string                            `tfschema:"scan_dns_record_id"`
	ScanIPIds                 []string                          `tfschema:"scan_ip_ids"`
	ScanListenerPortTcp       int64                             `tfschema:"scan_listener_port_tcp"`
	ScanListenerPortTcpSsl    int64                             `tfschema:"scan_listener_port_tcp_ssl"`
	Shape                     string                            `tfschema:"shape"`
	SnapshotFileSystemStorage []ExadbVmClusterStorageModel      `tfschema:"snapshot_file_system_storage"`
	SshPublicKeys             []string                          `tfschema:"ssh_public_keys"`
	SubnetId                  string                            `tfschema:"subnet_id"`
	SubnetOcid                string                            `tfschema:"subnet_ocid"`
	SystemVersion             string                            `tfschema:"system_version"`
	TimeZone                  string                            `tfschema:"time_zone"`
	TotalEcpuCount            int64                             `tfschema:"total_ecpu_count"`
	TotalFileSystemStorage    []ExadbVmClusterStorageModel      `tfschema:"total_file_system_storage"`
	VipIds                    []string                          `tfschema:"vip_ids"`
	VmFileSystemStorage       []ExadbVmClusterStorageModel      `tfschema:"vm_file_system_storage"`
	VnetId                    string                            `tfschema:"virtual_network_id"`
	ZoneOcid                  string                            `tfschema:"zone_ocid"`
}

type ExadbDataCollectionOptionsModel struct {
	IsDiagnosticsEventsEnabled bool `tfschema:"diagnostics_events_enabled"`
	IsHealthMonitoringEnabled  bool `tfschema:"health_monitoring_enabled"`
	IsIncidentLogsEnabled      bool `tfschema:"incident_logs_enabled"`
}

type IormConfigModel struct {
	DbPlans          []ExadbDbIormConfigModel `tfschema:"db_plans"`
	LifecycleDetails string                   `tfschema:"lifecycle_details"`
	LifecycleState   string                   `tfschema:"lifecycle_state"`
	Objective        string                   `tfschema:"objective"`
}

type ExadbDbIormConfigModel struct {
	DbName          string `tfschema:"db_name"`
	FlashCacheLimit string `tfschema:"flash_cache_limit"`
	Share           int64  `tfschema:"share"`
}

type ExadbVmClusterStorageModel struct {
	TotalSizeInGbs int64 `tfschema:"total_size_in_gbs"`
}

type NsgCidrModel struct {
	DestinationPortRange []PortRangeModel `tfschema:"destination_port_range"`
	Source               string           `tfschema:"source"`
}

type PortRangeModel struct {
	Max int64 `tfschema:"max"`
	Min int64 `tfschema:"min"`
}

func (d ExadbVmClusterDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
		"zones":               commonschema.ZonesMultipleOptional(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ExadbVMClusterName,
		},
	}
}

func (d ExadbVmClusterDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// ExadbVMClusterProperties
		"backup_subnet_cidr": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
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

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"domain": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"enabled_ecpu_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"exascale_db_storage_vault_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"gi_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"grid_image_ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"grid_image_type": {
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

		"nsg_cidrs": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destination_port_range": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"min": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},

					"source": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
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

		"private_zone_ocid": {
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

		"snapshot_file_system_storage": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"ssh_public_keys": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
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

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"total_ecpu_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"total_file_system_storage": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"vip_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"vm_file_system_storage": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"zone_ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d ExadbVmClusterDataSource) ModelObject() interface{} {
	return &ExadbVmClusterDataModel{}
}

func (d ExadbVmClusterDataSource) ResourceType() string {
	return "azurerm_oracle_exa_db_vm_cluster"
}

func (d ExadbVmClusterDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exadbvmclusters.ValidateExadbVMClusterID
}

func (d ExadbVmClusterDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient25.ExadbVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ExadbVmClusterDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := exadbvmclusters.NewExadbVMClusterID(subscriptionId, state.ResourceGroupName, state.Name)

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
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.BackupSubnetCidr = pointer.From(props.BackupSubnetCidr)
					state.ClusterName = pointer.From(props.ClusterName)
					state.DataCollectionOptions = FlattenExadbDataCollectionOptions(props.DataCollectionOptions)
					state.DisplayName = props.DisplayName
					state.Domain = pointer.From(props.Domain)
					state.EnabledEcpuCount = props.EnabledEcpuCount
					state.ExascaleDbStorageVaultId = props.ExascaleDbStorageVaultId
					state.GiVersion = pointer.From(props.GiVersion)
					state.GridImageOcid = pointer.From(props.GridImageOcid)
					state.GridImageType = string(pointer.From(props.GridImageType))
					state.Hostname = removeHostnameSuffix(props.Hostname)
					state.HostnameActual = props.Hostname
					state.IormConfigCache = FlattenIormConfig(props.IormConfigCache)
					state.LicenseModel = string(pointer.From(props.LicenseModel))
					state.LifecycleDetails = pointer.From(props.LifecycleDetails)
					state.LifecycleState = string(*props.LifecycleState)
					state.ListenerPort = pointer.From(props.ListenerPort)
					state.MemorySizeInGbs = pointer.From(props.MemorySizeInGbs)
					state.NodeCount = props.NodeCount
					state.NsgCidrs = FlattenNsgCidrs(props.NsgCidrs)
					state.NsgUrl = pointer.From(props.NsgURL)
					state.OciUrl = pointer.From(props.OciURL)
					state.Ocid = pointer.From(props.Ocid)
					state.PrivateZoneOcid = pointer.From(props.PrivateZoneOcid)
					state.ScanDnsName = pointer.From(props.ScanDnsName)
					state.ScanDnsRecordId = pointer.From(props.ScanDnsRecordId)
					state.ScanIPIds = pointer.From(props.ScanIPIds)
					state.ScanListenerPortTcp = pointer.From(props.ScanListenerPortTcp)
					state.ScanListenerPortTcpSsl = pointer.From(props.ScanListenerPortTcpSsl)
					state.Shape = props.Shape
					state.SnapshotFileSystemStorage = FlattenExadbVmClusterStorage(props.SnapshotFileSystemStorage)
					state.SshPublicKeys = props.SshPublicKeys
					state.SubnetId = props.SubnetId
					state.SubnetOcid = pointer.From(props.SubnetOcid)
					state.SystemVersion = pointer.From(props.SystemVersion)
					state.TimeZone = pointer.From(props.TimeZone)
					state.TotalEcpuCount = props.TotalEcpuCount
					state.TotalFileSystemStorage = FlattenExadbVmClusterStorage(props.TotalFileSystemStorage)
					state.VipIds = *props.VipIds
					state.VmFileSystemStorage = FlattenVMFileSystemStorage(props.VMFileSystemStorage)
					state.VnetId = props.VnetId
					state.ZoneOcid = pointer.From(props.ZoneOcid)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func FlattenIormConfig(input *exadbvmclusters.ExadataIormConfig) []IormConfigModel {
	output := make([]IormConfigModel, 0)

	if input != nil {
		var dbIormConfigModel []ExadbDbIormConfigModel
		if input.DbPlans != nil {
			dbPlans := *input.DbPlans
			for _, dbPlan := range dbPlans {
				dbIormConfigModel = append(dbIormConfigModel, ExadbDbIormConfigModel{
					DbName:          pointer.From(dbPlan.DbName),
					FlashCacheLimit: pointer.From(dbPlan.FlashCacheLimit),
					Share:           pointer.From(dbPlan.Share),
				})
			}
		}
		return append(output, IormConfigModel{
			DbPlans:          dbIormConfigModel,
			LifecycleDetails: pointer.From(input.LifecycleDetails),
			LifecycleState:   string(pointer.From(input.LifecycleState)),
			Objective:        string(pointer.From(input.Objective)),
		})
	}

	return output
}

func FlattenExadbVmClusterStorage(input *exadbvmclusters.ExadbVMClusterStorageDetails) []ExadbVmClusterStorageModel {
	output := make([]ExadbVmClusterStorageModel, 0)
	if input != nil {
		return append(output, ExadbVmClusterStorageModel{
			TotalSizeInGbs: input.TotalSizeInGbs,
		})
	}
	return output
}

func FlattenVMFileSystemStorage(input exadbvmclusters.ExadbVMClusterStorageDetails) []ExadbVmClusterStorageModel {
	output := make([]ExadbVmClusterStorageModel, 0)
	return append(output, ExadbVmClusterStorageModel{
		TotalSizeInGbs: input.TotalSizeInGbs,
	})
}

func FlattenNsgCidrs(input *[]exadbvmclusters.NsgCidr) []NsgCidrModel {
	output := make([]NsgCidrModel, 0)

	if input != nil {
		for _, nsgCidr := range *input {
			var portRangeModel []PortRangeModel
			if nsgCidr.DestinationPortRange != nil {
				portRangeModel = append(portRangeModel, PortRangeModel{
					Max: nsgCidr.DestinationPortRange.Max,
					Min: nsgCidr.DestinationPortRange.Min,
				})
			}
			output = append(output, NsgCidrModel{
				DestinationPortRange: portRangeModel,
				Source:               nsgCidr.Source,
			})
		}
	}

	return output
}
