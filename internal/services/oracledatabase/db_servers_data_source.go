// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbservers"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DBServersDataSource struct{}

type DBServersDataModel struct {
	DBServers []DBServerDataModel `tfschema:"db_servers"`
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
	ProvisioningState           string   `tfschema:"provisioning_state"`
	Shape                       string   `tfschema:"shape"`
	TimeCreated                 string   `tfschema:"time_created"`
	VMClusterIds                []string `tfschema:"vm_cluster_ids"`
}

func (d DBServersDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cloud_exadata_infrastructure_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
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

					"provisioning_state": {
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
	return nil
}

func (d DBServersDataSource) ResourceType() string {
	return "azurerm_oracledatabase_db_servers"
}

func (d DBServersDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbservers.ValidateDbServerID
}

func (d DBServersDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.DbServers
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := dbservers.NewCloudExadataInfrastructureID(subscriptionId,
				metadata.ResourceData.Get("resource_group_name").(string),
				metadata.ResourceData.Get("cloud_exadata_infrastructure_name").(string))

			resp, err := client.ListByCloudExadataInfrastructure(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				output := DBServersDataModel{
					DBServers: make([]DBServerDataModel, 0),
				}
				for _, element := range *model {
					if element.Properties != nil {
						properties := element.Properties
						dbServer := DBServerDataModel{
							AutonomousVMClusterIds:      pointer.From(properties.AutonomousVMClusterIds),
							AutonomousVirtualMachineIds: pointer.From(properties.AutonomousVirtualMachineIds),
							CompartmentId:               pointer.From(properties.CompartmentId),
							CpuCoreCount:                pointer.From(properties.CpuCoreCount),
							DbNodeIds:                   pointer.From(properties.DbNodeIds),
							DbNodeStorageSizeInGbs:      pointer.From(properties.DbNodeStorageSizeInGbs),
							DisplayName:                 pointer.From(properties.DisplayName),
							ExadataInfrastructureId:     pointer.From(properties.ExadataInfrastructureId),
							LifecycleDetails:            pointer.From(properties.LifecycleDetails),
							LifecycleState:              string(pointer.From(properties.LifecycleState)),
							MaxCPUCount:                 pointer.From(properties.MaxCPUCount),
							MaxDbNodeStorageInGbs:       pointer.From(properties.MaxDbNodeStorageInGbs),
							MaxMemoryInGbs:              pointer.From(properties.MaxMemoryInGbs),
							MemorySizeInGbs:             pointer.From(properties.MemorySizeInGbs),
							Ocid:                        pointer.From(properties.Ocid),
							ProvisioningState:           string(pointer.From(properties.ProvisioningState)),
							Shape:                       pointer.From(properties.Shape),
							TimeCreated:                 pointer.From(properties.TimeCreated),
							VMClusterIds:                pointer.From(properties.VMClusterIds),
						}
						output.DBServers = append(output.DBServers, dbServer)
					}
				}
				metadata.SetID(id)
				if err := metadata.Encode(&output); err != nil {
					return fmt.Errorf("encoding %s: %+v", id, err)
				}
			}
			return nil
		},
	}
}
