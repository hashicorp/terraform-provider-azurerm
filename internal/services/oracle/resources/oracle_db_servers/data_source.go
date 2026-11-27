// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_db_servers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type DBServersDataModel struct {
	CloudExadataInfrastructureName string              `tfschema:"cloud_exadata_infrastructure_name"`
	DBServers                      []DBServerDataModel `tfschema:"db_servers"`
	ResourceGroupName              string              `tfschema:"resource_group_name"`
}

type DBServerDataModel struct {
	AutonomousVMClusterIds      []string `tfschema:"autonomous_vm_cluster_ids"`
	AutonomousVirtualMachineIds []string `tfschema:"autonomous_virtual_machine_ds"`
	CompartmentId               string   `tfschema:"compartment_id"`
	ComputeModel                string   `tfschema:"compute_model"`
	CpuCoreCount                int64    `tfschema:"cpu_core_count"`
	DbNodeIds                   []string `tfschema:"db_node_ids"`
	DbNodeStorageSizeInGbs      int64    `tfschema:"db_node_storage_size_in_gbs"`
	DisplayName                 string   `tfschema:"display_name"`
	ExadataInfrastructureId     string   `tfschema:"exadata_infrastructure_id"`
	LifecycleDetails            string   `tfschema:"lifecycle_details"`
	LifecycleState              string   `tfschema:"lifecycle_state"`
	MaxCPUCount                 int64    `tfschema:"max_cpu_count"`
	MaxDbNodeStorageInGbs       int64    `tfschema:"max_db_node_storage_in_gbs"`
	MaxMemoryInGbs              int64    `tfschema:"max_memory_in_gbs"`
	MemorySizeInGbs             int64    `tfschema:"memory_size_in_gbs"`
	Ocid                        string   `tfschema:"ocid"`
	Shape                       string   `tfschema:"shape"`
	TimeCreated                 string   `tfschema:"time_created"`
	VMClusterIds                []string `tfschema:"vm_cluster_ids"`
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"cloud_exadata_infrastructure_name": {
			Description:  "The name of the Cloud Exadata Infrastructure.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"db_servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"autonomous_vm_cluster_ids": {
						Description: "The list of [OCIDs](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Autonomous VM Clusters associated with the DB Server.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"autonomous_virtual_machine_ds": {
						Description: "The list of [OCIDs](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Autonomous Virtual Machines associated with the DB Server.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"compartment_id": {
						Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the compartment.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"compute_model": {
						Description: "The compute model of the Exadata Infrastructure.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"cpu_core_count": {
						Description: "The number of CPU cores enabled on the DB Server.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"db_node_ids": {
						Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Db nodes associated with the DB Server.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"db_node_storage_size_in_gbs": {
						Description: "The allocated local node storage in GBs on the DB Server.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"display_name": {
						Description: "The user-friendly name for the DB Server. The name does not need to be unique.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"exadata_infrastructure_id": {
						Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exadata infrastructure.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"lifecycle_details": {
						Description: "Additional information about the current lifecycle state.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"lifecycle_state": {
						Description: "The current state of the DB Server.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"max_cpu_count": {
						Description: "The total number of CPU cores available.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"max_db_node_storage_in_gbs": {
						Description: "The total local node storage available in GBs.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"max_memory_in_gbs": {
						Description: "The total memory available in GBs.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"memory_size_in_gbs": {
						Description: "The allocated memory in GBs on the DB Server.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"ocid": {
						Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DB Server.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"shape": {
						Description: "The shape of the DB Server. The shape determines the amount of CPU, storage, and memory resources available.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"time_created": {
						Description: "The date and time that the DB Server was created.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"vm_cluster_ids": {
						Description: "The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the VM Clusters associated with the DB Server.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (d DataSource) ModelObject() any {
	return &DBServersDataModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_oracle_db_servers"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about existing DB Servers."
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbservers.ValidateDbServerID
}

func (d DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbServers
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DBServersDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := dbservers.NewCloudExadataInfrastructureID(subscriptionId, state.ResourceGroupName, state.CloudExadataInfrastructureName)

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
						dbServer := DBServerDataModel{
							AutonomousVMClusterIds:      pointer.From(props.AutonomousVMClusterIds),
							AutonomousVirtualMachineIds: pointer.From(props.AutonomousVirtualMachineIds),
							CompartmentId:               pointer.From(props.CompartmentId),
							ComputeModel:                pointer.FromEnum(props.ComputeModel),
							CpuCoreCount:                pointer.From(props.CpuCoreCount),
							DbNodeIds:                   pointer.From(props.DbNodeIds),
							DbNodeStorageSizeInGbs:      pointer.From(props.DbNodeStorageSizeInGbs),
							DisplayName:                 pointer.From(props.DisplayName),
							ExadataInfrastructureId:     pointer.From(props.ExadataInfrastructureId),
							LifecycleDetails:            pointer.From(props.LifecycleDetails),
							LifecycleState:              pointer.FromEnum(props.LifecycleState),
							MaxCPUCount:                 pointer.From(props.MaxCPUCount),
							MaxDbNodeStorageInGbs:       pointer.From(props.MaxDbNodeStorageInGbs),
							MaxMemoryInGbs:              pointer.From(props.MaxMemoryInGbs),
							MemorySizeInGbs:             pointer.From(props.MemorySizeInGbs),
							Ocid:                        pointer.From(props.Ocid),
							Shape:                       pointer.From(props.Shape),
							TimeCreated:                 pointer.From(props.TimeCreated),
							VMClusterIds:                pointer.From(props.VMClusterIds),
						}
						state.DBServers = append(state.DBServers, dbServer)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
