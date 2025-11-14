// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exadbvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ExascaleDatabaseVirtualMachineClusterDataSource struct{}

type ExascaleDatabaseVirtualMachineClusterDataModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	BackupSubnetCidr                         string                                              `tfschema:"backup_subnet_cidr"`
	ClusterName                              string                                              `tfschema:"cluster_name"`
	DataCollection                           []ExascaleDatabaseDataCollectionModel               `tfschema:"data_collection"`
	DisplayName                              string                                              `tfschema:"display_name"`
	Domain                                   string                                              `tfschema:"domain"`
	EnabledEcpuCount                         int64                                               `tfschema:"enabled_ecpu_count"`
	ExascaleDatabaseStorageVaultId           string                                              `tfschema:"exascale_database_storage_vault_id"`
	GridInfrastructureVersion                string                                              `tfschema:"grid_infrastructure_version"`
	GridImageOcid                            string                                              `tfschema:"grid_image_ocid"`
	GridImageType                            string                                              `tfschema:"grid_image_type"`
	Hostname                                 string                                              `tfschema:"hostname"`
	HostnameActual                           string                                              `tfschema:"hostname_actual"`
	IormConfigCache                          []IormConfigModel                                   `tfschema:"iorm_config_cache"`
	LicenseModel                             string                                              `tfschema:"license_model"`
	LifecycleDetails                         string                                              `tfschema:"lifecycle_details"`
	LifecycleState                           string                                              `tfschema:"lifecycle_state"`
	ListenerPort                             int64                                               `tfschema:"listener_port"`
	Location                                 string                                              `tfschema:"location"`
	MemorySizeInGb                           int64                                               `tfschema:"memory_size_in_gb"`
	NodeCount                                int64                                               `tfschema:"node_count"`
	NetworkSecurityGroupCidr                 []NetworkSecurityGroupCidrModel                     `tfschema:"network_security_group_cidr"`
	NetworkSecurityGroupUrl                  string                                              `tfschema:"network_security_group_url"`
	OciUrl                                   string                                              `tfschema:"oci_url"`
	Ocid                                     string                                              `tfschema:"ocid"`
	PrivateZoneOcid                          string                                              `tfschema:"private_zone_ocid"`
	SingleClientAccessNameDnsName            string                                              `tfschema:"single_client_access_name_dns_name"`
	SingleClientAccessNameDnsRecordId        string                                              `tfschema:"single_client_access_name_dns_record_id"`
	SingleClientAccessNameIpIds              []string                                            `tfschema:"single_client_access_name_ip_ids"`
	SingleClientAccessNameListenerPortTcp    int64                                               `tfschema:"single_client_access_name_listener_port_tcp"`
	SingleClientAccessNameListenerPortTcpSsl int64                                               `tfschema:"single_client_access_name_listener_port_tcp_ssl"`
	Shape                                    string                                              `tfschema:"shape"`
	SnapshotFileSystemStorage                []ExascaleDatabaseVirtualMachineClusterStorageModel `tfschema:"snapshot_file_system_storage"`
	SshPublicKeys                            []string                                            `tfschema:"ssh_public_keys"`
	SubnetId                                 string                                              `tfschema:"subnet_id"`
	SubnetOcid                               string                                              `tfschema:"subnet_ocid"`
	SystemVersion                            string                                              `tfschema:"system_version"`
	TimeZone                                 string                                              `tfschema:"time_zone"`
	TotalEcpuCount                           int64                                               `tfschema:"total_ecpu_count"`
	TotalFileSystemStorage                   []ExascaleDatabaseVirtualMachineClusterStorageModel `tfschema:"total_file_system_storage"`
	VirtualIpIds                             []string                                            `tfschema:"virtual_ip_ids"`
	VirtualMachineFileSystemStorage          []ExascaleDatabaseVirtualMachineClusterStorageModel `tfschema:"virtual_machine_file_system_storage"`
	VnetId                                   string                                              `tfschema:"virtual_network_id"`
	ZoneOcid                                 string                                              `tfschema:"zone_ocid"`
}

type IormConfigModel struct {
	DatabasePlans    []ExascaleDatabaseIormConfigModel `tfschema:"database_plans"`
	LifecycleDetails string                            `tfschema:"lifecycle_details"`
	LifecycleState   string                            `tfschema:"lifecycle_state"`
	Objective        string                            `tfschema:"objective"`
}

type ExascaleDatabaseIormConfigModel struct {
	DatabaseName    string `tfschema:"database_name"`
	FlashCacheLimit string `tfschema:"flash_cache_limit"`
	Share           int64  `tfschema:"share"`
}

type ExascaleDatabaseVirtualMachineClusterStorageModel struct {
	TotalSizeInGb int64 `tfschema:"total_size_in_gb"`
}

