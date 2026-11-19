// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_cosmosdb_association

import (
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	cosmosValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func springCloudAppCosmosdbAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudAppAssociationName,
		},

		"spring_cloud_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(appplatform_rm.ParseAppIDInsensitively),
		},

		"cosmosdb_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cosmosdb.ValidateDatabaseAccountID,
		},

		"api_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				springCloudAppCosmosDbAssociationAPITypeCassandra,
				springCloudAppCosmosDbAssociationAPITypeGremlin,
				springCloudAppCosmosDbAssociationAPITypeMongo,
				springCloudAppCosmosDbAssociationAPITypeSql,
				springCloudAppCosmosDbAssociationAPITypeTable,
			}, false),
		},

		"cosmosdb_access_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cosmosdb_cassandra_keyspace_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  cosmosValidate.CosmosEntityName,
			ConflictsWith: []string{"cosmosdb_gremlin_database_name", "cosmosdb_gremlin_graph_name", "cosmosdb_mongo_database_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_gremlin_database_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  cosmosValidate.CosmosEntityName,
			RequiredWith:  []string{"cosmosdb_gremlin_graph_name"},
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_mongo_database_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_gremlin_graph_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  cosmosValidate.CosmosEntityName,
			RequiredWith:  []string{"cosmosdb_gremlin_database_name"},
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_mongo_database_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_mongo_database_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  cosmosValidate.CosmosEntityName,
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_gremlin_database_name", "cosmosdb_gremlin_graph_name", "cosmosdb_sql_database_name"},
		},

		"cosmosdb_sql_database_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  cosmosValidate.CosmosEntityName,
			ConflictsWith: []string{"cosmosdb_cassandra_keyspace_name", "cosmosdb_gremlin_database_name", "cosmosdb_gremlin_graph_name", "cosmosdb_mongo_database_name"},
		},
	}
}
