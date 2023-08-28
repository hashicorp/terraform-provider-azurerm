// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudAppCosmosDbAssociationV0ToV1 struct{}

func (s SpringCloudAppCosmosDbAssociationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"spring_cloud_app_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"cosmosdb_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"api_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"cosmosdb_access_key": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cosmosdb_cassandra_keyspace_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"cosmosdb_gremlin_database_name", "cosmosdb_gremlin_graph_name", "cosmosdb_mongo_database_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_gremlin_database_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			RequiredWith:  []string{"cosmosdb_gremlin_graph_name"},
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_mongo_database_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_gremlin_graph_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			RequiredWith:  []string{"cosmosdb_gremlin_database_name"},
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_mongo_database_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_mongo_database_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_gremlin_database_name", "cosmosdb_gremlin_graph_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_sql_database_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_gremlin_database_name", "cosmosdb_gremlin_graph_name", "cosmosdb_mongo_database_name"},
		},
	}
}

func (s SpringCloudAppCosmosDbAssociationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudAppAssociationIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