func (d ExascaleDatabaseVirtualMachineClusterDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
		"zones":               commonschema.ZonesMultipleOptional(),

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_]`), "Name must start with a letter or underscore (_)"),
				validation.StringDoesNotContainAny("--"),
			),
		},
	}
}

func (d ExascaleDatabaseVirtualMachineClusterDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"backup_subnet_cidr": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"data_collection": {
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

		"exascale_database_storage_vault_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"grid_infrastructure_version": {
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
					"database_plans": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"database_name": {
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

		"memory_size_in_gb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"network_security_group_cidr": {
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

		"network_security_group_url": {
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

		"single_client_access_name_dns_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"single_client_access_name_dns_record_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"single_client_access_name_ip_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"single_client_access_name_listener_port_tcp": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"single_client_access_name_listener_port_tcp_ssl": {
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
					"total_size_in_gb": {
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
					"total_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"virtual_ip_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_machine_file_system_storage": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gb": {
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

func (d ExascaleDatabaseVirtualMachineClusterDataSource) ModelObject() interface{} {
	return &ExascaleDatabaseVirtualMachineClusterDataModel{}
}

func (d ExascaleDatabaseVirtualMachineClusterDataSource) ResourceType() string {
	return "azurerm_oracle_exascale_database_virtual_machine_cluster"
}

func (d ExascaleDatabaseVirtualMachineClusterDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exadbvmclusters.ValidateExadbVMClusterID
}

func (d ExascaleDatabaseVirtualMachineClusterDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExadbVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ExascaleDatabaseVirtualMachineClusterDataModel
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
					state.DataCollection = FlattenExadbDataCollectionOption(props.DataCollectionOptions)
					state.DisplayName = props.DisplayName
					state.Domain = pointer.From(props.Domain)
					state.EnabledEcpuCount = props.EnabledEcpuCount
					state.ExascaleDatabaseStorageVaultId = props.ExascaleDbStorageVaultId
					state.GridInfrastructureVersion = pointer.From(props.GiVersion)
					state.GridImageOcid = pointer.From(props.GridImageOcid)
					state.GridImageType = string(pointer.From(props.GridImageType))
					state.Hostname = removeHostnameSuffix(props.Hostname)
					state.HostnameActual = props.Hostname
					state.IormConfigCache = flattenIormConfig(props.IormConfigCache)
					state.LicenseModel = string(pointer.From(props.LicenseModel))
					state.LifecycleDetails = pointer.From(props.LifecycleDetails)
					state.LifecycleState = string(*props.LifecycleState)
					state.ListenerPort = pointer.From(props.ListenerPort)
					state.MemorySizeInGb = pointer.From(props.MemorySizeInGbs)
					state.NodeCount = props.NodeCount
					state.NetworkSecurityGroupCidr = FlattenNetworkSecurityGroupCidr(props.NsgCidrs)
					state.NetworkSecurityGroupUrl = pointer.From(props.NsgURL)
					state.OciUrl = pointer.From(props.OciURL)
					state.Ocid = pointer.From(props.Ocid)
					state.PrivateZoneOcid = pointer.From(props.PrivateZoneOcid)
					state.SingleClientAccessNameDnsName = pointer.From(props.ScanDnsName)
					state.SingleClientAccessNameDnsRecordId = pointer.From(props.ScanDnsRecordId)
					state.SingleClientAccessNameIpIds = pointer.From(props.ScanIPIds)
					state.SingleClientAccessNameListenerPortTcp = pointer.From(props.ScanListenerPortTcp)
					state.SingleClientAccessNameListenerPortTcpSsl = pointer.From(props.ScanListenerPortTcpSsl)
					state.Shape = props.Shape
					state.SnapshotFileSystemStorage = flattenExadbVmClusterStorage(props.SnapshotFileSystemStorage)
					state.SshPublicKeys = props.SshPublicKeys
					state.SubnetId = props.SubnetId
					state.SubnetOcid = pointer.From(props.SubnetOcid)
					state.SystemVersion = pointer.From(props.SystemVersion)
					state.TimeZone = pointer.From(props.TimeZone)
					state.TotalEcpuCount = props.TotalEcpuCount
					state.TotalFileSystemStorage = flattenExadbVmClusterStorage(props.TotalFileSystemStorage)
					state.VirtualIpIds = *props.VipIds
					state.VirtualMachineFileSystemStorage = FlattenVMFileSystemStorage(props.VMFileSystemStorage)
					state.VnetId = props.VnetId
					state.ZoneOcid = pointer.From(props.ZoneOcid)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenIormConfig(input *exadbvmclusters.ExadataIormConfig) []IormConfigModel {
	output := make([]IormConfigModel, 0)

	if input != nil {
		var dbIormConfigModel []ExascaleDatabaseIormConfigModel
		if input.DbPlans != nil {
			dbPlans := *input.DbPlans
			for _, dbPlan := range dbPlans {
				dbIormConfigModel = append(dbIormConfigModel, ExascaleDatabaseIormConfigModel{
					DatabaseName:    pointer.From(dbPlan.DbName),
					FlashCacheLimit: pointer.From(dbPlan.FlashCacheLimit),
					Share:           pointer.From(dbPlan.Share),
				})
			}
		}
		return append(output, IormConfigModel{
			DatabasePlans:    dbIormConfigModel,
			LifecycleDetails: pointer.From(input.LifecycleDetails),
			LifecycleState:   string(pointer.From(input.LifecycleState)),
			Objective:        string(pointer.From(input.Objective)),
		})
	}

	return output
}

func flattenExadbVmClusterStorage(input *exadbvmclusters.ExadbVMClusterStorageDetails) []ExascaleDatabaseVirtualMachineClusterStorageModel {
	output := make([]ExascaleDatabaseVirtualMachineClusterStorageModel, 0)
	if input != nil {
		return append(output, ExascaleDatabaseVirtualMachineClusterStorageModel{
			TotalSizeInGb: input.TotalSizeInGbs,
		})
	}
	return output
}
