// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbnodes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DBNodesDataSource struct{}

type DBNodesDataModel struct {
	CloudVmClusterId string            `tfschema:"cloud_vm_cluster_id"`
	DBNodes          []DBNodeDataModel `tfschema:"db_nodes"`
}

type DBNodeDataModel struct {
	AdditionalDetails          string `tfschema:"additional_details"`
	BackupIPId                 string `tfschema:"backup_ip_id"`
	BackupVnic2Id              string `tfschema:"backup_vnic_2_id"`
	BackupVnicId               string `tfschema:"backup_vnic_id"`
	CpuCoreCount               int64  `tfschema:"cpu_core_count"`
	DbNodeStorageSizeInGbs     int64  `tfschema:"db_node_storage_size_in_gbs"`
	DbServerId                 string `tfschema:"db_server_id"`
	DbSystemId                 string `tfschema:"db_system_id"`
	FaultDomain                string `tfschema:"fault_domain"`
	HostIPId                   string `tfschema:"host_ip_id"`
	Hostname                   string `tfschema:"hostname"`
	LifecycleDetails           string `tfschema:"lifecycle_details"`
	LifecycleState             string `tfschema:"lifecycle_state"`
	MaintenanceType            string `tfschema:"maintenance_type"`
	MemorySizeInGbs            int64  `tfschema:"memory_size_in_gbs"`
	Ocid                       string `tfschema:"ocid"`
	SoftwareStorageSizeInGb    int64  `tfschema:"software_storage_size_in_gb"`
	TimeCreated                string `tfschema:"time_created"`
	TimeMaintenanceWindowEnd   string `tfschema:"time_maintenance_window_end"`
	TimeMaintenanceWindowStart string `tfschema:"time_maintenance_window_start"`
	Vnic2Id                    string `tfschema:"vnic_2_id"`
	VnicId                     string `tfschema:"vnic_id"`
}

func (d DBNodesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cloud_vm_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: cloudvmclusters.ValidateCloudVMClusterID,
		},
	}
}

func (d DBNodesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"db_nodes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_details": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"backup_ip_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"backup_vnic_2_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"backup_vnic_id": {
						Type:     pluginsdk.TypeString,
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

					"db_server_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"db_system_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"fault_domain": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"host_ip_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hostname": {
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

					"maintenance_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"memory_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"ocid": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"software_storage_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"time_created": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"time_maintenance_window_end": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"time_maintenance_window_start": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"vnic_2_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"vnic_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d DBNodesDataSource) ModelObject() interface{} {
	return &DBNodesDataModel{}
}

func (d DBNodesDataSource) ResourceType() string {
	return "azurerm_oracle_db_nodes"
}

func (d DBNodesDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbnodes.ValidateDbNodeID
}

func (d DBNodesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbNodes
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := DBNodesDataModel{
				DBNodes: make([]DBNodeDataModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parsedCloudVmClusterId, err := cloudvmclusters.ParseCloudVMClusterID(state.CloudVmClusterId)
			if err != nil {
				return fmt.Errorf("decoding id: %+v", err)
			}
			id := dbnodes.NewCloudVMClusterID(subscriptionId, parsedCloudVmClusterId.ResourceGroupName, parsedCloudVmClusterId.CloudVmClusterName)

			resp, err := client.ListByCloudVMCluster(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if props := element.Properties; props != nil {
						dbNode := DBNodeDataModel{
							AdditionalDetails:          pointer.From(props.AdditionalDetails),
							BackupIPId:                 pointer.From(props.BackupIPId),
							BackupVnicId:               pointer.From(props.BackupVnicId),
							BackupVnic2Id:              pointer.From(props.BackupVnic2Id),
							CpuCoreCount:               pointer.From(props.CpuCoreCount),
							DbNodeStorageSizeInGbs:     pointer.From(props.DbNodeStorageSizeInGbs),
							DbServerId:                 pointer.From(props.DbServerId),
							DbSystemId:                 props.DbSystemId,
							FaultDomain:                pointer.From(props.FaultDomain),
							HostIPId:                   pointer.From(props.HostIPId),
							Hostname:                   pointer.From(props.Hostname),
							LifecycleDetails:           pointer.From(props.LifecycleDetails),
							LifecycleState:             string(props.LifecycleState),
							MaintenanceType:            string(pointer.From(props.MaintenanceType)),
							MemorySizeInGbs:            pointer.From(props.MemorySizeInGbs),
							Ocid:                       props.Ocid,
							SoftwareStorageSizeInGb:    pointer.From(props.SoftwareStorageSizeInGb),
							TimeCreated:                props.TimeCreated,
							TimeMaintenanceWindowEnd:   pointer.From(props.TimeMaintenanceWindowEnd),
							TimeMaintenanceWindowStart: pointer.From(props.TimeMaintenanceWindowStart),
							VnicId:                     props.VnicId,
							Vnic2Id:                    pointer.From(props.Vnic2Id),
						}
						state.DBNodes = append(state.DBNodes, dbNode)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
