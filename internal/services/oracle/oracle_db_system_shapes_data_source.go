// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbsystemshapes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DbSystemShapesDataSource struct{}

type DbSystemShapesModel struct {
	DbSystemShapes []DbSystemShapeModel `tfschema:"db_system_shapes"`
	Location       string               `tfschema:"location"`
	Zone           string               `tfschema:"zone"`
}

type DbSystemShapeModel struct {
	AreServerTypesSupported            bool     `tfschema:"are_server_types_supported"`
	AvailableCoreCount                 int64    `tfschema:"available_core_count"`
	AvailableCoreCountPerNode          int64    `tfschema:"available_core_count_per_node"`
	AvailableDataStorageInTbs          int64    `tfschema:"available_data_storage_in_tbs"`
	AvailableDataStoragePerServerInTbs float64  `tfschema:"available_data_storage_per_server_in_tbs"`
	AvailableDbNodePerNodeInGbs        int64    `tfschema:"available_db_node_per_node_in_gbs"`
	AvailableDbNodeStorageInGbs        int64    `tfschema:"available_db_node_storage_in_gbs"`
	AvailableMemoryInGbs               int64    `tfschema:"available_memory_in_gbs"`
	AvailableMemoryPerNodeInGbs        int64    `tfschema:"available_memory_per_node_in_gbs"`
	ComputeModel                       string   `tfschema:"compute_model"`
	CoreCountIncrement                 int64    `tfschema:"core_count_increment"`
	DisplayName                        string   `tfschema:"display_name"`
	MaxStorageCount                    int64    `tfschema:"maximum_storage_count"`
	MaximumNodeCount                   int64    `tfschema:"maximum_node_count"`
	MinCoreCountPerNode                int64    `tfschema:"minimum_core_count_per_node"`
	MinDataStorageInTbs                int64    `tfschema:"minimum_data_storage_in_tbs"`
	MinDbNodeStoragePerNodeInGbs       int64    `tfschema:"minimum_db_node_storage_per_node_in_gbs"`
	MinMemoryPerNodeInGbs              int64    `tfschema:"minimum_memory_per_node_in_gbs"`
	MinStorageCount                    int64    `tfschema:"minimum_storage_count"`
	MinimumCoreCount                   int64    `tfschema:"minimum_core_count"`
	MinimumNodeCount                   int64    `tfschema:"minimum_node_count"`
	RuntimeMinimumCoreCount            int64    `tfschema:"runtime_minimum_core_count"`
	ShapeAttributes                    []string `tfschema:"shape_attributes"`
	ShapeFamily                        string   `tfschema:"shape_family"`
}

func (d DbSystemShapesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"zone": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Filter the versions by zone",
		},
	}
}

func (d DbSystemShapesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"db_system_shapes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"are_server_types_supported": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"available_core_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"available_core_count_per_node": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"available_data_storage_in_tbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"available_data_storage_per_server_in_tbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"available_db_node_per_node_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"available_db_node_storage_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"available_memory_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"available_memory_per_node_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"compute_model": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"core_count_increment": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"display_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"maximum_storage_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"maximum_node_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"minimum_core_count_per_node": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"minimum_data_storage_in_tbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"minimum_db_node_storage_per_node_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"minimum_memory_per_node_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"minimum_storage_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"minimum_core_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"minimum_node_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"runtime_minimum_core_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"shape_attributes": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"shape_family": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d DbSystemShapesDataSource) ModelObject() interface{} {
	return &DbSystemShapesModel{}
}

func (d DbSystemShapesDataSource) ResourceType() string {
	return "azurerm_oracle_db_system_shapes"
}

func (d DbSystemShapesDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbsystemshapes.ValidateDbSystemShapeID
}

func (d DbSystemShapesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbSystemShapes
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := DbSystemShapesModel{
				DbSystemShapes: make([]DbSystemShapeModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := dbsystemshapes.NewLocationID(subscriptionId, state.Location)

			options := dbsystemshapes.ListByLocationOperationOptions{}

			if state.Zone != "" {
				options.Zone = &state.Zone
			}

			if state.Zone == "" {
				log.Printf("[WARN] DbSystem shapes data source: Zone parameter is empty. This may result in unfiltered results from the API. Consider specifying Zone for more precise results in the desired zone.")
			}

			resp, err := client.ListByLocation(ctx, id, options)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if props := element.Properties; props != nil {
						state.DbSystemShapes = append(state.DbSystemShapes, DbSystemShapeModel{
							AreServerTypesSupported:            pointer.From(props.AreServerTypesSupported),
							AvailableCoreCount:                 props.AvailableCoreCount,
							AvailableCoreCountPerNode:          pointer.From(props.AvailableCoreCountPerNode),
							AvailableDataStorageInTbs:          pointer.From(props.AvailableDataStorageInTbs),
							AvailableDataStoragePerServerInTbs: pointer.From(props.AvailableDataStoragePerServerInTbs),
							AvailableDbNodePerNodeInGbs:        pointer.From(props.AvailableDbNodePerNodeInGbs),
							AvailableDbNodeStorageInGbs:        pointer.From(props.AvailableDbNodeStorageInGbs),
							AvailableMemoryInGbs:               pointer.From(props.AvailableMemoryInGbs),
							AvailableMemoryPerNodeInGbs:        pointer.From(props.AvailableMemoryPerNodeInGbs),
							ComputeModel:                       pointer.FromEnum(props.ComputeModel),
							CoreCountIncrement:                 pointer.From(props.CoreCountIncrement),
							DisplayName:                        pointer.From(props.DisplayName),
							MaxStorageCount:                    pointer.From(props.MaxStorageCount),
							MaximumNodeCount:                   pointer.From(props.MaximumNodeCount),
							MinCoreCountPerNode:                pointer.From(props.MinCoreCountPerNode),
							MinDataStorageInTbs:                pointer.From(props.MinDataStorageInTbs),
							MinDbNodeStoragePerNodeInGbs:       pointer.From(props.MinDbNodeStoragePerNodeInGbs),
							MinMemoryPerNodeInGbs:              pointer.From(props.MinMemoryPerNodeInGbs),
							MinStorageCount:                    pointer.From(props.MinStorageCount),
							MinimumCoreCount:                   pointer.From(props.MinimumCoreCount),
							MinimumNodeCount:                   pointer.From(props.MinimumNodeCount),
							RuntimeMinimumCoreCount:            pointer.From(props.RuntimeMinimumCoreCount),
							ShapeAttributes:                    pointer.From(props.ShapeAttributes),
							ShapeFamily:                        pointer.From(props.ShapeFamily),
						})
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
