// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_db_nodes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbnodes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

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

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cloud_vm_cluster_id": {
			Description:  "The ID of the Cloud VM cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: cloudvmclusters.ValidateCloudVMClusterID,
		},
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"db_nodes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_details": {
						Description: "Additional information about the planned maintenance.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"backup_ip_id": {
						Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the backup IP address associated with the database node. Use this OCID with either the [GetPrivateIp](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PrivateIp/GetPrivateIp) or the [GetPublicIpByPrivateIpId](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PublicIp/GetPublicIpByPrivateIpId) API to get the IP address needed to make a database connection.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"backup_vnic_2_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"backup_vnic_id": {
						Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the backup VNIC.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"cpu_core_count": {
						Description: "The number of CPU cores enabled on the DB node.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"db_node_storage_size_in_gbs": {
						Description: "The allocated local node storage in GBs on the DB node.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"db_server_id": {
						Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the ExaCC DB server associated with the database node.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"db_system_id": {
						Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the DB system.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"fault_domain": {
						Description: "The name of the Fault Domain the instance is contained in.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"host_ip_id": {
						Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the host IP address associated with the database node. Use this OCID with either the [GetPrivateIp](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PrivateIp/GetPrivateIp) or the [GetPublicIpByPrivateIpId](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PublicIp/GetPublicIpByPrivateIpId) API to get the IP address needed to make a database connection.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"hostname": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"lifecycle_details": {
						Description: "Information about the current lifecycle details.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"lifecycle_state": {
						Description: "Information about the current lifecycle state.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"maintenance_type": {
						Description: "The type of database node maintenance.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"memory_size_in_gbs": {
						Description: "The allocated memory in GBs on the DB Node.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"ocid": {
						Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DB node.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"software_storage_size_in_gb": {
						Description: "The size (in GB) of the block storage volume allocation for the DB system. This attribute applies only for Virtual Machine DB systems.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"time_created": {
						Description: "The date and time that the DB node was created.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"time_maintenance_window_end": {
						Description: "End date and time of maintenance window.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"time_maintenance_window_start": {
						Description: "Start date and time of maintenance window.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"vnic_2_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"vnic_id": {
						Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the VNIC.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d DataSource) ModelObject() any {
	return &DBNodesDataModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_oracle_db_nodes"
}

func (d DataSource) Description() string {
	return "Lists the database nodes for the specified Cloud VM Cluster."
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbnodes.ValidateDbNodeID
}

func (d DataSource) Read() sdk.ResourceFunc {
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
				return err
			}
			id := dbnodes.NewCloudVMClusterID(subscriptionId, parsedCloudVmClusterId.ResourceGroupName, parsedCloudVmClusterId.CloudVmClusterName)

			resp, err := client.ListByParent(ctx, id)
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
							MaintenanceType:            pointer.FromEnum(props.MaintenanceType),
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
