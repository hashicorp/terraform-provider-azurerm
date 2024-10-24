// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbsystemshapes"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DbSystemShapesDataSource struct{}

type DbSystemShapesModel struct {
	DbSystemShapes []DbSystemShapeModel `tfschema:"db_system_shapes"`
	Location       string               `tfschema:"location"`
}

type DbSystemShapeModel struct {
	AvailableCoreCount                 int64   `tfschema:"available_core_count"`
	AvailableCoreCountPerNode          int64   `tfschema:"available_core_count_per_node"`
	AvailableDataStorageInTbs          int64   `tfschema:"available_data_storage_in_tbs"`
	AvailableDataStoragePerServerInTbs float64 `tfschema:"available_data_storage_per_server_in_tbs"`
	AvailableDbNodePerNodeInGbs        int64   `tfschema:"available_db_node_per_node_in_gbs"`
	AvailableDbNodeStorageInGbs        int64   `tfschema:"available_db_node_storage_in_gbs"`
	AvailableMemoryInGbs               int64   `tfschema:"available_memory_in_gbs"`
	AvailableMemoryPerNodeInGbs        int64   `tfschema:"available_memory_per_node_in_gbs"`
	CoreCountIncrement                 int64   `tfschema:"core_count_increment"`
	MaxStorageCount                    int64   `tfschema:"max_storage_count"`
	MaximumNodeCount                   int64   `tfschema:"maximum_node_count"`
	MinCoreCountPerNode                int64   `tfschema:"min_core_count_per_node"`
	MinDataStorageInTbs                int64   `tfschema:"min_data_storage_in_tbs"`
	MinDbNodeStoragePerNodeInGbs       int64   `tfschema:"min_db_node_storage_per_node_in_gbs"`
	MinMemoryPerNodeInGbs              int64   `tfschema:"min_memory_per_node_in_gbs"`
	MinStorageCount                    int64   `tfschema:"min_storage_count"`
	MinimumCoreCount                   int64   `tfschema:"minimum_core_count"`
	MinimumNodeCount                   int64   `tfschema:"minimum_node_count"`
	RuntimeMinimumCoreCount            int64   `tfschema:"runtime_minimum_core_count"`
	ShapeFamily                        string  `tfschema:"shape_family"`
}

func (d DbSystemShapesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
	}
}

func (d DbSystemShapesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"db_system_shapes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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
					"core_count_increment": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"max_storage_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"maximum_node_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"min_core_count_per_node": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"min_data_storage_in_tbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"min_db_node_storage_per_node_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"min_memory_per_node_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"min_storage_count": {
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

			resp, err := client.ListByLocation(ctx, id)
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
							AvailableCoreCount:                 pointer.From(properties.AvailableCoreCount),
							AvailableCoreCountPerNode:          pointer.From(properties.AvailableCoreCountPerNode),
							AvailableDataStorageInTbs:          pointer.From(properties.AvailableDataStorageInTbs),
							AvailableDataStoragePerServerInTbs: pointer.From(properties.AvailableDataStoragePerServerInTbs),
							AvailableDbNodePerNodeInGbs:        pointer.From(properties.AvailableDbNodePerNodeInGbs),
							AvailableDbNodeStorageInGbs:        pointer.From(properties.AvailableDbNodeStorageInGbs),
							AvailableMemoryInGbs:               pointer.From(properties.AvailableMemoryInGbs),
							AvailableMemoryPerNodeInGbs:        pointer.From(properties.AvailableMemoryPerNodeInGbs),
							CoreCountIncrement:                 pointer.From(properties.CoreCountIncrement),
							MaxStorageCount:                    pointer.From(properties.MaxStorageCount),
							MaximumNodeCount:                   pointer.From(properties.MaximumNodeCount),
							MinCoreCountPerNode:                pointer.From(properties.MinCoreCountPerNode),
							MinDataStorageInTbs:                pointer.From(properties.MinDataStorageInTbs),
							MinDbNodeStoragePerNodeInGbs:       pointer.From(properties.MinDbNodeStoragePerNodeInGbs),
							MinMemoryPerNodeInGbs:              pointer.From(properties.MinMemoryPerNodeInGbs),
							MinStorageCount:                    pointer.From(properties.MinStorageCount),
							MinimumCoreCount:                   pointer.From(properties.MinimumCoreCount),
							MinimumNodeCount:                   pointer.From(properties.MinimumNodeCount),
							RuntimeMinimumCoreCount:            pointer.From(properties.RuntimeMinimumCoreCount),
							ShapeFamily:                        pointer.From(properties.ShapeFamily),
						})
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
