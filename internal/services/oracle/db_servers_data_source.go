// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DBServersDataSource struct{}

type DBServersDataModel struct {
	CloudExadataInfrastructureName string              `tfschema:"cloud_exadata_infrastructure_name"`
	DBServers                      []DBServerDataModel `tfschema:"db_servers"`
	ResourceGroupName              string              `tfschema:"resource_group_name"`
}

type DBServerDataModel struct {
	AutonomousVMClusterIds      []string `tfschema:"autonomous_vm_cluster_ids"`
	AutonomousVirtualMachineIds []string `tfschema:"autonomous_virtual_machine_ds"`
	CompartmentId               string   `tfschema:"compartment_id"`
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

func (d DBServersDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"cloud_exadata_infrastructure_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ExadataName,
		},
	}
}

func (d DBServersDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"db_servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"autonomous_vm_cluster_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"autonomous_virtual_machine_ds": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"compartment_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"cpu_core_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"db_node_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"db_node_storage_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"display_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"exadata_infrastructure_id": {
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

					"max_cpu_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"max_db_node_storage_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"max_memory_in_gbs": {
						Type:     pluginsdk.TypeInt,
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

					"shape": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"time_created": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"vm_cluster_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (d DBServersDataSource) ModelObject() interface{} {
	return &DBServersDataModel{}
}

func (d DBServersDataSource) ResourceType() string {
	return "azurerm_oracle_db_servers"
}

func (d DBServersDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbservers.ValidateDbServerID
}

func (d DBServersDataSource) Read() sdk.ResourceFunc {
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

			resp, err := client.ListByCloudExadataInfrastructure(ctx, id)
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
							CpuCoreCount:                pointer.From(props.CpuCoreCount),
							DbNodeIds:                   pointer.From(props.DbNodeIds),
							DbNodeStorageSizeInGbs:      pointer.From(props.DbNodeStorageSizeInGbs),
							DisplayName:                 pointer.From(props.DisplayName),
							ExadataInfrastructureId:     pointer.From(props.ExadataInfrastructureId),
							LifecycleDetails:            pointer.From(props.LifecycleDetails),
							LifecycleState:              string(pointer.From(props.LifecycleState)),
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
