// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exadbvmclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbnodes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExascaleDatabaseNodesDataSource struct{}

type ExascaleDbNodesDataModel struct {
	ExascaleVmClusterId string                    `tfschema:"exascale_database_virtual_machine_cluster_id"`
	ExascaleDBNodes     []ExascaleDbNodeDataModel `tfschema:"exascale_database_nodes"`
}

type ExascaleDbNodeDataModel struct {
	AdditionalDetails           string `tfschema:"additional_details"`
	CpuCoreCount                int64  `tfschema:"cpu_core_count"`
	DatabaseNodeStorageSizeInGb int64  `tfschema:"database_node_storage_size_in_gb"`
	FaultDomain                 string `tfschema:"fault_domain"`
	Hostname                    string `tfschema:"hostname"`
	LifecycleState              string `tfschema:"lifecycle_state"`
	MaintenanceType             string `tfschema:"maintenance_type"`
	MemorySizeInGb              int64  `tfschema:"memory_size_in_gb"`
	Ocid                        string `tfschema:"ocid"`
	SoftwareStorageSizeInGb     int64  `tfschema:"software_storage_size_in_gb"`
	TimeMaintenanceWindowEnd    string `tfschema:"time_maintenance_window_end"`
	TimeMaintenanceWindowStart  string `tfschema:"time_maintenance_window_start"`
	TotalCpuCoreCount           int64  `tfschema:"total_cpu_core_count"`
}

func (d ExascaleDatabaseNodesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"exascale_database_virtual_machine_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: exadbvmclusters.ValidateExadbVMClusterID,
		},
	}
}

func (d ExascaleDatabaseNodesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"exascale_database_nodes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_details": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"cpu_core_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"database_node_storage_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"fault_domain": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hostname": {
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

					"memory_size_in_gb": {
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

					"time_maintenance_window_end": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"time_maintenance_window_start": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"total_cpu_core_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d ExascaleDatabaseNodesDataSource) ModelObject() interface{} {
	return &ExascaleDbNodesDataModel{}
}

func (d ExascaleDatabaseNodesDataSource) ResourceType() string {
	return "azurerm_oracle_exascale_database_nodes"
}

func (d ExascaleDatabaseNodesDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exascaledbnodes.ValidateExadbVMClusterDbNodeID
}

func (d ExascaleDatabaseNodesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbNodes
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := ExascaleDbNodesDataModel{
				ExascaleDBNodes: make([]ExascaleDbNodeDataModel, 0),
			}

			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			parsedExadbVmClusterId, err := exadbvmclusters.ParseExadbVMClusterID(state.ExascaleVmClusterId)
			if err != nil {
				return fmt.Errorf("decoding id: %+v", err)
			}
			id := exascaledbnodes.NewExadbVMClusterID(subscriptionId, parsedExadbVmClusterId.ResourceGroupName, parsedExadbVmClusterId.ExadbVmClusterName)

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
						dbNode := ExascaleDbNodeDataModel{
							AdditionalDetails:           pointer.From(props.AdditionalDetails),
							CpuCoreCount:                pointer.From(props.CpuCoreCount),
							DatabaseNodeStorageSizeInGb: pointer.From(props.DbNodeStorageSizeInGbs),
							FaultDomain:                 pointer.From(props.FaultDomain),
							Hostname:                    pointer.From(props.Hostname),
							LifecycleState:              string(*props.LifecycleState),
							MaintenanceType:             pointer.From(props.MaintenanceType),
							MemorySizeInGb:              pointer.From(props.MemorySizeInGbs),
							Ocid:                        props.Ocid,
							SoftwareStorageSizeInGb:     pointer.From(props.SoftwareStorageSizeInGb),
							TimeMaintenanceWindowEnd:    pointer.From(props.TimeMaintenanceWindowEnd),
							TimeMaintenanceWindowStart:  pointer.From(props.TimeMaintenanceWindowStart),
							TotalCpuCoreCount:           pointer.From(props.TotalCPUCoreCount),
						}
						state.ExascaleDBNodes = append(state.ExascaleDBNodes, dbNode)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
